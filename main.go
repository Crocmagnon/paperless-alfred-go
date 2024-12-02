package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Crocmagnon/paperless-alfred-go/internal/alfred"
	"github.com/Crocmagnon/paperless-alfred-go/internal/paperless"
	"github.com/carlmjohnson/requests"
	"golang.org/x/text/unicode/norm"
	"io"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout io.Writer) error {
	if len(args) != 3 {
		return fmt.Errorf("usage: ppl-go <base_url> <token> '<query>'")
	}

	baseURL := norm.NFC.String(args[0])
	token := norm.NFC.String(args[1])
	query := norm.NFC.String(args[2])

	if len(query) == 0 {
		return fmt.Errorf("no query specified")
	}

	res, err := search(ctx, baseURL, token, query)
	if err != nil {
		return err
	}

	correspondents, err := getCorrespondents(ctx, baseURL, token)
	if err != nil {
		return err
	}

	alfredItems := paperlessToAlfred(res, baseURL, query, correspondents)

	out, err := json.Marshal(alfred.Result{Items: alfredItems})
	if err != nil {
		return fmt.Errorf("marshalling alfred results: %w", err)
	}

	_, _ = fmt.Fprintln(stdout, string(out))

	return nil
}

func search(ctx context.Context, baseURL, token, query string) ([]paperless.Result, error) {
	var resp paperless.SearchResponse

	err := requests.URL(baseURL).
		Path("/api/documents/").
		Header("Authorization", "Token "+token).
		Param("query", query).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying documents with %q: %w", query, err)
	}

	return resp.Results, nil
}

func getCorrespondents(ctx context.Context, baseURL, token string) (map[int]paperless.Correspondent, error) {
	var resp paperless.CorrespondentsResponse

	corr := make(map[int]paperless.Correspondent)

	err := requests.URL(baseURL).
		Path("/api/correspondents/").
		Header("Authorization", "Token "+token).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching correspondents: %w", err)
	}

	for _, res := range resp.Results {
		corr[res.Id] = res
	}

	for resp.Next != nil {
		err := requests.URL(*resp.Next).
			Header("Authorization", "Token "+token).
			ToJSON(&resp).
			Fetch(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching correspondents: %w", err)
		}

		for _, res := range resp.Results {
			corr[res.Id] = res
		}
	}

	return corr, nil
}

func paperlessToAlfred(results []paperless.Result, baseURL, query string, correspondents map[int]paperless.Correspondent) []alfred.Item {
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
			Subtitle: strings.Join(result.Metadata(correspondents), " - "),
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
