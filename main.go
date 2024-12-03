package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Crocmagnon/paperless-alfred-go/internal/alfred"
	"github.com/Crocmagnon/paperless-alfred-go/internal/paperless"
	"golang.org/x/text/unicode/norm"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args[1:], os.Stdout, http.DefaultClient); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout io.Writer, httpClient *http.Client) error {
	if len(args) != 3 {
		return fmt.Errorf("usage: ppl-go <base_url> <token> '<query>'")
	}

	baseURL := norm.NFC.String(args[0])
	token := norm.NFC.String(args[1])
	query := norm.NFC.String(args[2])

	if len(query) == 0 {
		return fmt.Errorf("no query specified")
	}

	res, err := paperless.Search(ctx, httpClient, baseURL, token, query)
	if err != nil {
		return err
	}

	correspondents, err := paperless.GetCorrespondents(ctx, httpClient, baseURL, token)
	if err != nil {
		return err
	}

	docTypes, err := paperless.GetDocTypes(ctx, httpClient, baseURL, token)
	if err != nil {
		return err
	}

	alfredItems := paperlessToAlfred(res, baseURL, query, correspondents, docTypes)

	out, err := json.Marshal(alfred.Result{Items: alfredItems})
	if err != nil {
		return fmt.Errorf("marshalling alfred results: %w", err)
	}

	_, _ = fmt.Fprintln(stdout, string(out))

	return nil
}

func paperlessToAlfred(
	results []paperless.DocumentSearch,
	baseURL, query string,
	correspondents map[int]paperless.Correspondent,
	docTypes map[int]paperless.DocumentType,
) []alfred.Item {
	var items []alfred.Item

	if len(results) == 0 {
		items = append(items, alfred.Item{
			Title:    "No result found",
			Arg:      fmt.Sprintf("%s/documents?query=%s", baseURL, query),
			Subtitle: "Open search in Paperless",
		})

		return items
	}

	for _, result := range results {
		items = append(items, alfred.Item{
			UID:      result.DetailsURL(baseURL),
			Title:    result.Title,
			Subtitle: strings.Join(result.Metadata(correspondents, docTypes), " - "),
			Arg:      result.DetailsURL(baseURL),
			Icon: &alfred.Icon{
				Type: "filetype",
				Path: "com.adobe.pdf",
			},
			Mods: map[string]alfred.Mod{
				"cmd": {
					Arg:      fmt.Sprintf("%s/documents?query=%s", baseURL, query),
					Subtitle: "Open search in Paperless",
				},
			},
		})
	}

	return items
}
