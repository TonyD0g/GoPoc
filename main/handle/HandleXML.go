package handle

import (
	"Scanner/main/format"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func HandleXMLFunc(inputXml string) format.RequestPackage {
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

func CheckFileCorrectness(requestPackage format.RequestPackage) {
	if requestPackage.Method != "GET" && requestPackage.Url != "POST" {
		fmt.Print("[-] error! You must provide the correct method")
		flag.Usage()
		os.Exit(1)
	}
}
