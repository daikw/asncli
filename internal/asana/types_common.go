package asana

type User struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type Workspace struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type Page struct {
	Offset string `json:"offset"`
	URI    string `json:"uri"`
}
