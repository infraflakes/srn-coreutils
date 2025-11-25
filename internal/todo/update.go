package todo

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
)

// Update is the central message loop for the application, a core part of the Bubble Tea
// architecture. It's called by the runtime whenever an event occurs (e.g., a key press,
// window resize, or a command finishing). Its job is to update the model's state
// based on the event and return the updated model and any new command to run.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle window resize events.
	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		m.Help.Width = msg.Width
		return m, tea.ClearScreen

	// Handle key press events.
	case tea.KeyMsg:
		// On any key press, we clear a previous error message, so it's not sticky.
		m.ErrorMessage = ""

		// The update logic is first delegated based on the current ViewMode.
		// This creates a state machine where key presses have different meanings
		// depending on what the user is currently doing.
		switch m.ViewMode {
		case InputView:
			return m.UpdateInputMode(msg)
		case DateInputView:
			return m.UpdateDateInputMode(msg)
		case RemoveTagView:
			return m.UpdateRemoveTagMode(msg)
		}

		// If not in a special input mode, delegate to the handler for the current major view.
		switch m.ViewMode {
		case NormalView:
			return m.UpdateNormalView(msg)
		case KanbanView:
			return m.UpdateKanbanView(msg)
		case StatsView:
			return m.UpdateStatsView(msg)
		}
	}

	return m, nil
}

// UpdateInputMode handles all key presses when the user is in an input dialog.
func (m Model) UpdateInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	// 'esc' exits the input dialog and returns to the normal task view.
	case key.Matches(msg, m.KeyMap.Back):
		m.ViewMode = NormalView
		return m, nil

	// 'enter' confirms the input.
	case key.Matches(msg, m.KeyMap.Enter):
		input := strings.TrimSpace(m.TextInput.Value())
		m.TextInput.SetValue("")

		// The action taken depends on the specific InputMode we're in.
		switch m.InputMode {
		case AddTaskInput:
			if input != "" {
				m.SaveStateForUndo()
				m.AddTask(input)
			}
		case EditTaskInput:
			if input != "" {
				m.SaveStateForUndo()
				m.EditCurrentTask(input)
			}
		case AddContextInput:
			if input != "" {
				m.AddContext(input)
			}
		case RenameContextInput:
			if input != "" && input != m.CurrentContext {
				m.RenameContext(input)
			}
		case AddTagInput:
			if input != "" {
				m.SaveStateForUndo()
				m.AddTagToCurrentTask(input)
			}
		case DeleteConfirmInput:
			if strings.ToLower(input) == "y" {
				m.SaveStateForUndo()
				m.DeleteContext()
			}
		}

		m.ViewMode = NormalView
		return m, nil
	}

	// For any other key, update the text input component.
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

// UpdateDateInputMode handles key presses for the due date input dialog.
func (m Model) UpdateDateInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.KeyMap.Back):
		m.ViewMode = NormalView
		return m, nil

	case key.Matches(msg, m.KeyMap.Enter):
		day := m.DateInputs[0].Value()
		month := m.DateInputs[1].Value()
		year := m.DateInputs[2].Value()
		dateStr := fmt.Sprintf("%s-%s-%s", year, month, day)
		m.SaveStateForUndo()
		m.SetDueDateForCurrentTask(dateStr)
		m.ViewMode = NormalView
		return m, nil

	case key.Matches(msg, m.KeyMap.Up):
		m.DateInputs[m.DateInputIndex].Blur()
		m.DateInputIndex = (m.DateInputIndex - 1 + 3) % 3
		m.DateInputs[m.DateInputIndex].Focus()

	case key.Matches(msg, m.KeyMap.Down):
		m.DateInputs[m.DateInputIndex].Blur()
		m.DateInputIndex = (m.DateInputIndex + 1) % 3
		m.DateInputs[m.DateInputIndex].Focus()
	}

	m.DateInputs[m.DateInputIndex], cmd = m.DateInputs[m.DateInputIndex].Update(msg)
	return m, cmd
}

// UpdateRemoveTagMode handles key presses for the remove tag dialog.
func (m Model) UpdateRemoveTagMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.KeyMap.Back):
		m.ViewMode = NormalView
		return m, nil

	case key.Matches(msg, m.KeyMap.Enter):
		m.SaveStateForUndo()
		m.RemoveTagsFromCurrentTask()
		m.ViewMode = NormalView
		return m, nil

	case key.Matches(msg, m.KeyMap.Up):
		if m.RemoveTagIndex > 0 {
			m.RemoveTagIndex--
		}

	case key.Matches(msg, m.KeyMap.Down):
		task := m.GetCurrentTask()
		if m.RemoveTagIndex < len(task.Tags)-1 {
			m.RemoveTagIndex++
		}

	case key.Matches(msg, m.KeyMap.Toggle):
		m.RemoveTagChecks[m.RemoveTagIndex] = !m.RemoveTagChecks[m.RemoveTagIndex]
	}

	return m, nil
}

