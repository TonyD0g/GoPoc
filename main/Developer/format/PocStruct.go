package Format

type PocStruct struct {
	RequestPackage  RequestPackage  `json:"RequestPackage"`
	ResponsePackage ResponsePackage `json:"ResponsePackage"`
	Fofa            string          `json:"fofa"`
}
