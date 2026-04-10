package cmd

type Root struct {
	JSON bool `help:"Output JSON."`

	Auth         AuthCmd         `cmd:"" help:"Manage authentication."`
	Config       ConfigCmd       `cmd:"" help:"Manage configuration."`
	CustomFields CustomFieldsCmd `cmd:"" help:"Manage custom fields."`
	Projects     ProjectsCmd     `cmd:"" help:"Manage projects."`
	Sections     SectionsCmd     `cmd:"" help:"Manage sections."`
	Tasks        TasksCmd        `cmd:"" help:"Manage tasks."`
	Users        UsersCmd        `cmd:"" help:"Manage users."`
	Workspaces   WorkspacesCmd   `cmd:"" help:"Manage workspaces."`
}
