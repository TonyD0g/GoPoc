package Format

type PocStruct struct {
	RequestPackage  RequestPackage  `json:"Request"`
	ResponsePackage ResponsePackage `json:"Response"`
	Fofa            string          `json:"fofa"`
}
