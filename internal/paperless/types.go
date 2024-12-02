package paperless

import (
	"fmt"
	"time"
)

type SearchResponse struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	All      []int       `json:"all"`
	Results  []Result    `json:"results"`
}

type Result struct {
	Id                  int           `json:"id"`
	Correspondent       *int          `json:"correspondent"`
	DocumentType        int           `json:"document_type"`
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

func (r Result) Metadata(correspondents map[int]Correspondent) []string {
	meta := make([]string, 0, 3)

	if asn := r.ASN(); asn != "" {
		meta = append(meta, asn)
	}

	if corr := r.CorrespondentName(correspondents); corr != "" {
		meta = append(meta, corr)
	}

	meta = append(meta, r.CreatedDate)

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

type SearchHit struct {
	Score          float64 `json:"score"`
	Highlights     string  `json:"highlights"`
	NoteHighlights string  `json:"note_highlights"`
	Rank           int     `json:"rank"`
}

type CorrespondentsResponse struct {
	Count    int             `json:"count"`
	Next     *string         `json:"next"`
	Previous *string         `json:"previous"`
	All      []int           `json:"all"`
	Results  []Correspondent `json:"results"`
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
