package Format

type ResponsePackage struct {
	Operation string  `json:"Operation"`
	Group     []Group `json:"Group"`
}

type Group struct {
	Regexp string                 `json:"Regexp"`
	Header map[string]interface{} `json:"Header"`
	Body   []string               `json:"Body"`
}
