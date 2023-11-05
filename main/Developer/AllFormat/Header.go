package Format

type Header struct {
	UserAgent      string `json:"User-Agent"`
	AcceptEncoding string `json:"Accept-Encoding"`
	Accept         string `json:"Accept"`
	Host           string `json:"Host"`
	ContentType    string `json:"Content-Type"`
}
