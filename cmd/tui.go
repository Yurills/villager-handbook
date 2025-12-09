package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	// Your internal packages
	"github.com/Yurills/villager-handbook/internal/engine"
	"github.com/Yurills/villager-handbook/internal/model"
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
	stepEventActor
	stepEventTarget
	stepEventType
	stepEventResult
)

// --- 2. MAIN MODEL ---
type bubbleModel struct {
	step      step
	textInput textinput.Model
	
	// Data Storage
	playerInfo PlayerInfo          // Stores Setup counts
	gameEngine *engine.Engine      // Stores the actual Game Engine
	tempEvent  model.Interaction   // Stores input while creating an event

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

func initialModel() bubbleModel {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	return bubbleModel{
		step:       stepWelcome,
		textInput:  ti,
		playerInfo: PlayerInfo{},
		menuItems:  []string{"Add Event", "Show Recommend", "Exit"},
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
				if m.menuCursor > 0 { m.menuCursor-- }
			case "down", "j":
				if m.menuCursor < len(m.menuItems)-1 { m.menuCursor++ }
			case "enter":
				selected := m.menuItems[m.menuCursor]
				
				if selected == "Exit" {
					return m, tea.Quit
				
				} else if selected == "Add Event" {
					// Start the Event Input Flow
					m.step = stepEventActor
					m.textInput.Placeholder = "Actor ID (e.g. 0)"
					m.textInput.Reset()
				
				} else if selected == "Show Recommend" {
					// --- CALL YOUR ENGINE STATS ---
					stats := m.gameEngine.GetStats()
					
					// Format the stats slice into a single string for display
					var sb strings.Builder
					sb.WriteString("=== RECOMMENDATIONS ===\n")
					for _, s := range stats {
						sb.WriteString(fmt.Sprintf("- %s\n", s))
					}
					m.feedback = sb.String()
				}
			}
			return m, nil
		}

		// --- D. EVENT INPUT FLOW (Actor -> Target -> Type -> Result) ---
		if m.step >= stepEventActor && m.step <= stepEventResult {
			if msg.Type == tea.KeyEnter {
				input := m.textInput.Value()

				switch m.step {
				case stepEventActor:
					id, _ := strconv.Atoi(input)
					m.tempEvent.Actor = id
					m.step = stepEventTarget
					m.textInput.Placeholder = "Target ID (e.g. 1)"

				case stepEventTarget:
					id, _ := strconv.Atoi(input)
					m.tempEvent.Target = id
					m.step = stepEventType
					m.textInput.Placeholder = "Type (accuse/claim/fact)"

				case stepEventType:
					m.tempEvent.Type = input
					m.step = stepEventResult
					m.textInput.Placeholder = "Result (Werewolf/Seer/Villager)"

				case stepEventResult:
					m.tempEvent.Result = input
					
					// --- EXECUTE MOVE IN ENGINE ---
					m.gameEngine.ProcessMove(m.tempEvent)
					
					m.feedback = fmt.Sprintf("Processed: Actor %d %s Target %d -> %s", 
						m.tempEvent.Actor, m.tempEvent.Type, m.tempEvent.Target, m.tempEvent.Result)
					
					// Return to menu
					m.step = stepGameMenu
				}
				m.textInput.Reset()
				return m, nil
			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

// --- 4. VIEW (RENDER) ---
func (m bubbleModel) View() string {
	s := strings.Builder{}

	switch m.step {
	case stepWelcome:
		s.WriteString("=== VILLAGER HANDBOOK ===\nPress Enter to Start")

	case stepInputVillager:
		s.WriteString("Enter Villager Count:\n" + m.textInput.View())

	case stepInputSear:
		s.WriteString("Enter Seer Count:\n" + m.textInput.View())

	case stepInputWarewolf:
		s.WriteString("Enter Werewolf Count:\n" + m.textInput.View())

	case stepGameMenu:
		s.WriteString("=== MAIN MENU ===\n")
		s.WriteString(fmt.Sprintf("Players: %d (v:%d s:%d w:%d) | Worlds: %d\n\n", 
			m.playerInfo.TotalPlayer, m.playerInfo.VillagerCount, m.playerInfo.SearCount, m.playerInfo.WarewolfCount, len(m.gameEngine.Worlds)))

		for i, item := range m.menuItems {
			cursor := " "
			if m.menuCursor == i { cursor = ">" }
			s.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
		}
		
		if m.feedback != "" {
			s.WriteString("\n----------------\n" + m.feedback + "\n----------------")
		}

	// VIEW FOR EVENT INPUTS
	case stepEventActor:
		s.WriteString("Interaction: Who is acting? (Actor ID)\n" + m.textInput.View())
	case stepEventTarget:
		s.WriteString(fmt.Sprintf("Actor: %d\nInteraction: Who are they targeting? (Target ID)\n%s", m.tempEvent.Actor, m.textInput.View()))
	case stepEventType:
		s.WriteString(fmt.Sprintf("Actor: %d -> Target: %d\nInteraction: What is the action? (accuse, claim, fact)\n%s", m.tempEvent.Actor, m.tempEvent.Target, m.textInput.View()))
	case stepEventResult:
		s.WriteString(fmt.Sprintf("Action: %s\nInteraction: What is the result/role? (Werewolf, Seer, Villager)\n%s", m.tempEvent.Type, m.textInput.View()))
	}

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}