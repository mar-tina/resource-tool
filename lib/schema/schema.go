package schema

type Collection struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Owner       string      `json:"owner"`
	UID         string      `json:"uid"`
	Descendants []Service   `json:"descendants"`
	Ancestors   []Service   `json:"ancestors"`
	Resources   []Resource  `json:"resources"`
	Info        Info        `json:"info"`
	Item        []ItemEntry `json:"item"`
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
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
}

type ItemEntry struct {
	Item []NestedItem `json:"item"`
	Name string       `json:"name"`
}

type NestedItem struct {
	Name    string  `json:"name"`
	Request Request `json:"request"`
}

type Request struct {
	URL    string   `json:"url"`
	Method string   `json:"method"`
	Header []Header `json:"header"`
	Body   Body     `json:"body"`
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
}

type FetchResoure struct {
	ServiceID string `json:"service_id"`
	Key       string `json:"key"`
	Resource  string `json:"resource"`
}

type UpdateEnv struct {
	Env Environment `json:"environment"`
}
