package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/bananazon/raindrop/pkg/data"
)

func (ac *APIClient) AddBookmark(payload data.AddBookmarkPayload) (data.AddBookmarkResult, error) {
	var (
		addBookmarkResult data.AddBookmarkResult
		addUrl            url.URL
		err               error
		response          APIResponse
	)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return addBookmarkResult, err
	}

	addUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/raindrop", apiVersion)}
	response = ac.Request(APIRequest{Method: "POST", URL: addUrl, Body: string(jsonData)})
	if !response.Success {
		return addBookmarkResult, response.Error
	}

	err = json.Unmarshal(response.Body, &addBookmarkResult)
	if err != nil {
		return addBookmarkResult, err
	}

	if !addBookmarkResult.Result {
		return addBookmarkResult, fmt.Errorf("add bookmark returned false: %s", addBookmarkResult.ErrorMessage)
	}

	return addBookmarkResult, nil
}

func (ac *APIClient) ListBookmarks(page int) (data.ListBookmarksResult, error) {
	var (
		err                 error
		listBookmarksResult data.ListBookmarksResult
		listUrl             url.URL
		queryMap            map[string]string
		response            APIResponse
	)
	queryMap = map[string]string{"sort": "-lastUpdate&perpage", "perpage": strconv.Itoa(PageSize), "page": strconv.Itoa(page), "version": "2"}
	listUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/raindrops/0", apiVersion), RawQuery: MapToQueryString(queryMap)}
	response = ac.Request(APIRequest{Method: "GET", URL: listUrl})
	if !response.Success {
		return listBookmarksResult, response.Error
	}

	err = json.Unmarshal(response.Body, &listBookmarksResult)
	if err != nil {
		return listBookmarksResult, err
	}

	if !listBookmarksResult.Result {
		return listBookmarksResult, fmt.Errorf("list bookmarks returned false: %s", listBookmarksResult.ErrorMessage)
	}

	return listBookmarksResult, nil
}

func (ac *APIClient) RemoveBookmark(bookmarkId int64) (data.RemoveBookmarkResult, error) {
	var (
		err                  error
		removeBookmarkResult data.RemoveBookmarkResult
		removeUrl            url.URL
		response             APIResponse
	)
	removeUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/raindrop/%d", apiVersion, bookmarkId)}
	response = ac.Request(APIRequest{Method: "DELETE", URL: removeUrl})
	if !response.Success {
		return removeBookmarkResult, response.Error
	}

	err = json.Unmarshal(response.Body, &removeBookmarkResult)
	if err != nil {
		return removeBookmarkResult, err
	}

	if !removeBookmarkResult.Result {
		return removeBookmarkResult, fmt.Errorf("remove bookmark returned false: %s", removeBookmarkResult.ErrorMessage)
	}

	return removeBookmarkResult, nil
}

func (ac *APIClient) UpdateBookmark(bookmarkId int64, payload data.UpdateBookmarkPayload) (data.UpdateBookmarkResult, error) {
	var (
		err                  error
		response             APIResponse
		updateBookmarkResult data.UpdateBookmarkResult
		updateUrl            url.URL
	)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return updateBookmarkResult, err
	}

	updateUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/raindrop/%d", apiVersion, bookmarkId)}
	response = ac.Request(APIRequest{Method: "PUT", URL: updateUrl, Body: string(jsonData)})
	if !response.Success {
		return updateBookmarkResult, response.Error
	}

	fmt.Println(string(response.Body))

	err = json.Unmarshal(response.Body, &updateBookmarkResult)
	if err != nil {
		return updateBookmarkResult, err
	}

	if !updateBookmarkResult.Result {
		return updateBookmarkResult, fmt.Errorf("add bookmark returned false: %s", updateBookmarkResult.ErrorMessage)
	}

	return updateBookmarkResult, nil
}
