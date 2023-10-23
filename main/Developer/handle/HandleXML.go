package handle

import (
	"Scanner/main/Developer/format"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

func ProcessXML(inputXml string) format.RequestPackage {
	xmlData, err := ioutil.ReadFile(inputXml)
	if err != nil {
		fmt.Println("[-] Error reading file:", err)
		os.Exit(1)
	}
	var requestPackage format.RequestPackage
	err = xml.Unmarshal(xmlData, &requestPackage)
	if err != nil {
		fmt.Println("[-] Error unmarshalling XML:", err)
		os.Exit(1)
	}
	return requestPackage
}
