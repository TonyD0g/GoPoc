package Handle

import (
	"GoPoc/main/Developer/Format"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

func ProcessXML(inputXml string) Format.RequestPackage {
	xmlData, err := ioutil.ReadFile(inputXml)
	if err != nil {
		fmt.Println("[-] Error reading file:", err)
		os.Exit(1)
	}
	var requestPackage Format.RequestPackage
	err = xml.Unmarshal(xmlData, &requestPackage)
	if err != nil {
		fmt.Println("[-] Error unmarshalling XML:", err)
		os.Exit(1)
	}
	return requestPackage
}
