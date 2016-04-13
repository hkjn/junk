package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Staff struct {
	XMLName   xml.Name `xml:"staff"`
	ID        int      `xml:"id"`
	FirstName string   `xml:"firstname"`
	LastName  string   `xml:"lastname"`
	UserName  string   `xml:"username"`
}

type Company struct {
	XMLName xml.Name `xml:"company"`
	Staffs  []Staff  `xml:"staff"`
}

func (s Staff) String() string {
	return fmt.Sprintf("\t ID : %d - FirstName : %s - LastName : %s - UserName : %s \n", s.ID, s.FirstName, s.LastName, s.UserName)
}

type (
	Scan struct {
		XMLName xml.Name `xml:"document"`
		Tests   []struct {
			Host       string `xml:"host,attr"`
			Port       int    `xml:"port,attr"`
			Heartbleed []struct {
				SSLVersion string `xml:"sslversion,attr"`
				Vulnerable int    `xml:"vulnerable,attr"`
			} `xml:"heartbleed"`
			Cipher []struct {
				Status     string `xml:"status,attr"`
				SSLVersion string `xml:"sslversion,attr"`
				Bits       int    `xml:"bits,attr"`
				Cipher     string `xml:"cipher,attr"`
			} `xml:"cipher"`
			DefaultCipher []struct {
				SSLVersion string `xml:"sslversion,attr"`
				Bits       int    `xml:"bits,attr"`
				Cipher     string `xml:"cipher,attr"`
				DHEBits    int    `xml:"dhebits,attr"`
			} `xml:"defaultcipher"`
			Certificate []struct {
				SignAlgorithm struct {
					Algorithm string `xml:",chardata"`
				} `xml:"signature-algorithm"`
				PK struct {
					Error bool   `xml:"error,attr"`
					Type  string `xml:"type,attr"`
					Bits  int    `xml:"bits,attr"`
				} `xml:"pk"`
				Subject    string `xml:"subject"`
				AltNames   string `xml:"altnames"`
				Issuer     string `xml:"issuer"`
				SelfSigned bool   `xml:"self-signed"`
				FirstValid string `xml:"not-valid-before"`
				LastValid  string `xml:"not-valid-after"`
			} `xml:"certificate"`
		} `xml:"ssltest"`
	}
)

func main() {
	xmlFile, err := os.Open("emp.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var c Company
	//var d struct{}
	xml.Unmarshal(XMLdata, &c)
	//	fmt.Printf("%v\n", c)
	//	fmt.Println(c.Staffs)

	xmlFile, err = os.Open("scan.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)

	var s Scan
	//var d struct{}
	xml.Unmarshal(data, &s)
	fmt.Printf("%+v\n", s)
}
