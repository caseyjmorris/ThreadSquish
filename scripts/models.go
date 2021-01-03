package scripts

type MenuOption struct {
	Default     string            `json:"default"`
	Description string            `json:"description"`
	Cases       map[string]string `json:"cases"`
}

type FormFields struct {
	Name        string       `json:"name"`
	FormatName  string       `json:"formatName"`
	Format      string       `json:"format"`
	Example     string       `json:"example"`
	Description string       `json:"description"`
	Options     []MenuOption `json:"options"`
}

type CommandRequest struct {
	Script    string   `json:"script"`
	Directory string   `json:"directory"`
	Arguments []string `json:"arguments"`
}
