package Format

type RequestPackage struct {
	Method string                 `json:"Method"`
	Url    string                 `json:"Url"`
	Uri    []string               `json:"Uri"`
	Header map[string]interface{} `json:"Header"`
	Body   string                 `json:"Body"`
}
