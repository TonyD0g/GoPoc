package format

type Header struct {
	UserAgent      string `json:"User-Agent"`
	AcceptEncoding string `json:"Accept-Encoding"`
	Accept         string `json:"Accept"`
	Cookie         string `json:"Cookie"`
	Host           string `json:"Host"`
	ContentType    string `json:"Content-Type"`
}
