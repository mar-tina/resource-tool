package schema

type CollectionPayload struct {
	Collection *Collection `json:"collection"`
}

type EnvPayload struct {
	Environment *Environment `json:"environment"`
}

type Collection struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Owner       string      `json:"owner,omitempty"`
	UID         string      `json:"uid,omitempty"`
	Descendants []Service   `json:"descendants,omitempty"`
	Ancestors   []Service   `json:"ancestors,omitempty"`
	Resources   []Resource  `json:"resources,omitempty"`
	Info        Info        `json:"info,omitempty"`
	Item        []ItemEntry `json:"item,omitempty"`
}

type Service struct {
	Name      string
	Resources []string
}

type Resource struct {
	Name     string
	Params   string
	Defaults map[string]interface{}
	Body     map[string]interface{}
}

type ResourceCheckPayload struct {
	Collection string `json:"collection"`
	Resource   string `json:"resource"`
}

type Info struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema,omitempty"`
}

type ItemEntry struct {
	Item []NestedItem `json:"item,omitempty"`
	Name string       `json:"name,omitempty"`
}

type NestedItem struct {
	Name    string  `json:"name,omitempty"`
	Request Request `json:"request,omitempty"`
}

type Request struct {
	URL    string   `json:"url,omitempty"`
	Method string   `json:"method,omitempty"`
	Header []Header `json:"header,omitempty"`
	Body   Body     `json:"body,omitempty"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Body struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type Environment struct {
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Values []*EnvValues `json:"values"`
}

type EnvValues struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
}

type FetchResoure struct {
	ServiceID string `json:"service_id"`
	Key       string `json:"key"`
	Resource  string `json:"resource"`
}

type FetchEnv struct {
	ID string `json:"id"`
}

type FetchAllResult struct {
	Environments []*Env `json:"environments"`
}

type Env struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
	UID   string `json:"uid"`
}

type UpdateEnv struct {
	Env Environment `json:"environment"`
}

type FetchAllCollections struct {
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}
