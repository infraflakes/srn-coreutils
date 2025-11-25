package todo

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// ============================================================================
// UI & Dialog Helpers
// ============================================================================

// ShowInputDialog switches the view to the text input mode.
// It resets the text field to ensure no stale data remains from previous edits.
func (m *Model) ShowInputDialog(mode InputMode, prompt string) {
	m.ViewMode = InputView
	m.InputMode = mode
	m.InputPrompt = prompt
	m.TextInput.SetValue("") // Clear previous input
	m.TextInput.Focus()
}

// ShowDateInputDialog prepares the three-part date input fields.
// It pre-fills the inputs with the current date to provide a convenient starting point
// for the user, rather than forcing them to type the full date from scratch.
func (m *Model) ShowDateInputDialog() {
	m.ViewMode = DateInputView
	m.DateInputIndex = 0

	now := time.Now()
	// Pre-fill with today's date
	m.DateInputs[0].SetValue(fmt.Sprintf("%02d", now.Day()))
	m.DateInputs[1].SetValue(fmt.Sprintf("%02d", now.Month()))
	m.DateInputs[2].SetValue(fmt.Sprintf("%d", now.Year()))

	for i := range m.DateInputs {
		m.DateInputs[i].Focus()
	}
}

// ShowRemoveTagDialog initiates the tag removal flow.
// It creates a dynamic list of checkboxes corresponding to the current task's tags.
// If the task has no tags, it aborts early to prevent showing an empty modal.
func (m *Model) ShowRemoveTagDialog() {
	task := m.GetCurrentTask()
	if len(task.Tags) == 0 {
		m.ErrorMessage = "No tags to remove"
		return
	}

	m.ViewMode = RemoveTagView
	m.RemoveTagIndex = 0
	// precise allocation: we know exactly how many checks we need
	m.RemoveTagChecks = make([]bool, len(task.Tags))
}

// ============================================================================
// Data Retrieval & Navigation
// ============================================================================

// GetFilteredTasks returns the subset of tasks visible to the user.
// ARCHITECTURE NOTE: This creates a NEW slice. Modifying the order of this slice
// does not affect m.Tasks (the source of truth). To modify data, we must always
// lookup the task in m.Tasks by its unique ID.
func (m *Model) GetFilteredTasks() []Task {
	// We delegate to GetTasksForContext to centralize the filtering logic.
	return m.GetTasksForContext(m.CurrentContext)
}

