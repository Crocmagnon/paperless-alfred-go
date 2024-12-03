package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Crocmagnon/paperless-alfred-go/internal/alfred"
	"github.com/Crocmagnon/paperless-alfred-go/internal/paperless"
	"github.com/jarcoal/httpmock"
	"gotest.tools/v3/assert"
	"net/http"
	"testing"
)

func Test_run(t *testing.T) {
	const fakePaperlessHost = "http://localhost:1234"
	t.Parallel()
	type args struct {
		args       []string
		httpClient func() *http.Client
	}
	tests := []struct {
		name       string
		args       args
		wantStdout alfred.Result
		wantErr    bool
	}{
		{
			"simple query",
			args{
				[]string{fakePaperlessHost, "1234", "simple"},
				func() *http.Client {
					transport := httpmock.NewMockTransport()
					transport.RegisterResponder(
						http.MethodGet, "/api/documents/?query=simple",
						httpmock.NewJsonResponderOrPanic(http.StatusOK, paperless.PageResponse[paperless.DocumentSearch]{
							Results: []paperless.DocumentSearch{
								{
									Id:                  12,
									Correspondent:       nil,
									DocumentType:        nil,
									Title:               "Simple doc",
									CreatedDate:         "2024-10-12",
									ArchiveSerialNumber: nil,
								},
							},
						}))
					transport.RegisterResponder(http.MethodGet, "/api/correspondents/", httpmock.NewStringResponder(http.StatusOK, "{}"))
					transport.RegisterResponder(http.MethodGet, "/api/document_types/", httpmock.NewStringResponder(http.StatusOK, "{}"))

					return &http.Client{Transport: transport}
				},
			},
			alfred.Result{
				Items: []alfred.Item{
					{
						UID:      "http://localhost:1234/documents/12/details",
						Title:    "Simple doc",
						Subtitle: "2024-10-12",
						Arg:      "http://localhost:1234/documents/12/details",
						Icon:     &alfred.Icon{Type: "filetype", Path: "com.adobe.pdf"},
						Mods: map[string]alfred.Mod{
							"cmd": {Arg: "http://localhost:1234/documents?query=simple", Subtitle: "Open search in Paperless"},
						},
					},
				},
			},
			false,
		},
		{
			"two docs, one with correspondent, doc type and ASN",
			args{
				[]string{fakePaperlessHost, "1234", "complex querÿ"},
				func() *http.Client {
					corID := 1
					docTypeID := 2
					asn := 123
					transport := httpmock.NewMockTransport()
					transport.RegisterResponder(
						http.MethodGet, "/api/documents/?query=complex+quer%C3%BF",
						httpmock.NewJsonResponderOrPanic(http.StatusOK, paperless.PageResponse[paperless.DocumentSearch]{
							Results: []paperless.DocumentSearch{
								{
									Id:                  11,
									Correspondent:       &corID,
									DocumentType:        &docTypeID,
									Title:               "Complete doc",
									CreatedDate:         "2024-10-11",
									ArchiveSerialNumber: &asn,
								},
								{
									Id:                  12,
									Correspondent:       nil,
									DocumentType:        nil,
									Title:               "Simple doc",
									CreatedDate:         "2024-10-12",
									ArchiveSerialNumber: nil,
								},
							},
						}))

					nextPage := fakePaperlessHost + "/api/correspondents/?page=2"
					transport.RegisterResponder(http.MethodGet, "/api/correspondents/",
						httpmock.NewJsonResponderOrPanic(http.StatusOK, paperless.PageResponse[paperless.Correspondent]{
							PaginationEnvelope: paperless.PaginationEnvelope{
								Next: &nextPage,
							},
							Results: nil,
						}),
					)
					transport.RegisterResponder(http.MethodGet, "/api/correspondents/?page=2",
						httpmock.NewJsonResponderOrPanic(http.StatusOK, paperless.PageResponse[paperless.Correspondent]{
							Results: []paperless.Correspondent{
								{Id: 1, Name: "Fake corresp."},
							},
						}),
					)

					transport.RegisterResponder(http.MethodGet, "/api/document_types/",
						httpmock.NewJsonResponderOrPanic(http.StatusOK, paperless.PageResponse[paperless.DocumentType]{
							Results: []paperless.DocumentType{
								{Id: 2, Name: "Doc type"},
							},
						}))

					return &http.Client{Transport: transport}
				},
			},
			alfred.Result{
				Items: []alfred.Item{
					{
						UID:      "http://localhost:1234/documents/11/details",
						Title:    "Complete doc",
						Subtitle: "2024-10-11 - Doc type - Fake corresp. - ASN 123",
						Arg:      "http://localhost:1234/documents/11/details",
						Icon:     &alfred.Icon{Type: "filetype", Path: "com.adobe.pdf"},
						Mods: map[string]alfred.Mod{
							"cmd": {Arg: "http://localhost:1234/documents?query=complex querÿ", Subtitle: "Open search in Paperless"},
						},
					},
					{
						UID:      "http://localhost:1234/documents/12/details",
						Title:    "Simple doc",
						Subtitle: "2024-10-12",
						Arg:      "http://localhost:1234/documents/12/details",
						Icon:     &alfred.Icon{Type: "filetype", Path: "com.adobe.pdf"},
						Mods: map[string]alfred.Mod{
							"cmd": {Arg: "http://localhost:1234/documents?query=complex querÿ", Subtitle: "Open search in Paperless"},
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stdout := &bytes.Buffer{}
			ctx := context.Background()
			err := run(ctx, tt.args.args, stdout, tt.args.httpClient())
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotStdout := alfred.Result{}
			err = json.Unmarshal(stdout.Bytes(), &gotStdout)
			assert.NilError(t, err)

			assert.DeepEqual(t, tt.wantStdout, gotStdout)
		})
	}
}
