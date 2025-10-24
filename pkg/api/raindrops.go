package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/bananazon/rdbak/pkg/data"
)

func (ac *APIClient) AddRaindrop(link string, title string, collectionId int64) (data.AddRaindropResult, error) {
	var (
		addRaindropResult data.AddRaindropResult
		addUrl            url.URL
		err               error
		response          APIResponse
	)
	jsonBody := map[string]string{"link": link, "title": title, "collectionId": strconv.FormatInt(collectionId, 10)}
	jsonStr, err := json.Marshal(&jsonBody)
	if err != nil {
		return addRaindropResult, err
	}

	addUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/raindrop", apiVersion)}
	response = ac.Request(APIRequest{Method: "POST", URL: addUrl, Body: string(jsonStr)})
	if !response.Success {
		return addRaindropResult, response.Error
	}

	err = json.Unmarshal(response.Body, &addRaindropResult)
	if err != nil {
		return addRaindropResult, err
	}

	if !addRaindropResult.Result {
		return addRaindropResult, fmt.Errorf("add raindrop returned false: %s", addRaindropResult.ErrorMessage)
	}

	return addRaindropResult, nil
}

func (ac *APIClient) ListRaindrops(page int) (data.ListRaindropsResult, error) {
	var (
		err                 error
		listRaindropsResult data.ListRaindropsResult
		listUrl             url.URL
		queryMap            map[string]string
		response            APIResponse
	)
	queryMap = map[string]string{"sort": "-lastUpdate&perpage", "perpage": strconv.Itoa(PageSize), "page": strconv.Itoa(page), "version": "2"}
	listUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/raindrops/0", apiVersion), RawQuery: MapToQueryString(queryMap)}
	response = ac.Request(APIRequest{Method: "GET", URL: listUrl})
	if !response.Success {
		return listRaindropsResult, response.Error
	}

	err = json.Unmarshal(response.Body, &listRaindropsResult)
	if err != nil {
		return listRaindropsResult, err
	}

	if !listRaindropsResult.Result {
		return listRaindropsResult, fmt.Errorf("list raindrops returned false: %s", listRaindropsResult.ErrorMessage)
	}

	return listRaindropsResult, nil
}

func (ac *APIClient) RemoveRaindrop(raindropId int64) (data.RemoveRaindropResult, error) {
	var (
		err                  error
		removeRaindropResult data.RemoveRaindropResult
		removeUrl            url.URL
		response             APIResponse
	)
	removeUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/raindrop/%d", apiVersion, raindropId)}
	response = ac.Request(APIRequest{Method: "DELETE", URL: removeUrl})
	if !response.Success {
		return removeRaindropResult, response.Error
	}

	err = json.Unmarshal(response.Body, &removeRaindropResult)
	if err != nil {
		return removeRaindropResult, err
	}

	if !removeRaindropResult.Result {
		return removeRaindropResult, fmt.Errorf("remove raindrop returned false: %s", removeRaindropResult.ErrorMessage)
	}

	return removeRaindropResult, nil
}
