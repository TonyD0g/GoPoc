package format

type RequestPackage struct {
	Method string `xml:"method"`
	Url    string `xml:"url"`
	Header string `xml:"header"`
}

// 示例xml
// <data>
//    <person>
//        <name>John Doe</name>
//        <age>30</age>
//        <country>USA</country>
//    </person>
//    <person>
//        <name>Jane Smith</name>
//        <age>25</age>
//        <country>Canada</country>
//    </person>
//</data>
