package pdf_parser

import (
	logger "github.com/sirupsen/logrus"
	//"fmt"
	"os"
	"path/filepath"

	//"path/filepath"
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
	path, err := filepath.Abs("./resources/test.pdf")
	if err != nil {
		t.Error(err)
	}
	file, err := ParsePdf(path)

	if err != nil {
		t.Error(err)
	}

	if file == nil {
		t.Error("general fail send regards")
	}
}

func TestPdfPageCount(t *testing.T) {
	file, err := os.Open("./resources/sample.pdf")
	if err != nil {
		t.Error(err)
	}
	if count := countPages(file) != 2; count {
		t.Error("Wrong page count, should be 2 ")
	}
}

func TestParseTestPdfInfo(t *testing.T) {
	file, err := filepath.Abs("./resources/test.pdf")
	if err != nil {
		t.Error(err)
	}

	pdfInf, err := ParsePdf(file)
	if err != nil {
		t.Error(err)
	}

	if pdfInf.PdfVersion != "2.0" {
		t.Error("Invalid pdf version parsing")
	}

	if pdfInf.PagesCount != 1 {
		t.Error("Invalid pages count parsing")
	}

	testName := "Soda PDF Online"
	if pdfInf.Info.Author != testName || pdfInf.Info.Creator != testName ||
		pdfInf.Info.Producer != testName {
		t.Error("Invalid info block parsing")
	}

	if pdfInf.OriginalXrefOffset != 16062 {
		t.Error("Invalid original offset value parse")
	}

	if pdfInf.OriginalTrailerSection.Info.ObjectNumber != 11 ||
		pdfInf.OriginalTrailerSection.Root.ObjectNumber != 1 ||
		pdfInf.OriginalTrailerSection.Prev != 0 ||
		pdfInf.OriginalTrailerSection.IdRaw != "[<DB0B127ADB050BC7FA6212CE9DDBFC70><246EA46DE075715E82A47F70DB5527D6>]" {
		t.Error("Invalid original section parsing")
	}

	if len(pdfInf.XrefTable) != 1 {
		t.Error("Invalid Xreftable parse")
	}

	if pdfInf.XrefTable[0].SectionStart != 16062 ||
		len(pdfInf.XrefTable[0].ObjectSubsections) != 2 ||
		pdfInf.XrefTable[0].ObjectSubsections[0].ObjectsCount != 6 ||
		pdfInf.XrefTable[0].ObjectSubsections[11].ObjectsCount != 7 ||
		pdfInf.XrefTable[0].ObjectSubsections[11].Id != 11 ||
		pdfInf.XrefTable[0].ObjectSubsections[11].LastSubsectionObjectId != 17 {
		t.Error("Invalid xref object parsing")
	}

	if pdfInf.Root.Type != "/Catalog" ||
		pdfInf.Root.PageLabels != nil ||
		pdfInf.Root.Lang != "" {
		t.Error("Invalid root section parse")
	}

	if pdfInf.Root.Metadata.ObjectNumber != 12 {
		t.Error("Invalid root section metadata link parse")
	}

	if pdfInf.Root.Pages.ObjectNumber != 2 {
		t.Error("Invalid root section pages link parse")
	}

	if pdfInf.Metadata.Subtype != "XML" ||
		pdfInf.Metadata.Length != 3175 ||
		pdfInf.Metadata.Type != "Metadata" ||
		pdfInf.Metadata.RdfMeta.Creator != testName {
		t.Error("Invalid metadata section parse")
	}
}

func TestGetTitle(t *testing.T) {
	file, _ := filepath.Abs("./resources/test.pdf")
	pdf, err := ParsePdf(file)

	if err != nil {
		t.Error(err)
	}

	if pdf.GetTitle() != pdf.Info.Title ||
		pdf.GetTitle() != pdf.Metadata.RdfMeta.Title {
		t.Error("Invalid GetTitle behaviour")
	}
}

func TestGetAuthor(t *testing.T) {
	file, _ := filepath.Abs("./resources/test.pdf")
	pdf, err := ParsePdf(file)

	if err != nil {
		t.Error(err)
	}

	if pdf.GetAuthor() != pdf.Info.Author ||
		pdf.GetAuthor() != pdf.Metadata.RdfMeta.Creator {
		t.Error("Invalid getAuthor behaviour")
	}
}

func TestSetLogger(t *testing.T) {
	lg := logger.New()
	lg.SetOutput(os.Stdout)
	lg.SetReportCaller(true)
	lg.SetFormatter(&logger.JSONFormatter{})

	SetLogger(lg)
	file, _ := filepath.Abs("./resources/test.pdf")
	pdf, err := ParsePdf(file)

	if err != nil {
		t.Error(err)
	}

	if pdf.GetPagesCount() != 1 {
		t.Error("Invalid amount of pages")
	}
}

func TestNonExistedFileWithLof(t *testing.T) {
	file, _ := filepath.Abs("./resources/test_non_existed.pdf")
	_, err := ParsePdf(file)

	if err == nil {
		t.Error(err)
	}
}

func TestWriteLogToFile(t *testing.T) {
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		t.Error(err)
	}

	defer f.Close()

	lg := logger.New()
	lg.SetOutput(f)
	lg.SetFormatter(&logger.JSONFormatter{})

	SetLogger(lg)
	file, _ := filepath.Abs("./resources/test.pdf")
	pdf, err := ParsePdf(file)

	if err != nil {
		t.Error(err)
	}

	if pdf.GetPagesCount() != 1 {
		t.Error("Invalid amount of pages")
	}

	logFile, errF := os.Stat("test.log")
	if errF != nil {
		t.Error(errF)
	}

	if logFile.Size() == 0 {
		t.Error("File should not be empty")
	}

	if os.IsNotExist(errF) {
		t.Error(errF)
	}

	e := os.Remove("test.log")
	if e != nil {
		t.Error(err)
	}
}
