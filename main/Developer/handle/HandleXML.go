package Handle

import (
	"GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Log"
	"encoding/xml"
	"io/ioutil"
)

func ProcessXML(inputXml string) Format.RequestPackage {
	xmlData, err := ioutil.ReadFile(inputXml)
	if err != nil {
		Log.Log.Fatal("[-] Error reading file:", err)
	}
	var requestPackage Format.RequestPackage
	err = xml.Unmarshal(xmlData, &requestPackage)
	if err != nil {
		Log.Log.Fatal("[-] Error unmarshalling XML:", err)
	}
	return requestPackage
}
