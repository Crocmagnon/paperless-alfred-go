package paperless

import (
	"fmt"
	"time"
)

type SearchResponse struct {
	PaginationEnvelope
	Results []Result `json:"results"`
}

type Result struct {
	Id                  int           `json:"id"`
	Correspondent       *int          `json:"correspondent"`
	DocumentType        *int          `json:"document_type"`
	StoragePath         interface{}   `json:"storage_path"`
	Title               string        `json:"title"`
	Content             string        `json:"content"`
	Tags                []int         `json:"tags"`
	Created             time.Time     `json:"created"`
	CreatedDate         string        `json:"created_date"`
	Modified            time.Time     `json:"modified"`
	Added               time.Time     `json:"added"`
	DeletedAt           time.Time     `json:"deleted_at"`
	ArchiveSerialNumber *int          `json:"archive_serial_number"`
	OriginalFileName    string        `json:"original_file_name"`
	ArchivedFileName    string        `json:"archived_file_name"`
	Owner               int           `json:"owner"`
	UserCanChange       bool          `json:"user_can_change"`
	IsSharedByRequester bool          `json:"is_shared_by_requester"`
	Notes               []interface{} `json:"notes"`
	CustomFields        []interface{} `json:"custom_fields"`
	PageCount           int           `json:"page_count"`
	SearchHit           SearchHit     `json:"__search_hit__"`
}

func (r Result) Metadata(correspondents map[int]Correspondent, docTypes map[int]DocumentType) []string {
	meta := make([]string, 0, 4)

	meta = append(meta, r.CreatedDate)

	if docType := r.DocumentTypeName(docTypes); docType != "" {
		meta = append(meta, docType)
	}

	if corr := r.CorrespondentName(correspondents); corr != "" {
		meta = append(meta, corr)
	}

	if asn := r.ASN(); asn != "" {
		meta = append(meta, asn)
	}

	return meta
}

func (r Result) ASN() string {
	if r.ArchiveSerialNumber == nil {
		return ""
	}

	return fmt.Sprintf("ASN %d", *r.ArchiveSerialNumber)
}

func (r Result) DetailsURL(baseURL string) string {
	return fmt.Sprintf("%s/documents/%v/details", baseURL, r.Id)
}

func (r Result) CorrespondentName(correspondents map[int]Correspondent) string {
	if r.Correspondent == nil {
		return ""
	}

	if corr, ok := correspondents[*r.Correspondent]; ok {
		return corr.Name
	}

	return ""
}

func (r Result) DocumentTypeName(documentTypes map[int]DocumentType) string {
	if r.DocumentType == nil {
		return ""
	}

	if docType, ok := documentTypes[*r.DocumentType]; ok {
		return docType.Name
	}

	return ""
}

type SearchHit struct {
	Score          float64 `json:"score"`
	Highlights     string  `json:"highlights"`
	NoteHighlights string  `json:"note_highlights"`
	Rank           int     `json:"rank"`
}

type PaginationEnvelope struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	All      []int   `json:"all"`
}

type CorrespondentsResponse struct {
	PaginationEnvelope
	Results []Correspondent `json:"results"`
}

type Correspondent struct {
	Id                int    `json:"id"`
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	Match             string `json:"match"`
	MatchingAlgorithm int    `json:"matching_algorithm"`
	IsInsensitive     bool   `json:"is_insensitive"`
	DocumentCount     int    `json:"document_count"`
	Owner             int    `json:"owner"`
	UserCanChange     bool   `json:"user_can_change"`
}

type DocumentTypesResponse struct {
	PaginationEnvelope
	Results []DocumentType `json:"results"`
}

type DocumentType struct {
	Id                int    `json:"id"`
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	Match             string `json:"match"`
	MatchingAlgorithm int    `json:"matching_algorithm"`
	IsInsensitive     bool   `json:"is_insensitive"`
	DocumentCount     int    `json:"document_count"`
	Owner             int    `json:"owner"`
	UserCanChange     bool   `json:"user_can_change"`
}
