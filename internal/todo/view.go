package todo

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	// Base styles
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1).
			Bold(true)

	// Task styles
	taskStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedTaskStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#EE6FF8")).
				Background(lipgloss.Color("#313244")).
				PaddingLeft(2)

	completedTaskStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A6E3A1")).
				Strikethrough(true)

	// Priority styles
	highPriorityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F38BA8"))

	mediumPriorityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAB387"))

	lowPriorityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F9E2AF"))

	// Context styles
	contextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#89B4FA")).
			Bold(true)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F38BA8")).
			Bold(true)

	// Input styles
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Margin(1)
)

// View is the main rendering function for the application, also a core part of the
// Bubble Tea architecture. It returns a string that represents the UI to be drawn
// to the terminal. The runtime calls this whenever the model is updated.
func (m Model) View() string {
	if m.HelpVisible {
		return m.renderFullHelpView()
	}

	// Delegate to a specific rendering function based on the current ViewMode.
	// This acts as a router for the UI, ensuring the correct screen is displayed.
	switch m.ViewMode {
	case InputView:
		return m.RenderInputView()
	case DateInputView:
		return m.RenderDateInputView()
	case RemoveTagView:
		return m.RenderRemoveTagView()
	case KanbanView:
		return m.RenderKanbanView()
	case StatsView:
		return m.RenderStatsView()
	default:
		return m.RenderNormalView()
	}
}

// renderFullHelpView renders a centered, modal-like view of all keybindings.
func (m Model) renderFullHelpView() string {
	// Style the help box to look like a modal
	helpBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89B4FA")).
		Padding(1, 2)

	m.Help.ShowAll = true
	helpContent := m.Help.View(m.KeyMap)
	// Add a title to the help view
	titledHelp := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Keybindings"),
		helpContent,
	)

	return lipgloss.Place(m.WindowWidth, m.WindowHeight, lipgloss.Center, lipgloss.Center, helpBoxStyle.Render(titledHelp))
}

// RenderNormalView renders the main task list screen.
func (m Model) RenderNormalView() string {
	var mainContent strings.Builder

	// 1. Render the header, showing the current context.
	contextText := fmt.Sprintf("Context: %s", m.CurrentContext)
	mainContent.WriteString(titleStyle.Render(contextText) + "\n\n")

	// 2. Get the tasks for the current context and render them.
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		if len(m.Contexts) == 0 {
			mainContent.WriteString("No contexts exist. Press 'n' to create one.\n")
		} else {
			mainContent.WriteString("No tasks in this context. Press 'a' to add one.\n")
		}
	} else {
		for i, task := range tasks {
			taskLine := m.RenderTask(task, i == m.SelectedIndex, m.MovingMode && task.ID == m.MovingTaskID)
			mainContent.WriteString(taskLine + "\n")
		}
	}

	// 3. Display an error message if one exists.
	if m.ErrorMessage != "" {
		mainContent.WriteString("\n" + errorStyle.Render(m.ErrorMessage) + "\n")
	}

	return baseStyle.Render(mainContent.String())
}

// RenderTask renders a single line representing a task.
// It applies different styling based on whether the task is completed,
// selected by the cursor, or currently being moved.
func (m Model) RenderTask(task Task, selected, moving bool) string {
	// Checkbox indicates completion status.
	checkbox := "[ ]"
	if task.Checked {
		checkbox = "[✓]"
	}

	// Priority is shown with exclamation marks.
	priority := ""
	switch task.Priority {
	case "high":
		priority = highPriorityStyle.Render("!!! ")
	case "medium":
		priority = mediumPriorityStyle.Render("!! ")
	case "low":
		priority = lowPriorityStyle.Render("! ")
	}

	taskText := task.Task

	// Tags are appended to the task text.
	tags := ""
	if len(task.Tags) > 0 {
		tags = " > " + strings.Join(task.Tags, ", ")
	}

	// Due date is shown at the end.
	dueDate := ""
	if task.DueDate != "" {
		dueDate = fmt.Sprintf(" [Due: %s]", task.DueDate)
	}

	text := fmt.Sprintf("%s %s%s%s", checkbox, taskText, tags, dueDate)

	// Apply styles based on the task's state.
	style := taskStyle
	if task.Checked {
		style = completedTaskStyle
	}

	// The 'selected' style has precedence over the base or completed style.
	if selected {
		style = style.Background(lipgloss.Color("#313244"))
	}

	// The 'moving' style is applied on top of other styles.
	if moving {
		style = style.Bold(true)
	}

	return priority + style.Render(text)
}

// RenderInputView renders input dialogs and places them in the center of the screen.
func (m Model) RenderInputView() string {
	content := inputStyle.Render(
		fmt.Sprintf("%s\n\n%s", m.InputPrompt, m.TextInput.View()),
	)
	return lipgloss.Place(m.WindowWidth, m.WindowHeight, lipgloss.Center, lipgloss.Center, content)
}

// RenderDateInputView renders due date input dialog
func (m Model) RenderDateInputView() string {
	var content strings.Builder
	content.WriteString("Set due date (YYYY-MM-DD):\n\n")
	inputs := []string{
		fmt.Sprintf("Day: %s", m.DateInputs[0].View()),
		fmt.Sprintf("Month: %s", m.DateInputs[1].View()),
		fmt.Sprintf("Year: %s", m.DateInputs[2].View()),
	}
	for i, input := range inputs {
		if i == m.DateInputIndex {
			content.WriteString(selectedTaskStyle.Render(input) + "\n")
		} else {
			content.WriteString(input + "\n")
		}
	}
	return inputStyle.Render(content.String())
}

