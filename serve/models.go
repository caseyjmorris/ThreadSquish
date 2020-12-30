package serve

type MenuOption struct {
	Default string
	Description string
	Cases map[string]string
}

type FormFields struct {
	Name string
	FormatName string
	Format string
	Example string
	Description string
	Options []MenuOption
}
