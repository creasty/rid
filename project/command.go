package project

// Command represents a custom sub-command
type Command struct {
	Name           string
	Summary        string
	Description    string
	RunInContainer bool
	HelpFile       string
}
