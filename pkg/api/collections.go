package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gdanko/rdbak/pkg/data"
)

func (ac *APIClient) ListCollections() (listCollectionsResult data.ListCollectionsResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()

	listUrl := url.URL{
		Scheme: "https",
		Host:   apiBase,
		Path:   fmt.Sprintf("%s/collections", apiVersion),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", listUrl.String(), nil)
	if err != nil {
		return listCollectionsResult, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	resp, err := ac.Client.Do(req)
	if err != nil {
		return listCollectionsResult, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return listCollectionsResult, fmt.Errorf("bad status at list bookmarks: %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return listCollectionsResult, err
	}

	err = json.Unmarshal(body, &listCollectionsResult)
	if err != nil {
		return listCollectionsResult, err
	}

	if !listCollectionsResult.Result {
		return listCollectionsResult, fmt.Errorf("list bookmarks returned false: %s", listCollectionsResult.ErrorMessage)
	}

	return listCollectionsResult, nil
}
