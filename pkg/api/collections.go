package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/bananazon/rdbak/pkg/data"
)

func (ac *APIClient) AddCollection(title string, parentId int64, public bool) (data.AddCollectionResult, error) {
	var (
		addCollectionResult data.AddCollectionResult
		addUrl              url.URL
		err                 error
		response            APIResponse
	)

	jsonBody := map[string]any{"title": title, "parent": strconv.FormatInt(parentId, 10), "public": public}
	jsonStr, err := json.Marshal(&jsonBody)
	if err != nil {
		return addCollectionResult, err
	}

	addUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/collection", apiVersion)}
	response = ac.Request(APIRequest{Method: "POST", URL: addUrl, Body: string(jsonStr)})
	if !response.Success {
		return addCollectionResult, response.Error
	}

	err = json.Unmarshal(response.Body, &addCollectionResult)
	if err != nil {
		return addCollectionResult, err
	}

	if !addCollectionResult.Result {
		return addCollectionResult, fmt.Errorf("add collection returned false: %s", addCollectionResult.ErrorMessage)
	}

	return addCollectionResult, nil
}

func (ac *APIClient) ListCollections() (data.ListCollectionsResult, error) {
	var (
		err                   error
		listCollectionsResult data.ListCollectionsResult
		listUrl               url.URL
		response              APIResponse
	)

	listUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/collections", apiVersion)}
	response = ac.Request(APIRequest{Method: "GET", URL: listUrl})
	if !response.Success {
		return listCollectionsResult, response.Error
	}

	err = json.Unmarshal(response.Body, &listCollectionsResult)
	if err != nil {
		return listCollectionsResult, err
	}

	if !listCollectionsResult.Result {
		return listCollectionsResult, fmt.Errorf("list collections returned false: %s", listCollectionsResult.ErrorMessage)
	}

	return listCollectionsResult, nil
}

func (ac *APIClient) SortCollections(sortOrder string) (data.SortCollectionsResult, error) {

	var (
		err                   error
		sortCollectionsResult data.SortCollectionsResult
		sortUrl               url.URL
		response              APIResponse
	)

	jsonBody := map[string]any{"sort": sortOrder}
	jsonBytes, err := json.Marshal(&jsonBody)
	if err != nil {
		return sortCollectionsResult, err
	}

	sortUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/collections", apiVersion)}
	response = ac.Request(APIRequest{Method: "PUT", URL: sortUrl, Body: string(jsonBytes)})
	if !response.Success {
		return sortCollectionsResult, response.Error
	}

	err = json.Unmarshal(response.Body, &sortCollectionsResult)
	if err != nil {
		return sortCollectionsResult, err
	}

	if !sortCollectionsResult.Result {
		return sortCollectionsResult, fmt.Errorf("sort collections returned false: %s", sortCollectionsResult.ErrorMessage)
	}

	return sortCollectionsResult, nil
}

func (ac *APIClient) RemoveCollection(collectionId int64) (data.RemoveCollectionResult, error) {
	var (
		err                    error
		removeCollectionResult data.RemoveCollectionResult
		removeUrl              url.URL
		response               APIResponse
	)
	removeUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("%s/collection/%d", apiVersion, collectionId)}
	response = ac.Request(APIRequest{Method: "DELETE", URL: removeUrl})
	if !response.Success {
		return removeCollectionResult, response.Error
	}

	err = json.Unmarshal(response.Body, &removeCollectionResult)
	if err != nil {
		return removeCollectionResult, err
	}

	if !removeCollectionResult.Result {
		return removeCollectionResult, fmt.Errorf("remove raindrop returned false: %s", removeCollectionResult.ErrorMessage)
	}

	return removeCollectionResult, nil
}
