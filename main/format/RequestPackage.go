package format

type RequestPackage struct {
	Method   string     `json:"Method"`
	Url      string     `json:"Url"`
	PathList PathList   `json:"PathList"`
	Header   HeaderList `json:"HeaderList"`
	Body     string     `json:"Body"`
}

type PathList struct {
	Path string `json:"Path"`
}
type HeaderList struct {
	UserAgent      string `json:"User-Agent"`
	AcceptEncoding string `json:"Accept-Encoding"`
	Accept         string `json:"Accept"`
	Cookie         string `json:"Cookie"`
}
