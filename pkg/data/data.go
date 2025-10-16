package data

import "time"

type Cache struct {
	Status string `json:"status" yaml:"status"`
}

type RaindropCollection struct {
	Ref string `json:"$ref" yaml:"$ref"`
	Id  int64  `json:"$id" yaml:"$id"`
	Oid int64  `json:"oid" yaml:"oid"`
}

type RaindropCreatorRef struct {
	Id     int64  `json:"_id" yaml:"_id"`
	Avatar string `json:"avatar" yaml:"avatar"`
	Name   string `json:"name" yaml:"name"`
	Email  string `email:"link" yaml:"email"`
}

type RaindropLink struct {
	Link string `json:"link" yaml:"link"`
	Type string `json:"type" yaml:"type"`
}

type RaindropUser struct {
	Ref string `json:"$ref" yaml:"$ref"`
	Id  int64  `json:"$id" yaml:"$id"`
}

type Bookmark struct {
	Id           uint64             `json:"_id" yaml:"_id"`
	Link         string             `json:"link" yaml:"link"`
	Title        string             `json:"title" yaml:"title"`
	Excerpt      string             `json:"excerpt" yaml:"excerpt"`
	Note         string             `json:"note" yaml:"note"`
	Type         string             `json:"type" yaml:"type"`
	User         RaindropUser       `json:"user" yaml:"user"`
	Cover        string             `json:"cover" yaml:"cover"`
	Media        []RaindropLink     `json:"media" yaml:"media"`
	Tags         []string           `json:"tags" yaml:"tags"`
	Important    bool               `json:"important" yaml:"important"`
	Removed      bool               `json:"removed" yaml:"removed"`
	Created      time.Time          `json:"created" yaml:"created"`
	Collection   RaindropCollection `json:"collection" yaml:"collection"`
	Highlights   []any              `json:"highlights" yaml:"highlights"`
	LastUpdate   time.Time          `json:"lastUpdate" yaml:"lastUpdate"`
	Domain       string             `json:"domain" yaml:"domain"`
	CreatorRef   RaindropCreatorRef `json:"creatorRef" yaml:"creatorRef"`
	Sort         int64              `json:"sort" yaml:"sort"`
	CollectionId int64              `json:"collectionId" yaml:"collectionId"`
	Cache        Cache              `json:"cache" yaml:"cache"`
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
