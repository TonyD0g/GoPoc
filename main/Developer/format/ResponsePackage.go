package Format

type ResponsePackage struct {
	Operation string  `json:"Operation"`
	Group     []Group `json:"Group"`
}

type Group struct {
	Header Header   `json:"Header"`
	Body   []string `json:"Body"`
}
