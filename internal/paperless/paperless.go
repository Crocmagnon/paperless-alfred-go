package paperless

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"net/http"
)

func Search(ctx context.Context, client *http.Client, baseURL, token, query string) ([]DocumentSearch, error) {
	var resp PageResponse[DocumentSearch]

	err := requests.URL(baseURL).
		Client(client).
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

func GetCorrespondents(ctx context.Context, client *http.Client, baseURL, token string) (map[int]Correspondent, error) {
	return getAllPagesByID[Correspondent](ctx, client, baseURL, "/api/correspondents/", token)
}

func GetDocTypes(ctx context.Context, client *http.Client, baseURL, token string) (map[int]DocumentType, error) {
	return getAllPagesByID[DocumentType](ctx, client, baseURL, "/api/document_types/", token)
}

type Identifiable interface {
	GetID() int
}

func getAllPagesByID[T Identifiable](ctx context.Context, client *http.Client, baseURL, path, token string) (map[int]T, error) {
	var resp PageResponse[T]

	types := make(map[int]T)

	err := requests.URL(baseURL).
		Client(client).
		Path(path).
		Header("Authorization", "Token "+token).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching doc types: %w", err)
	}

	for _, res := range resp.Results {
		types[res.GetID()] = res
	}

	for resp.Next != nil {
		err := requests.URL(*resp.Next).
			Client(client).
			Header("Authorization", "Token "+token).
			ToJSON(&resp).
			Fetch(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching doc types: %w", err)
		}

		for _, res := range resp.Results {
			types[res.GetID()] = res
		}
	}

	return types, nil
}
