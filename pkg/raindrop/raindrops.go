package raindrop

import (
	"os"

	"github.com/bananazon/rdbak/pkg/api"
	"github.com/bananazon/rdbak/pkg/data"
	"github.com/bananazon/rdbak/pkg/util"
	"gopkg.in/yaml.v3"
)

func (r *Raindrop) ListRaindrops() (raindrops map[uint64]*data.Raindrop, err error) {
	raindrops, err = r.getAllRaindrops()
	if err != nil {
		return raindrops, err
	}

	return raindrops, nil
}

func (r *Raindrop) LoadRaindrops() (err error) {
	bookmarks := make([]*data.Raindrop, 0)

	if util.PathExists(r.Config.BookmarksFile) {
		contents, err := os.ReadFile(r.Config.BookmarksFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(contents, &bookmarks)
		if err != nil {
			return err
		}
	}

	for _, bookmark := range bookmarks {
		r.Raindrops[bookmark.Id] = bookmark
	}

	return nil
}

func (r *Raindrop) GetChanges() (new []*data.Raindrop, changed []*data.Raindrop, removed []uint64, err error) {
	raindrops, err := r.getAllRaindrops()
	if err != nil {
		return new, changed, removed, err
	}

	// Find new and changed bookmarks
	for _, bookmark := range raindrops {
		storedBookmark, exists := r.Raindrops[bookmark.Id]
		if !exists {
			new = append(new, bookmark)
		} else if bookmark.LastUpdate.After(storedBookmark.LastUpdate) {
			changed = append(changed, bookmark)
		}
	}

	// See if any need deleting
	for _, bookmark := range r.Raindrops {
		_, exists := raindrops[bookmark.Id]
		if !exists {
			removed = append(removed, bookmark.Id)
		}
	}

	return new, changed, removed, nil
}

func (r *Raindrop) getAllRaindrops() (raindrops map[uint64]*data.Raindrop, err error) {
	raindrops = make(map[uint64]*data.Raindrop)
	page := 0

	for {
		listRaindropsResult, err := r.API.ListRaindrops(page)
		if err != nil {
			return raindrops, err
		}

		for _, bookmark := range listRaindropsResult.Items {
			raindrops[bookmark.Id] = bookmark
		}

		over := len(listRaindropsResult.Items) < api.PageSize

		if over {
			break
		}
		page += 1
	}

	return raindrops, nil
}
