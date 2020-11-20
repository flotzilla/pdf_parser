package pdf_parser

import "errors"

const BufferSize = 50
const BufferSize300 = 300

var (
	fileIsNotPdfError       = errors.New("file is not pdf")
	cannotReadXrefOffset    = errors.New("cannot read OriginalXrefOffset")
	cannotParseXrefOffset   = errors.New("cannot parse OriginalXrefOffset")
	cannotParseXrefSection  = errors.New("cannot parse XrefSection")
	cannotFindObjectById    = errors.New("cannot find object in Xref table")
	cannotParseTrailer      = errors.New("cannot parse trailer section")
	cannotParseObject       = errors.New("cannot parse xref Object")
	unsupportedParseContent = errors.New("unsupported stream decode content")
	cannotFindStreamContent = errors.New("cannot find stream content")
)

type PdfInfo struct {
	PdfVersion               string
	OriginalXrefOffset       int64
	OriginalTrailerSection   TrailerSection
	AdditionalTrailerSection []*TrailerSection
	XrefTable                []*XrefTable
	Root                     RootObject
	Info                     InfoObject
	Metadata                 Metadata
	PagesCount               int
}

func (pdf *PdfInfo) getTitle() string {
	if pdf.Info.Title != "" {
		return pdf.Info.Title
	}
	return pdf.Metadata.RdfMeta.Title
}

func (pdf *PdfInfo) getAuthor() string {
	if pdf.Info.Author != "" {
		return pdf.Info.Author
	}
	return pdf.Metadata.RdfMeta.Creator
}

func getCreator(pdf *PdfInfo) string {
	if pdf.Info.Creator != "" {
		return pdf.Info.Creator
	}
	return ""
}

func getISBN(pdf *PdfInfo) string {
	return pdf.Metadata.RdfMeta.Isbn
}

func getPublisher(pdf *PdfInfo) []string {
	return pdf.Metadata.RdfMeta.Publishers
}

func getLanguages(pdf *PdfInfo) []string {
	return pdf.Metadata.RdfMeta.Languages
}

func getDescription(pdf *PdfInfo) string {
	return pdf.Metadata.RdfMeta.Description
}

type TrailerSection struct {
	IdRaw string
	Info  ObjectIdentifier
	Root  ObjectIdentifier
	Size  string
	Prev  int64
}

type ObjectIdentifier struct {
	ObjectNumber     int
	GenerationNumber int
	KeyWord          string
}

type ObjectSubsectionElement struct {
	Id               int
	ObjectNumber     int
	GenerationNumber int
	KeyWord          string
}

/*
	Object subsection that contain list of objects for this object
*/
type ObjectSubsection struct {
	Id                      int // objectId
	ObjectsCount            int
	FirstSubsectionObjectId int
	LastSubsectionObjectId  int
	Elements                map[int]*ObjectSubsectionElement
}

type XrefTable struct {
	Objects           map[int]*ObjectSubsectionElement
	ObjectSubsections map[int]*ObjectSubsection
	SectionStart      int64
}

type InfoObject struct {
	Title        string
	Author       string
	Creator      string
	CreationDate string
	Producer     string
	ModDate      string
}

type RootObject struct {
	Type       string
	Pages      *ObjectIdentifier
	Metadata   *ObjectIdentifier
	PageLabels *ObjectIdentifier
	Lang       string
}

type Metadata struct {
	Type          string
	Subtype       string
	Length        int64
	DL            int64
	RawStreamData []byte
	RdfMeta       *MetaDataRdf
}

type MetaDataRdf struct {
	Title       string
	Description string
	Creator     string
	Date        string
	Isbn        string

	Publishers []string
	Languages  []string
}