// RenderRemoveTagView renders remove tag view
func (m Model) RenderRemoveTagView() string {
	var content strings.Builder
	content.WriteString("Select tags to remove:\n\n")
	task := m.GetCurrentTask()
	for i, tag := range task.Tags {
		checkbox := "[ ]"
		if m.RemoveTagChecks[i] {
			checkbox = "[✓]"
		}
		line := fmt.Sprintf("%s %s", checkbox, tag)
		if i == m.RemoveTagIndex {
			content.WriteString(selectedTaskStyle.Render(line) + "\n")
		} else {
			content.WriteString(line + "\n")
		}
	}
	return inputStyle.Render(content.String())
}

// RenderKanbanView renders the kanban board with horizontal and vertical scrolling.
func (m Model) RenderKanbanView() string {
	var content strings.Builder
	title := titleStyle.Render("Kanban View (←/→/↑/↓ scroll, esc to return)")
	content.WriteString(title + "\n")

	if len(m.Contexts) == 0 {
		content.WriteString("No contexts available.\n")
		return baseStyle.Render(content.String())
	}

	// --- Horizontal Scrolling Logic ---
	const (
		fixedColWidth  = 35
		separatorWidth = 3
	)

	// Calculate how many columns can fit on screen
	numVisibleCols := m.WindowWidth / (fixedColWidth + separatorWidth)
	if numVisibleCols < 1 {
		numVisibleCols = 1
	}

	// Ensure scroll position is within bounds
	if m.KanbanScrollX > len(m.Contexts)-numVisibleCols {
		m.KanbanScrollX = max(0, len(m.Contexts)-numVisibleCols)
	}
	if m.KanbanScrollX < 0 {
		m.KanbanScrollX = 0
	}

	// Get the slice of contexts that are currently visible
	startCol := m.KanbanScrollX
	endCol := min(startCol+numVisibleCols, len(m.Contexts))
	visibleContexts := m.Contexts[startCol:endCol]

	// Style for wrapping text within a column.
	columnStyle := lipgloss.NewStyle().Width(fixedColWidth).Padding(0, 1)
	taskTextStyle := lipgloss.NewStyle().Width(fixedColWidth - 2)

	// Render each visible context as a column.
	var columns []string
	for _, context := range visibleContexts {
		var column strings.Builder
		header := contextStyle.Render(context)
		column.WriteString(header + "\n")
		column.WriteString(strings.Repeat("─", fixedColWidth) + "\n")

		tasks := m.GetTasksForContext(context)
		for _, task := range tasks {
			var taskLine strings.Builder
			if task.Checked {
				taskLine.WriteString("✓ ")
			} else {
				taskLine.WriteString("• ")
			}
			fullTaskText := task.Task
			if len(task.Tags) > 0 {
				fullTaskText += " > " + strings.Join(task.Tags, ", ")
			}
			if task.DueDate != "" {
				fullTaskText += fmt.Sprintf(" [Due: %s]", task.DueDate)
			}
			wrappedText := taskTextStyle.Render(fullTaskText)
			if task.Checked {
				taskLine.WriteString(completedTaskStyle.Render(wrappedText))
			} else {
				taskLine.WriteString(wrappedText)
			}
			column.WriteString(taskLine.String() + "\n")
		}
		columns = append(columns, columnStyle.Render(column.String()))
	}

	// Join columns horizontally.
	board := lipgloss.JoinHorizontal(lipgloss.Top, columns...)
	boardLines := strings.Split(board, "\n")

	// --- Vertical Scrolling Logic ---
	top := m.KanbanScrollY
	bottom := top + m.WindowHeight - lipgloss.Height(title) - 1
	if top < 0 {
		top = 0
	}
	if bottom > len(boardLines) {
		bottom = len(boardLines)
	}
	if top > bottom {
		top = max(0, bottom-m.WindowHeight)
		m.KanbanScrollY = top
	}

	visibleLines := boardLines[top:bottom]
	content.WriteString(strings.Join(visibleLines, "\n"))

	return baseStyle.Render(content.String())
}

// RenderStatsView renders the statistics view
func (m Model) RenderStatsView() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("Statistics (ESC to return)") + "\n\n")

	// Overall stats
	total := len(m.Tasks)
	completed := 0
	for _, task := range m.Tasks {
		if task.Checked {
			completed++
		}
	}

	completionRate := 0.0
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
	}

	content.WriteString(fmt.Sprintf("Total Tasks: %d\n", total))
	content.WriteString(fmt.Sprintf("Completed: %d (%.1f%%)\n\n", completed, completionRate))

	// Context stats
	content.WriteString("Context Statistics:\n")
	for _, context := range m.Contexts {
		tasks := m.GetTasksForContext(context)
		ctxTotal := len(tasks)
		ctxCompleted := 0
		for _, task := range tasks {
			if task.Checked {
				ctxCompleted++
			}
		}

		ctxRate := 0.0
		if ctxTotal > 0 {
			ctxRate = float64(ctxCompleted) / float64(ctxTotal) * 100
		}

		content.WriteString(fmt.Sprintf("  %s: %d/%d (%.1f%%)\n",
			contextStyle.Render(context), ctxCompleted, ctxTotal, ctxRate))
	}

	return baseStyle.Render(content.String())
}
