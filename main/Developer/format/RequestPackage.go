package Format

type RequestPackage struct {
	Method string   `json:"Method"`
	Url    string   `json:"Url"`
	Uri    []string `json:"Uri"`
	Header Header   `json:"Header"`
	Body   string   `json:"Body"`
}
