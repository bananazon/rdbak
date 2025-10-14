package data

import "time"

type Cache struct {
	Status string `json:"status" yaml:"status"`
}

type Bookmark struct {
	Id           uint64    `json:"_id" yaml:"_id"`
	Link         string    `json:"link" yaml:"link"`
	Title        string    `json:"title" yaml:"title"`
	Excerpt      string    `json:"excerpt" yaml:"excerpt"`
	Type         string    `json:"type" yaml:"type"`
	Created      time.Time `json:"created" yaml:"created"`
	LastUpdate   time.Time `json:"lastUpdate" yaml:"lastUpdate"`
	CollectionId int64     `json:"collectionId" yaml:"collectionId"`
	Cache        Cache     `json:"cache" yaml:"cache"`
}

type ListResult struct {
	Result       bool        `json:"result" yaml:"result"`
	Items        []*Bookmark `json:"items" yaml:"items"`
	Count        int         `json:"count" yaml:"count"`
	ErrorMessage string      `json:"errorMessage" yaml:"errorMessage"`
}

type LoginResult struct {
	Result       bool   `json:"result" yaml:"result"`
	ErrorMessage string `json:"errorMessage" yaml:"errorMessage"`
}