// GetTasksForContext filters the master task list for a specific context.
// It is used both for the main view and for context management logic.
func (m *Model) GetTasksForContext(context string) []Task {
	var filtered []Task
	for _, task := range m.Tasks {
		if task.Context == context {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// GetCurrentTask safely retrieves the task under the user's cursor.
// It handles edge cases where the list might be empty or the cursor index
// might be stale (out of bounds) due to a recent deletion or context switch.
func (m *Model) GetCurrentTask() Task {
	tasks := m.GetFilteredTasks()

	// Safety check: If list is empty or index is invalid, return zero-value Task.
	// This prevents index-out-of-range panics during render cycles.
	if len(tasks) == 0 || m.SelectedIndex >= len(tasks) {
		return Task{}
	}
	return tasks[m.SelectedIndex]
}

// MoveUp moves the cursor visually upwards.
// It implements "wrapping" logic: moving up from the top item jumps to the bottom.
func (m *Model) MoveUp() {
	tasks := m.GetFilteredTasks()
	if len(tasks) > 0 {
		// (Current - 1 + Length) % Length handles the wrap-around for index 0 safely.
		m.SelectedIndex = (m.SelectedIndex - 1 + len(tasks)) % len(tasks)
	}
}

// MoveDown moves the cursor visually downwards.
// It wraps around to the top if the user moves down past the last item.
func (m *Model) MoveDown() {
	tasks := m.GetFilteredTasks()
	if len(tasks) > 0 {
		m.SelectedIndex = (m.SelectedIndex + 1) % len(tasks)
	}
}

// findTaskIndexByID searches the MASTER list (m.Tasks) for a specific ID.
// This is the bridge between the "Filtered View" and the "Source of Truth".
// Returns -1 if the ID is not found (though this should theoretically be impossible
// if the ID came from the filtered list).
func (m *Model) findTaskIndexByID(id int) int {
	return slices.IndexFunc(m.Tasks, func(t Task) bool {
		return t.ID == id
	})
}

// MoveTaskUp reorders items in the master list based on a visual action.
// COMPLEXITY: We cannot simply swap indices in the filtered list. We must:
// 1. Identify the two tasks involved in the visual swap.
// 2. Locate their ACTUAL positions in the master m.Tasks slice.
// 3. Swap them in the master slice.
func (m *Model) MoveTaskUp() {
	tasks := m.GetFilteredTasks()

	// Cannot move up if we are already at the top
	if m.SelectedIndex > 0 {
		taskToMove := tasks[m.SelectedIndex]
		taskToSwapWith := tasks[m.SelectedIndex-1]

		// Find where these tasks live in the real data
		idxMove := m.findTaskIndexByID(taskToMove.ID)
		idxSwap := m.findTaskIndexByID(taskToSwapWith.ID)

		// Perform the swap only if both tasks exist in the master list
		if idxMove != -1 && idxSwap != -1 {
			m.Tasks[idxMove], m.Tasks[idxSwap] = m.Tasks[idxSwap], m.Tasks[idxMove]

			// Move the cursor along with the item so the user follows the task they are moving
			m.SelectedIndex--
		}
	}
}

// MoveTaskDown reorders items downwards.
// See MoveTaskUp for architectural details on why we map back to IDs.
func (m *Model) MoveTaskDown() {
	tasks := m.GetFilteredTasks()

	// Cannot move down if we are at the bottom
	if m.SelectedIndex < len(tasks)-1 {
		taskToMove := tasks[m.SelectedIndex]
		taskToSwapWith := tasks[m.SelectedIndex+1]

		idxMove := m.findTaskIndexByID(taskToMove.ID)
		idxSwap := m.findTaskIndexByID(taskToSwapWith.ID)

		if idxMove != -1 && idxSwap != -1 {
			m.Tasks[idxMove], m.Tasks[idxSwap] = m.Tasks[idxSwap], m.Tasks[idxMove]
			m.SelectedIndex++
		}
	}
}

// ============================================================================
// Context Management
// ============================================================================

// NextContext cycles to the next available context tab.
func (m *Model) NextContext() {
	if len(m.Contexts) > 0 {
		// Modern Go: Find index of current string in slice
		currentIdx := slices.Index(m.Contexts, m.CurrentContext)
		if currentIdx == -1 {
			currentIdx = 0 // Fallback if state is desynchronized
		}

		nextIdx := (currentIdx + 1) % len(m.Contexts)
		m.CurrentContext = m.Contexts[nextIdx]

		// Reset cursor to top when switching views to prevent out-of-bounds
		m.SelectedIndex = 0
	}
}

// PreviousContext cycles to the previous context tab.
func (m *Model) PreviousContext() {
	if len(m.Contexts) > 0 {
		currentIdx := slices.Index(m.Contexts, m.CurrentContext)
		if currentIdx == -1 {
			currentIdx = 0
		}

		// Add len() before modulo to handle negative wrapping correctly
		prevIdx := (currentIdx - 1 + len(m.Contexts)) % len(m.Contexts)
		m.CurrentContext = m.Contexts[prevIdx]
		m.SelectedIndex = 0
	}
}

// AddContext inserts a new context if it is unique.
func (m *Model) AddContext(contextName string) {
	// Validation: Contexts must be unique
	if slices.Contains(m.Contexts, contextName) {
		m.ErrorMessage = "Context already exists"
		return
	}

	m.Contexts = append(m.Contexts, contextName)
	// Immediately switch user to the new context
	m.CurrentContext = contextName
	m.SelectedIndex = 0
}

// RenameContext updates a context name in the list AND in all associated tasks.
// This ensures data consistency (referential integrity) between the definition
// of a context and the tasks assigned to it.
func (m *Model) RenameContext(newName string) {
	if newName == m.CurrentContext {
		return // No change needed
	}
	if slices.Contains(m.Contexts, newName) {
		m.ErrorMessage = "Context name already exists"
		return
	}

	oldName := m.CurrentContext

	// 1. Update the context registry
	if idx := slices.Index(m.Contexts, oldName); idx != -1 {
		m.Contexts[idx] = newName
	}

	// 2. Update all tasks that belonged to the old context name
	for i := range m.Tasks {
		if m.Tasks[i].Context == oldName {
			m.Tasks[i].Context = newName
		}
	}

	m.CurrentContext = newName
}

// DeleteContext removes a context and ALL tasks within it.
// This is a destructive operation.
func (m *Model) DeleteContext() {
	// Prevent deleting the last remaining context to keep the UI usable.
	if len(m.Contexts) <= 1 {
		m.ErrorMessage = "Cannot delete the only context"
		return
	}

	// 1. Cascade Delete: Remove tasks associated with this context.
	// slices.DeleteFunc modifies the slice in-place, removing elements where the func returns true.
	m.Tasks = slices.DeleteFunc(m.Tasks, func(t Task) bool {
		return t.Context == m.CurrentContext
	})

	// 2. Remove the context from the registry.
	if idx := slices.Index(m.Contexts, m.CurrentContext); idx != -1 {
		m.Contexts = slices.Delete(m.Contexts, idx, idx+1)
	}

	// 3. Fallback: Switch view to the first available context.
	if len(m.Contexts) > 0 {
		m.CurrentContext = m.Contexts[0]
		m.SelectedIndex = 0
	}
}

// UpdateContexts synchronizes the m.Contexts list with the actual tasks.
// Use Case: This is often called after Undo operations or file loads to ensure
// that if a task exists with context "Work", "Work" appears in the tabs list.
func (m *Model) UpdateContexts() {
	// Use a map to deduplicate context names efficiently
	uniqueContexts := make(map[string]bool)

	// 1. Gather contexts from current tasks
	for _, task := range m.Tasks {
		uniqueContexts[task.Context] = true
	}

	// 2. Gather currently known contexts (preserves empty contexts if we want to keep them)
	for _, ctx := range m.Contexts {
		uniqueContexts[ctx] = true
	}

	// 3. Flatten map back to slice
	m.Contexts = make([]string, 0, len(uniqueContexts))
	for context := range uniqueContexts {
		m.Contexts = append(m.Contexts, context)
	}

	// Sort for consistent UI navigation
	slices.Sort(m.Contexts)

	// 4. Validation: Ensure CurrentContext points to something valid.
	// If the current context disappeared (e.g., via Undo), reset to the first available.
	if m.CurrentContext == "" || !slices.Contains(m.Contexts, m.CurrentContext) {
		if len(m.Contexts) > 0 {
			m.CurrentContext = m.Contexts[0]
		} else {
			// Absolute fallback if everything is empty
			m.CurrentContext = "Work"
			m.Contexts = []string{"Work"}
		}
	}
}

// ============================================================================
// Task Modification (Edit, Check, Delete)
// ============================================================================

// ToggleCurrentTask switches the "checked" status of the selected task.
func (m *Model) ToggleCurrentTask() {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}

	// Map visual selection back to master ID
	targetID := tasks[m.SelectedIndex].ID
	if idx := m.findTaskIndexByID(targetID); idx != -1 {
		m.Tasks[idx].Checked = !m.Tasks[idx].Checked
	}
}

// AddTask appends a new task to the master list.
func (m *Model) AddTask(taskText string) {
	newTask := Task{
		ID:      m.NextID,
		Task:    taskText,
		Checked: false,
		Context: m.CurrentContext, // Inherit context from current view
	}
	m.Tasks = append(m.Tasks, newTask)
	m.NextID++

	// Auto-scroll to the bottom of the list so the user sees the new item
	filtered := m.GetFilteredTasks()
	m.SelectedIndex = len(filtered) - 1
}

// EditCurrentTask updates the text content of the selected task.
func (m *Model) EditCurrentTask(newText string) {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}

	targetID := tasks[m.SelectedIndex].ID
	if idx := m.findTaskIndexByID(targetID); idx != -1 {
		m.Tasks[idx].Task = newText
	}
}

// DeleteCurrentTask removes the selected task from the master list.
func (m *Model) DeleteCurrentTask() {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}

	targetID := tasks[m.SelectedIndex].ID
	if idx := m.findTaskIndexByID(targetID); idx != -1 {
		// slices.Delete handles the slice manipulation (shifting elements) efficiently.
		m.Tasks = slices.Delete(m.Tasks, idx, idx+1)
	}

	// Post-deletion cleanup: Ensure cursor doesn't hang off the end of the list.
	newTasks := m.GetFilteredTasks()
	if m.SelectedIndex >= len(newTasks) && len(newTasks) > 0 {
		m.SelectedIndex = len(newTasks) - 1
	}
}

