package Format

type PocStruct struct {
	RequestPackage  RequestPackage  `json:"Request"`
	ResponsePackage ResponsePackage `json:"Response"`
	Fofa            string          `json:"Fofa"`
	Uri             string          `json:"Uri"`
	Url             string          `json:"Url"`
	File            string          `json:"File"`
	Coroutine       string          `json:"Coroutine"`
	CheckIP         string          `json:"CheckIP"`
	VulnName        string          `json:"VulnName"`
}
