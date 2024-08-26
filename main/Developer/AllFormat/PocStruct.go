package Format

type PocStruct struct {
	RequestPackage  RequestPackage  `json:"Request"`
	ResponsePackage ResponsePackage `json:"Response"`
	Fofa            string          `json:"Fofa"`
	Uri             string          `json:"Uri"`
	Url             string          `json:"Url"`
	File            string          `json:"File"`
}
