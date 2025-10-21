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

func (ac *APIClient) ListRaindrops(page int) (listRaindropsResult data.ListRaindropsResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()

	listUrl := url.URL{
		Scheme:   "https",
		Host:     apiBase,
		Path:     fmt.Sprintf("%s/raindrops/0", apiVersion),
		RawQuery: fmt.Sprintf("sort=-lastUpdate&perpage=%d&page=%d&version=2", PageSize, page),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", listUrl.String(), nil)
	if err != nil {
		return listRaindropsResult, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	resp, err := ac.Client.Do(req)
	if err != nil {
		return listRaindropsResult, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return listRaindropsResult, fmt.Errorf("bad status at list bookmarks: %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return listRaindropsResult, err
	}

	err = json.Unmarshal(body, &listRaindropsResult)
	if err != nil {
		return listRaindropsResult, err
	}

	if !listRaindropsResult.Result {
		return listRaindropsResult, fmt.Errorf("list bookmarks returned false: %s", listRaindropsResult.ErrorMessage)
	}

	return listRaindropsResult, nil
}
