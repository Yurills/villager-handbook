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
	stepEventClaim
	stepEventFact
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
		menuItems:  []string{"Add Event", "Show Recommend", "Show Stat", "Show Log", "Exit"},
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
					m.step = stepEventType
					m.textInput.Placeholder = "Type (accuse/claim/fact) "
					m.textInput.Reset()

				} else if selected == "Show Recommend" {
					m.gameEngine.PredictMove()
					stats := m.gameEngine.GetPredictStat()
					recommendations := m.gameEngine.GetRecommend(stats, m.playerInfo.WarewolfCount)

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
		if m.step >= stepEventType && m.step <= stepEventResult {
			if msg.Type == tea.KeyEnter {
				input := m.textInput.Value()

				switch input {
				// --- D. EVENT INPUT FLOW (Type(accuse) -> Actor -> Target -> Result) ---
				case "accuse":
					switch m.step {
					case stepEventType:
						m.tempEvent.Type = input
						m.step = stepEventActor
						m.textInput.Placeholder = "Player ID (e.g. 0)"

					case stepEventActor:
						id, _ := strconv.Atoi(input)
						m.tempEvent.Actor = id
						m.step = stepEventTarget
						m.textInput.Placeholder = "Player ID (e.g. 1)"

					case stepEventTarget:
						id, _ := strconv.Atoi(input)
						m.tempEvent.Target = id
						m.step = stepEventResult
						m.textInput.Placeholder = "Result (Werewolf/Seer/Villager)"

					case stepEventResult:
						m.tempEvent.Result = input

						// --- EXECUTE MOVE IN ENGINE ---
						m.gameEngine.ProcessMove(m.tempEvent)

						log := fmt.Sprintf("Player %d %s Player %d -> %s",
							m.tempEvent.Actor, m.tempEvent.Type, m.tempEvent.Target, m.tempEvent.Result)

						m.feedback = log
						logEvents = append(logEvents, log) // Store log
						// Return to menu
						m.step = stepGameMenu
					}
					m.textInput.Reset()
					return m, nil
				// --- D. EVENT INPUT FLOW (Type(claim) -> Actor -> Result) ---
				case "claim":
					switch m.step {
					case stepEventType:
						m.tempEvent.Type = input
						m.step = stepEventClaim
						m.textInput.Placeholder = "Player ID (e.g. 0)"

					case stepEventClaim:
						id, _ := strconv.Atoi(input)
						m.tempEvent.Actor = id
						m.step = stepEventResult
						m.textInput.Placeholder = "Player ID (e.g. 1)"

					case stepEventResult:
						m.tempEvent.Result = input

						// --- EXECUTE MOVE IN ENGINE ---
						m.gameEngine.ProcessMove(m.tempEvent)

						log := fmt.Sprintf("Player %d %s themself as %s",
							m.tempEvent.Actor, m.tempEvent.Type, m.tempEvent.Result)

						m.feedback = log
						logEvents = append(logEvents, log) // Store log
						// Return to menu
						m.step = stepGameMenu
					}
					m.textInput.Reset()
					return m, nil
				// --- D. EVENT INPUT FLOW (Type(fact) -> Target -> Result) ---
				case "fact":
					switch m.step {
					case stepEventType:
						m.tempEvent.Type = input
						m.step = stepEventFact
						m.textInput.Placeholder = "Player ID (e.g. 0)"

					case stepEventFact:
						id, _ := strconv.Atoi(input)
						m.tempEvent.Target = id
						m.step = stepEventResult
						m.textInput.Placeholder = "Player ID (e.g. 1)"

					case stepEventResult:
						m.tempEvent.Result = input

						// --- EXECUTE MOVE IN ENGINE ---
						m.gameEngine.ProcessMove(m.tempEvent)

						log := fmt.Sprintf("Player %d is %s",
							m.tempEvent.Target, m.tempEvent.Result)

						m.feedback = log
						logEvents = append(logEvents, log) // Store log
						// Return to menu
						m.step = stepGameMenu
					}
					m.textInput.Reset()
					return m, nil
				}
			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
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
	case stepEventType:
		return sectionStyle.Render("What is the action? (accuse, claim, fact)\n" + m.textInput.View())
	case stepEventActor:
		return sectionStyle.Render("Who is accusing? (Player ID)\n" + m.textInput.View())
	case stepEventTarget:
		return sectionStyle.Render(
			fmt.Sprintf("Actor: %d\nWho are being accused? (Player ID)\n%s",
				m.tempEvent.Actor,
				m.textInput.View(),
			),
		)
	case stepEventClaim:
		return sectionStyle.Render("Who is Claiming? (Player ID)\n" + m.textInput.View())
	case stepEventFact:
		return sectionStyle.Render("This Player role is certainly a fact. (Player ID)\n" + m.textInput.View())
	case stepEventResult:
		return sectionStyle.Render(
			fmt.Sprintf("Action: %s\nInteraction: What is the result/role? (Werewolf, Seer, Villager)\n%s",
				m.tempEvent.Type, m.textInput.View()),
		)
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
