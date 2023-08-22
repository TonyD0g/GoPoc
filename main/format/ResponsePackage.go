package format

type ResponsePackage struct {
	Operation string  `json:"Operation"`
	Group     []Group `json:"Group"`
}

type Group struct {
	Status  string  `json:"Status"`
	Contain Contain `json:"Contain"`
}

type Contain struct {
	Header Header   `json:"Header"`
	Body   []string `json:"Body"`
}
