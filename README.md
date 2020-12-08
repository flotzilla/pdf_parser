Pdf metadata parser
====
Go library for parsing pdf metadata information

###Description
This library can:
 * parse pdf info from multiple xref tables
 * obtain metadata
 * extract pdf cover (if cover is an image)

###Usage
```go
import "github.com/flotzilla/pdf_parser.pdf"

// parse file
pdf, errors := pdf_parser.ParsePdf("filepath/file.pdf")

// main functions
pdf.GetTitle()
pdf.GetAuthor()
pdf.GetCreator()
pdf.GetISBN()
pdf.GetPublishers() []string
pdf.GetLanguages() []string
pdf.GetDescription()
pdf.GetPagesCount()
```
