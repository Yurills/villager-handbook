package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// Your internal packages
	"github.com/Yurills/villager-handbook/internal/engine"
	"github.com/Yurills/villager-handbook/internal/model"
)

// ─── STYLES ─────────────────────────────────────────

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD166")).
			Background(lipgloss.Color("#073B4C")).
			Padding(1, 2).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD166"))

	sectionStyle = lipgloss.NewStyle().
			MarginTop(1).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#118AB2"))

	menuCursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF476F")).
			Bold(true)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#06D6A0"))

	// Style for unselected items
	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#06D6A0")).
			Padding(0, 1)

	// Style for the SELECTED item (looks like a button)
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#073B4C")).
				Background(lipgloss.Color("#06D6A0")).
				Padding(0, 1).
				Bold(true)

	feedbackStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD166")).
			Foreground(lipgloss.Color("#073B4C")).
			Background(lipgloss.Color("#FFEFD5"))
)

// --- 1. APP STATES ---
type step int

const (
	stepWelcome       step = iota
	stepInputVillager      // Setup: Villagers
	stepInputSear          // Setup: Seers
	stepInputWarewolf      // Setup: Wolves
	stepGameMenu           // Main Menu

	// New States for "Add Event" Flow
	stepEventType
	stepEventActor
	stepEventTarget
	stepEventResult
)

// --- 2. MAIN MODEL ---
type bubbleModel struct {
	step      step
	textInput textinput.Model

	// Data Storage
	playerInfo PlayerInfo        // Stores Setup counts
	gameEngine *engine.Engine    // Stores the actual Game Engine
	tempEvent  model.Interaction // Stores input while creating an event

	// Menu UI
	menuCursor int
	menuItems  []string
	feedback   string // Displays stats or success messages

	// Event Cursor
	selectionCursor int
	typeOptions     []string
	resultOptions   []string
}

// Helper struct for setup
type PlayerInfo struct {
	VillagerCount int
	SearCount     int
	WarewolfCount int
	TotalPlayer   int
}

// Log storage
var logEvents []string

func initialModel() bubbleModel {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 50

	return bubbleModel{
		step:       stepWelcome,
		textInput:  ti,
		playerInfo: PlayerInfo{},
		menuItems:  []string{"Add Event", "Who to Vote", "Who to Investigate", "Show Stat", "Show Log", "Exit"},
		// Init event options
		typeOptions:     []string{"accuse", "claim", "fact"},
		resultOptions:   []string{"Werewolf", "Seer", "Villager"},
		selectionCursor: 0,
	}
}

func (m bubbleModel) Init() tea.Cmd {
	return textinput.Blink
}

