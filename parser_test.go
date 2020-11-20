package pdf_parser

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"
)

func TestParseBagElementsSingle(t *testing.T) {
	reg := regexp.MustCompile(`(?m)\<dc:publisher>((.|\n)*)\<rdf:Bag>((.|\n)*)\<\/rdf:Bag\>((.|\n)*)\<\/dc:publisher>`)
	data := []byte("<dc:publisher>\n        <rdf:Bag>\n          <rdf:li>No Starch Press</rdf:li>\n        </rdf:Bag>\n      </dc:publisher>")
	resp := parseBagElements(reg, &data)

	if len(resp) != 1 {
		t.Error("Expected to be string array ", resp)
	}
}

func TestParseBagElements(t *testing.T) {
	reg := regexp.MustCompile(`(?m)\<dc:publisher>((.|\n)*)\<rdf:Bag>((.|\n)*)\<\/rdf:Bag\>((.|\n)*)\<\/dc:publisher>`)
	data := []byte("<dc:publisher>\n        <rdf:Bag>\n          <rdf:li>No Starch Press</rdf:li>\n " +
		"<rdf:li>No Starch Press</rdf:li>\n " +
		"       </rdf:Bag>\n      </dc:publisher>")
	resp := parseBagElements(reg, &data)

	if len(resp) != 2 {
		t.Error("Expected to be string array with size 2", resp)
	}
}

func TestParseBagElementsZero(t *testing.T) {
	reg := regexp.MustCompile(`(?m)\<dc:publisher>((.|\n)*)\<rdf:Bag>((.|\n)*)\<\/rdf:Bag\>((.|\n)*)\<\/dc:publisher>`)
	data := []byte("<dc:publisher>\n    </dc:publisher>")
	resp := parseBagElements(reg, &data)

	if len(resp) != 0 {
		t.Error("Expected to be empty array", resp)
	}
}

func TestParsePdf(t *testing.T) {
	path, err := filepath.Abs("../resources/sample.pdf")
	if err != nil {
		t.Error(err)
	}
	file := ParsePdf(path)
	fmt.Print(file)
}