// SetDueDateForCurrentTask parses and assigns a date string.
// It enforces the "YYYY-MM-DD" format.
func (m *Model) SetDueDateForCurrentTask(dateStr string) {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}

	targetID := tasks[m.SelectedIndex].ID
	idx := m.findTaskIndexByID(targetID)
	if idx == -1 {
		return
	}

	// Handle clearing the date
	if strings.ToLower(dateStr) == "clear" {
		m.Tasks[idx].DueDate = ""
		return
	}

	if dateStr == "" {
		return
	}

	// Modern Go: time.DateOnly ("2006-01-02") is a built-in layout constant (Go 1.20+).
	// This avoids manual splitting and integer conversion.
	_, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		m.ErrorMessage = "Invalid date format. Use YYYY-MM-DD"
		return
	}

	m.Tasks[idx].DueDate = dateStr
}

// ToggleCurrentTaskPriority cycles through priority levels: None -> Low -> Medium -> High.
func (m *Model) ToggleCurrentTaskPriority() {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}

	targetID := tasks[m.SelectedIndex].ID
	idx := m.findTaskIndexByID(targetID)
	if idx == -1 {
		return
	}

	priorities := []string{"", "low", "medium", "high"}

	// Find current priority in the list
	currentPrioIdx := slices.Index(priorities, m.Tasks[idx].Priority)
	if currentPrioIdx == -1 {
		currentPrioIdx = 0 // Default to "" if current value is invalid
	}

	// Cycle to next priority
	nextIdx := (currentPrioIdx + 1) % len(priorities)
	m.Tasks[idx].Priority = priorities[nextIdx]
}

