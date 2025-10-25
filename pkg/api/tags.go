package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/bananazon/raindrop/pkg/data"
)

func (ac *APIClient) ListTags() (data.ListTagsResult, error) {
	var (
		err            error
		listTagsResult data.ListTagsResult
		listUrl        url.URL
		response       APIResponse
	)

	listUrl = url.URL{Scheme: "https", Host: apiBase, Path: fmt.Sprintf("rest/%s/tags", apiVersion)}
	response = ac.Request(APIRequest{Method: "GET", URL: listUrl})
	if !response.Success {
		return listTagsResult, response.Error
	}

	err = json.Unmarshal(response.Body, &listTagsResult)
	if err != nil {
		return listTagsResult, err
	}

	if !listTagsResult.Result {
		return listTagsResult, fmt.Errorf("list tags returned false: %s", listTagsResult.ErrorMessage)
	}

	return listTagsResult, nil
}

func (ac *APIClient) RenameTag(payload data.RenameTagPayload) (data.RenameTagResult, error) {
	var (
		err             error
		renameTagResult data.RenameTagResult
		renameUrl       url.URL
		response        APIResponse
		urlPath         string
	)

	if payload.CollectionId >= 0 {
		urlPath = fmt.Sprintf("rest/%s/tags/%d", apiVersion, payload.CollectionId)
	} else {
		urlPath = fmt.Sprintf("rest/%s/tags", apiVersion)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return renameTagResult, err
	}

	renameUrl = url.URL{Scheme: "https", Host: apiBase, Path: urlPath}
	response = ac.Request(APIRequest{Method: "PUT", URL: renameUrl, Body: string(jsonData)})

	if !response.Success {
		return renameTagResult, response.Error
	}

	err = json.Unmarshal(response.Body, &renameTagResult)
	if err != nil {
		return renameTagResult, err
	}

	if !renameTagResult.Result {
		return renameTagResult, fmt.Errorf("rename tag returned false: %s", renameTagResult.ErrorMessage)
	}

	return renameTagResult, nil
}

func (ac *APIClient) RemoveTags(payload data.RemoveTagsPayload) (data.RemoveTagsResult, error) {
	var (
		err              error
		removeTagsResult data.RemoveTagsResult
		removeUrl        url.URL
		response         APIResponse
		urlPath          string
	)

	if payload.CollectionId >= 0 {
		urlPath = fmt.Sprintf("rest/%s/tags/%d", apiVersion, payload.CollectionId)
	} else {
		urlPath = fmt.Sprintf("rest/%s/tags", apiVersion)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return removeTagsResult, err
	}

	removeUrl = url.URL{Scheme: "https", Host: apiBase, Path: urlPath}
	response = ac.Request(APIRequest{Method: "DELETE", URL: removeUrl, Body: string(jsonData)})
	if !response.Success {
		return removeTagsResult, response.Error
	}

	err = json.Unmarshal(response.Body, &removeTagsResult)
	if err != nil {
		return removeTagsResult, err
	}

	if !removeTagsResult.Result {
		return removeTagsResult, fmt.Errorf("remove tags returned false: %s", removeTagsResult.ErrorMessage)
	}

	return removeTagsResult, nil
}
