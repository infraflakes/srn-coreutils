package todo

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
)

// Task represents a single todo item
type Task struct {
	ID       int      `json:"id"`
	Task     string   `json:"task"`
	Checked  bool     `json:"checked"`
	Context  string   `json:"context"`
	Priority string   `json:"priority,omitempty"` // low, medium, high
	Tags     []string `json:"tags,omitempty"`
	DueDate  string   `json:"due_date,omitempty"` // YYYY-MM-DD format
}

// ViewMode represents the current view
type ViewMode int

const (
	NormalView ViewMode = iota
	KanbanView
	StatsView
	InputView
	DateInputView
	RemoveTagView
)

// InputMode represents different input dialogs
type InputMode int

const (
	AddTaskInput InputMode = iota
	EditTaskInput
	AddContextInput
	RenameContextInput
	AddTagInput
	DeleteConfirmInput
)

// Model represents the entire state of the todo application.
// It's the central data structure that gets passed around and updated
// by the Bubble Tea runtime.
type Model struct {
	// --- Core State ---
	// These fields represent the fundamental data of the application.
	Tasks          []Task   // A slice holding all tasks across all contexts. This is the single source of truth.
	Contexts       []string // A list of all available contexts (e.g., "Work", "Personal").
	CurrentContext string   // The context currently being viewed by the user.
	SelectedIndex  int      // The index of the currently selected task in the *filtered* view.
	NextID         int      // A counter to ensure all new tasks get a unique ID.

	// --- View State ---
	// These fields control what is currently being displayed on the screen.
	ViewMode  ViewMode // Determines which major view is active (e.g., Normal, Kanban, Input).
	InputMode InputMode // If ViewMode is InputView, this specifies the type of input (e.g., adding vs. editing a task).

	// State for task movement.
	MovingMode   bool // Flag to indicate if the user has initiated a task move.
	MovingTaskID int  // The unique ID of the task being moved. Using the ID is crucial as the task's index can change.

	// --- Input Handling ---
	// These fields manage the state of various user input components.
	TextInput       textinput.Model // A text input component for adding/editing tasks, contexts, etc.
	DateInputs      []textinput.Model
	DateInputIndex  int
	RemoveTagIndex  int
	RemoveTagChecks []bool
	InputPrompt     string // The prompt to display when in InputView (e.g., "Add new task:").

	// --- UI State ---
	// General UI-related state.
	WindowWidth  int
	WindowHeight int
	ErrorMessage string // A message to display to the user when an error occurs.

	// --- History for Undo ---
	History    [][]Task // A stack of previous task states to allow for undo operations.
	MaxHistory int      // The maximum number of undo states to store.

	// --- Keybindings & Help ---
	KeyMap KeyMap     // Holds the application's key bindings.
	Help   help.Model // The help bubble component.

	// --- Configuration ---
	ConfigPath string // The file path where the application's state is saved.
}

// KeyMap defines key bindings
type KeyMap struct {
	Up             key.Binding
	Down           key.Binding
	Left           key.Binding
	Right          key.Binding
	Toggle         key.Binding
	Add            key.Binding
	Edit           key.Binding
	Delete         key.Binding
	AddContext     key.Binding
	RenameContext  key.Binding
	DeleteContext  key.Binding
	TogglePriority key.Binding
	AddTag         key.Binding
	RemoveTag      key.Binding
	SetDueDate     key.Binding
	ClearDueDate   key.Binding
	KanbanView     key.Binding
	StatsView      key.Binding
	Undo           key.Binding
	Move           key.Binding
	Quit           key.Binding
	Back           key.Binding
	Enter          key.Binding
	Nav            key.Binding
}

// DefaultKeyMap returns default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "prev context"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next context"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add task"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		AddContext: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new context"),
		),
		RenameContext: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename context"),
		),
		DeleteContext: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "delete context"),
		),
		TogglePriority: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "priority"),
		),
		AddTag: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "add tag"),
		),
		RemoveTag: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "remove tag"),
		),
		SetDueDate: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "due date"),
		),
		ClearDueDate: key.NewBinding(
			key.WithKeys("U"),
			key.WithHelp("U", "clear due"),
		),
		KanbanView: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "kanban"),
		),
		StatsView: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "stats"),
		),
		Undo: key.NewBinding(
			key.WithKeys("z"),
			key.WithHelp("z", "undo"),
		),
		Move: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "move"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Nav: key.NewBinding(
			key.WithKeys("↑", "↓", "←", "→"),
			key.WithHelp("↑↓←→", "navigation"),
		),
	}
}

// KeyMap methods to implement help.KeyMap interface
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Nav, k.Toggle, k.Add, k.Edit, k.Delete, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Nav},
		{k.Toggle, k.Add, k.Edit, k.Delete, k.Move},
		{k.AddContext, k.RenameContext, k.DeleteContext},
		{k.TogglePriority, k.AddTag, k.RemoveTag, k.SetDueDate, k.ClearDueDate},
		{k.KanbanView, k.StatsView},
		{k.Undo, k.Back, k.Quit},
	}
}