// ============================================================================
// Tag Management
// ============================================================================

// AddTagToCurrentTask appends a tag if it doesn't already exist on the task.
func (m *Model) AddTagToCurrentTask(tag string) {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}

	targetID := tasks[m.SelectedIndex].ID
	if idx := m.findTaskIndexByID(targetID); idx != -1 {
		// Prevent duplicate tags on the same task
		if !slices.Contains(m.Tasks[idx].Tags, tag) {
			m.Tasks[idx].Tags = append(m.Tasks[idx].Tags, tag)
		}
	}
}

// RemoveTagsFromCurrentTask applies the user's checkbox selection to remove tags.
func (m *Model) RemoveTagsFromCurrentTask() {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}

	targetID := tasks[m.SelectedIndex].ID
	idx := m.findTaskIndexByID(targetID)
	if idx == -1 {
		return
	}

	// Rebuild the tag list, keeping only those NOT marked for removal.
	// We use a "filter-copy" approach here because we are comparing against
	// the separate m.RemoveTagChecks boolean slice.
	var newTags []string
	for j, tag := range m.Tasks[idx].Tags {
		// Ensure we don't index out of bounds if state got out of sync
		if j < len(m.RemoveTagChecks) && !m.RemoveTagChecks[j] {
			newTags = append(newTags, tag)
		}
	}
	m.Tasks[idx].Tags = newTags
}

// ============================================================================
// History & Undo
// ============================================================================

// SaveStateForUndo creates a snapshot of the current task list.
// NOTE: This uses slices.Clone() which creates a shallow copy of the slice structure.
// If the Task struct contains pointers or slices (like Tags), those internal
// structures are shared. Ideally, for a robust Undo, a deep copy is preferred,
// but for simple usage, this isolates the task list structure.
func (m *Model) SaveStateForUndo() {
	stateCopy := slices.Clone(m.Tasks)
	m.History = append(m.History, stateCopy)

	// Rolling buffer: Remove oldest history if we exceed the limit
	if len(m.History) > m.MaxHistory {
		m.History = m.History[1:]
	}
}

// Undo reverts the application state to the last snapshot.
func (m *Model) Undo() {
	if len(m.History) == 0 {
		m.ErrorMessage = "Nothing to undo"
		return
	}

	// Pop the last state
	m.Tasks = m.History[len(m.History)-1]
	m.History = m.History[:len(m.History)-1]

	// Refresh derived state (contexts) and UI cursor
	m.UpdateContexts()
	m.SelectedIndex = 0
}