// UpdateNormalView is the key handler for the main task list view.
// It maps keys to specific actions like navigation and task manipulation.
func (m Model) UpdateNormalView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.KeyMap.Quit):
		m.SaveConfig()
		return m, tea.Quit

	case key.Matches(msg, m.KeyMap.Back):
		return m, nil

	// --- Navigation ---
	case key.Matches(msg, m.KeyMap.Up):
		// If in moving mode, the 'up' key moves the selected task.
		// Otherwise, it just moves the cursor.
		if m.MovingMode {
			m.MoveTaskUp()
		} else {
			m.MoveUp()
		}

	case key.Matches(msg, m.KeyMap.Down):
		// If in moving mode, the 'down' key moves the selected task.
		// Otherwise, it just moves the cursor.
		if m.MovingMode {
			m.MoveTaskDown()
		} else {
			m.MoveDown()
		}

	case key.Matches(msg, m.KeyMap.Left):
		m.PreviousContext()

	case key.Matches(msg, m.KeyMap.Right):
		m.NextContext()

	// --- Task Manipulation ---
	case key.Matches(msg, m.KeyMap.Toggle):
		if len(m.GetFilteredTasks()) > 0 {
			m.SaveStateForUndo()
			m.ToggleCurrentTask()
		}

	case key.Matches(msg, m.KeyMap.Add):
		m.ShowInputDialog(AddTaskInput, "Add new task:")

	case key.Matches(msg, m.KeyMap.Edit):
		if len(m.GetFilteredTasks()) > 0 {
			task := m.GetCurrentTask()
			m.ShowInputDialog(EditTaskInput, "Edit task:")
			m.TextInput.SetValue(task.Task)
		}

	case key.Matches(msg, m.KeyMap.Delete):
		if len(m.GetFilteredTasks()) > 0 {
			m.SaveStateForUndo()
			m.DeleteCurrentTask()
		}

	// --- Context Manipulation ---
	case key.Matches(msg, m.KeyMap.AddContext):
		m.ShowInputDialog(AddContextInput, "New context name:")

	case key.Matches(msg, m.KeyMap.RenameContext):
		m.ShowInputDialog(RenameContextInput, "Rename context to:")
		m.TextInput.SetValue(m.CurrentContext)

	case key.Matches(msg, m.KeyMap.DeleteContext):
		if len(m.Contexts) > 1 {
			m.ShowInputDialog(DeleteConfirmInput, fmt.Sprintf("Delete context '%s'? (y/n):", m.CurrentContext))
		} else {
			m.ErrorMessage = "Cannot delete the only context"
		}

	// --- Task Metadata ---
	case key.Matches(msg, m.KeyMap.TogglePriority):
		if len(m.GetFilteredTasks()) > 0 {
			m.SaveStateForUndo()
			m.ToggleCurrentTaskPriority()
		}

	case key.Matches(msg, m.KeyMap.AddTag):
		if len(m.GetFilteredTasks()) > 0 {
			m.ShowInputDialog(AddTagInput, "Add tag:")
		}

	case key.Matches(msg, m.KeyMap.RemoveTag):
		if len(m.GetFilteredTasks()) > 0 {
			m.ShowRemoveTagDialog()
		}

	case key.Matches(msg, m.KeyMap.SetDueDate):
		if len(m.GetFilteredTasks()) > 0 {
			m.ShowDateInputDialog()
		}

	case key.Matches(msg, m.KeyMap.ClearDueDate):
		if len(m.GetFilteredTasks()) > 0 {
			m.SaveStateForUndo()
			m.SetDueDateForCurrentTask("clear")
		}

	// --- View & Mode Switching ---
	case key.Matches(msg, m.KeyMap.KanbanView):
		m.ViewMode = KanbanView

	case key.Matches(msg, m.KeyMap.StatsView):
		m.ViewMode = StatsView

	case key.Matches(msg, m.KeyMap.Undo):
		m.Undo()

	// This is the core of the move functionality. It toggles MovingMode on/off.
	// When entering moving mode, it records the ID of the currently selected task.
	// This ensures that even if the visual selection changes, the application
	// always knows which task was originally intended to be moved.
	case key.Matches(msg, m.KeyMap.Move):
		if len(m.GetFilteredTasks()) > 0 {
			m.MovingMode = !m.MovingMode
			if m.MovingMode {
				// "Pick up" the task by its ID.
				m.MovingTaskID = m.GetCurrentTask().ID
			} else {
				// "Drop" the task and save the new order to the undo history.
				m.SaveStateForUndo()
			}
		}
	}

	return m, nil
}

// UpdateKanbanView handles kanban view updates
func (m Model) UpdateKanbanView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.KeyMap.Back), key.Matches(msg, m.KeyMap.Quit), key.Matches(msg, m.KeyMap.KanbanView):
		m.ViewMode = NormalView
	}
	return m, nil
}

// UpdateStatsView handles stats view updates
func (m Model) UpdateStatsView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.KeyMap.Back), key.Matches(msg, m.KeyMap.Quit), key.Matches(msg, m.KeyMap.StatsView):
		m.ViewMode = NormalView
	}
	return m, nil
}
