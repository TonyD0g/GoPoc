package Format

type Fingerprint struct {
	Product string `json:"product"`
	Rule    string `json:"rule"`
}

type Rule struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}
