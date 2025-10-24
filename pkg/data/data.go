package data

import "time"

type Cache struct {
	Status string `json:"status" yaml:"status"`
}

/*
Refs
*/

type CollectionRef struct {
	Ref  string `json:"$ref" yaml:"$ref"`
	Id   int64  `json:"$id" yaml:"$id"`
	Oid  int64  `json:"oid" yaml:"oid"`
	Name string `json:"name" yaml:"name"`
}

type CreatorRef struct {
	Id     int64  `json:"_id" yaml:"_id"`
	Avatar string `json:"avatar" yaml:"avatar"`
	Name   string `json:"name" yaml:"name"`
	Email  string `email:"link" yaml:"email"`
}

type LinkRef struct {
	Link string `json:"link" yaml:"link"`
	Type string `json:"type" yaml:"type"`
}

type ReminderRef struct {
	Date time.Time `json:"date" yaml:"date"`
}

type UserRef struct {
	Ref string `json:"$ref" yaml:"$ref"`
	Id  int64  `json:"$id" yaml:"$id"`
}

/*
Results
*/

type LoginResult struct {
	Result       bool   `json:"result" yaml:"result"`
	ErrorMessage string `json:"errorMessage" yaml:"errorMessage"`
}

type AddRaindropResult struct {
	Result       bool     `json:"result" yaml:"result"`
	ErrorMessage string   `json:"errorMessage" yaml:"errorMessage"`
	Item         Raindrop `json:"item" yaml:"item"`
}

type ListRaindropsResult struct {
	Result       bool        `json:"result" yaml:"result"`
	Items        []*Raindrop `json:"items" yaml:"items"`
	Count        int         `json:"count" yaml:"count"`
	ErrorMessage string      `json:"errorMessage" yaml:"errorMessage"`
}
type RemoveRaindropResult struct {
	Result       bool   `json:"result" yaml:"result"`
	ErrorMessage string `json:"errorMessage" yaml:"errorMessage"`
}

type AddCollectionResult struct {
	Result       bool       `json:"result" yaml:"result"`
	ErrorMessage string     `json:"errorMessage" yaml:"errorMessage"`
	Item         Collection `json:"item" yaml:"item"`
}

type ListCollectionsResult struct {
	Result       bool          `json:"result" yaml:"result"`
	Items        []*Collection `json:"items" yaml:"items"`
	Count        int           `json:"count" yaml:"count"`
	ErrorMessage string        `json:"errorMessage" yaml:"errorMessage"`
}
type RemoveCollectionResult struct {
	Result       bool   `json:"result" yaml:"result"`
	ErrorMessage string `json:"errorMessage" yaml:"errorMessage"`
}

type SortCollectionsResult struct {
	Result       bool   `json:"result" yaml:"result"`
	ErrorMessage string `json:"errorMessage" yaml:"errorMessage"`
}

/*
Raindrop
*/

type Raindrop struct {
	Id           uint64        `json:"_id" yaml:"_id"`
	Link         string        `json:"link" yaml:"link"`
	Title        string        `json:"title" yaml:"title"`
	Excerpt      string        `json:"excerpt" yaml:"excerpt"`
	Note         string        `json:"note" yaml:"note"`
	Type         string        `json:"type" yaml:"type"`
	User         UserRef       `json:"user" yaml:"user"`
	Cover        string        `json:"cover" yaml:"cover"`
	Media        []LinkRef     `json:"media" yaml:"media"`
	Tags         []string      `json:"tags" yaml:"tags"`
	Important    bool          `json:"important" yaml:"important"`
	Reminder     ReminderRef   `json:"reminder" yaml:"reminder"`
	Removed      bool          `json:"removed" yaml:"removed"`
	Created      time.Time     `json:"created" yaml:"created"`
	Collection   CollectionRef `json:"collection" yaml:"collection"`
	Highlights   []any         `json:"highlights" yaml:"highlights"`
	LastUpdate   time.Time     `json:"lastUpdate" yaml:"lastUpdate"`
	Domain       string        `json:"domain" yaml:"domain"`
	CreatorRef   CreatorRef    `json:"creatorRef" yaml:"creatorRef"`
	Sort         int64         `json:"sort" yaml:"sort"`
	CollectionId int64         `json:"collectionId" yaml:"collectionId"`
	Cache        Cache         `json:"cache" yaml:"cache"`
}

/*
Collection
*/

type CollectionAccess struct {
	For       uint64 `json:"for" yaml:"for"`
	Level     int    `json:"level" yaml:"level"`
	Root      bool   `json:"root" yaml:"root"`
	Draggable bool   `json:"draggable" yaml:"draggable"`
}

type Collection struct {
	Id          uint64           `json:"_id" yaml:"_id"`
	Title       string           `json:"title" yaml:"title"`
	Description string           `json:"description" yaml:"description"`
	User        UserRef          `json:"user" yaml:"user"`
	Public      bool             `json:"public" yaml:"public"`
	View        string           `json:"view" yaml:"view"`
	Count       uint64           `json:"count" yaml:"count"`
	Cover       []string         `json:"cover" yaml:"cover"`
	Expanded    bool             `json:"expanded" yaml:"expanded"`
	CreatorRef  CreatorRef       `json:"creatorRef" yaml:"creatorRef"`
	LastAction  time.Time        `json:"lastAction" yaml:"lastAction"`
	Created     time.Time        `json:"created" yaml:"created"`
	LastUpdate  time.Time        `json:"lastUpdate" yaml:"lastUpdate"`
	Parent      uint64           `json:"parent" yaml:"parent"`
	Sort        uint64           `json:"sort" yaml:"sort"`
	Slug        string           `json:"slug" yaml:"slug"`
	Color       string           `json:"color" yaml:"color"`
	Access      CollectionAccess `json:"access" yaml:"access"`
	Author      bool             `json:"author" yaml:"autor"`
}