// --- 3. UPDATE (LOGIC) ---
func (m bubbleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// --- A. WELCOME ---
		if m.step == stepWelcome {
			if msg.String() == "enter" {
				m.step = stepInputVillager
				m.textInput.Reset()
			}
			return m, nil
		}

		// --- B. SETUP INPUTS ---
		if m.step == stepInputVillager || m.step == stepInputSear || m.step == stepInputWarewolf {
			if msg.Type == tea.KeyEnter {
				val, _ := strconv.Atoi(m.textInput.Value())

				if m.step == stepInputVillager {
					m.playerInfo.VillagerCount = val
					m.step = stepInputSear
				} else if m.step == stepInputSear {
					m.playerInfo.SearCount = val
					m.step = stepInputWarewolf
				} else if m.step == stepInputWarewolf {
					m.playerInfo.WarewolfCount = val // Note: User code used TotalPlayer for wolves in rules? Adjust if needed.

					// FINAL SETUP CALCULATION
					m.playerInfo.TotalPlayer = m.playerInfo.VillagerCount + m.playerInfo.SearCount + m.playerInfo.WarewolfCount

					// --- INITIALIZE YOUR ENGINE HERE ---
					players := make([]int, m.playerInfo.TotalPlayer)
					for i := 0; i < m.playerInfo.TotalPlayer; i++ {
						players[i] = i
					}

					rules := engine.GameRule{
						NumVillagers:  m.playerInfo.VillagerCount,
						NumSeers:      m.playerInfo.SearCount,
						NumWerewolves: m.playerInfo.WarewolfCount,
					}

					// Store the engine in the model
					m.gameEngine = engine.NewEngine(players, rules)

					m.step = stepGameMenu
					m.feedback = "Game Engine Initialized. Ready."
				}
				m.textInput.Reset()
				return m, nil
			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		// --- C. MAIN MENU ---
		if m.step == stepGameMenu {
			switch msg.String() {
			case "up", "k":
				if m.menuCursor > 0 {
					m.menuCursor--
				}
			case "down", "j":
				if m.menuCursor < len(m.menuItems)-1 {
					m.menuCursor++
				}
			case "enter":
				selected := m.menuItems[m.menuCursor]

				if selected == "Exit" {
					return m, tea.Quit

				} else if selected == "Add Event" {
					// Start the Event Input Flow
					m.step = stepEventType // Start from Event Type
					// m.textInput.Placeholder = "Type (accuse/claim/fact) "
					m.selectionCursor = 0
					// m.tempEvent = model.Interaction{} // Reset temp event
					m.textInput.Reset()

				} else if selected == "Who to Investigate" {
					stats := m.gameEngine.LookaheadBestCandidate()
					var sb strings.Builder
					sb.WriteString("=== BEST INVESTIGATE CANDIDATES ===\n")
					recommendation := GetVotingRecommend(stats)
					sb.WriteString(fmt.Sprintf(recommendation))
					m.feedback = sb.String()

				} else if selected == "Who to Vote" {
					stats := m.gameEngine.GetStats()
					recommendations := m.gameEngine.GetRecommend(stats)
					var sb strings.Builder
					sb.WriteString("=== RECOMMENDATIONS ===\n")
					sb.WriteString(fmt.Sprintf("%s\n", recommendations))
					m.feedback = sb.String()

				} else if selected == "Show Stat" {
					// --- CALL YOUR ENGINE STATS ---
					stats := m.gameEngine.GetStats()

					// Format the stats slice into a single string for display
					var sb strings.Builder
					sb.WriteString("=== STATS ===\n")
					for _, s := range stats {
						sb.WriteString(fmt.Sprintf("- %s\n", s))
					}
					m.feedback = sb.String()
				} else if selected == "Show Log" {
					// --- SHOW LOG ---
					showLog := logEvents
					// Format the logs into a single string for display
					var sb strings.Builder
					sb.WriteString("=== EVENT LOGS ===\n")
					for _, s := range showLog {
						sb.WriteString(fmt.Sprintf("- %s\n", s))
					}
					m.feedback = sb.String()
				}
			}
			return m, nil
		}

		// --- D. EVENT INPUT FLOW  ---
		// 1. SELECT TYPE
		if m.step == stepEventType {
			switch msg.String() {
			case "left", "h": // Left
				if m.selectionCursor > 0 {
					m.selectionCursor--
				}
			case "right", "l": // Right
				if m.selectionCursor < len(m.typeOptions)-1 {
					m.selectionCursor++
				}
			case "enter":
				m.tempEvent.Type = m.typeOptions[m.selectionCursor]
				m.step = stepEventActor
				// m.selectionCursor = 0
				m.textInput.Placeholder = "..."
				m.textInput.Reset()
				return m, nil
			}
		}

		if m.step == stepEventActor || m.step == stepEventTarget {
			if msg.Type == tea.KeyEnter {
				input := m.textInput.Value()

				if m.step == stepEventActor {
					id, _ := strconv.Atoi(input)
					m.tempEvent.Actor = id

					// --- BRANCHING LOGIC ---
					if m.tempEvent.Type == "accuse" {
						// Needs Target
						m.step = stepEventTarget
						m.textInput.Placeholder = "Target ID (e.g. 1)"
					} else {
						// Claim/Fact: Target is Actor (or self)
						m.tempEvent.Target = id
						m.step = stepEventResult
						m.selectionCursor = 0 // Reset for result menu
					}

				} else if m.step == stepEventTarget {
					id, _ := strconv.Atoi(input)
					m.tempEvent.Target = id
					m.step = stepEventResult
					m.selectionCursor = 0
				}

				m.textInput.Reset()
				return m, nil
			}

			// !!! CRITICAL FIX: Allow typing
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		if m.step == stepEventResult {
			switch msg.String() {
			case "left", "h": // Left
				if m.selectionCursor > 0 {
					m.selectionCursor--
				}
			case "right", "l": // Right
				if m.selectionCursor < len(m.resultOptions)-1 {
					m.selectionCursor++
				}
			case "enter":
				m.tempEvent.Result = m.resultOptions[m.selectionCursor]
				m.gameEngine.ProcessMove(m.tempEvent)

				log := ""
				if m.tempEvent.Type == "fact" {
					log = fmt.Sprintf("Processed: Player %d is %s",
						m.tempEvent.Target, m.tempEvent.Result)
				} else if m.tempEvent.Type == "claim" || m.tempEvent.Type == "accuse" {
					log = fmt.Sprintf("Processed: Player %d %ss Player %d as %s",
						m.tempEvent.Actor, m.tempEvent.Type, m.tempEvent.Target, m.tempEvent.Result)
				}
				m.feedback = log
				logEvents = append(logEvents, log)
				m.step = stepGameMenu
				return m, nil
			}
		}

	}

	return m, nil
}

func (m bubbleModel) LookaheadBestCandidate(param any) {
	panic("unimplemented")
}

// --- 4. VIEW (RENDER) ---
func (m bubbleModel) View() string {
	var out strings.Builder

	switch m.step {
	case stepWelcome:
		return titleStyle.Render("VILLAGER HANDBOOK") + "\n\nPress Enter to Start"

	case stepInputVillager:
		return sectionStyle.Render("Enter Villager Count:\n" + m.textInput.View())

	case stepInputSear:
		return sectionStyle.Render("Enter Seer Count:\n" + m.textInput.View())

	case stepInputWarewolf:
		return sectionStyle.Render("Enter Werewolf Count:\n" + m.textInput.View())

	case stepGameMenu:
		header := titleStyle.Render("MAIN MENU")
		topbar := fmt.Sprintf(
			"Players: %d (V:%d S:%d W:%d) | Worlds: %d\n",
			m.playerInfo.TotalPlayer,
			m.playerInfo.VillagerCount,
			m.playerInfo.SearCount,
			m.playerInfo.WarewolfCount,
			len(m.gameEngine.Worlds),
		)

		var menu strings.Builder

		for i, item := range m.menuItems {
			cursor := " "
			if m.menuCursor == i {
				cursor = menuCursorStyle.Render(">")
			}
			menu.WriteString(fmt.Sprintf("%s %s\n",
				cursor, menuItemStyle.Render(item)))
		}

		out.WriteString(header)
		out.WriteString("\n")
		out.WriteString(topbar)
		out.WriteString("\n")
		out.WriteString(sectionStyle.Render(menu.String()))

		if m.feedback != "" {
			out.WriteString("\n")
			out.WriteString(feedbackStyle.Render(m.feedback))
		}

		return out.String()

	// VIEW FOR EVENT INPUTS
	// --- 1. EVENT TYPE (Horizontal Selection) ---
	case stepEventType:
		var opts []string
		for i, option := range m.typeOptions {
			if m.selectionCursor == i {
				opts = append(opts, selectedItemStyle.Render(option))
			} else {
				opts = append(opts, normalItemStyle.Render(option))
			}
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, opts...)
		return sectionStyle.Render("What is the action?\n\n" + row)

	// --- 2. EVENT ACTOR (Dynamic Title) ---
	case stepEventActor:
		title := "Who is acting?"
		if m.tempEvent.Type == "accuse" {
			title = "Who is Accusing? (Actor ID)"
		} else if m.tempEvent.Type == "claim" {
			title = "Who is Claiming? (Actor ID)"
		} else if m.tempEvent.Type == "fact" {
			title = "Who is this Fact about? (Actor ID)"
		}
		return sectionStyle.Render(title + "\n" + m.textInput.View())

	// --- 3. EVENT TARGET (Only for Accuse) ---
	case stepEventTarget:
		return sectionStyle.Render(
			fmt.Sprintf("Actor: %d\nWho are they Accusing? (Target ID)\n%s",
				m.tempEvent.Actor, m.textInput.View()),
		)

	// --- 4. EVENT RESULT (Horizontal Selection) ---
	case stepEventResult:
		var opts []string
		for i, option := range m.resultOptions {
			if m.selectionCursor == i {
				opts = append(opts, selectedItemStyle.Render(option))
			} else {
				opts = append(opts, normalItemStyle.Render(option))
			}
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, opts...)

		// Dynamic prompt based on type
		prompt := "What is the result/role?"
		if m.tempEvent.Type == "claim" {
			prompt = fmt.Sprintf("Player %d claims to be:", m.tempEvent.Actor)
		} else if m.tempEvent.Type == "accuse" {
			prompt = fmt.Sprintf("Player %d accuses %d of being:", m.tempEvent.Actor, m.tempEvent.Target)
		}

		return sectionStyle.Render(prompt + "\n\n" + row)
	}

	return out.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func GetVotingRecommend(entropy []model.LookaheadResult) string {

	if len(entropy) == 0 {
		return "No players to recommend."
	}

	// 1. Find the lowest entropy
	lowest := entropy[0].Entropy
	for _, e := range entropy {
		if e.Entropy < lowest {
			lowest = e.Entropy
		}
	}

	// 2. Collect all players with the lowest entropy
	result := "Recommend investigating these players:\n"
	for _, e := range entropy {
		if e.Entropy == lowest {
			result += fmt.Sprintf(
				"- Player %d (Entropy: %f)\n",
				e.ID, e.Entropy,
			)
		}
	}

	return result
}
