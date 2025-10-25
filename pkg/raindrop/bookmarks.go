package raindrop

import (
	"os"

	"github.com/bananazon/raindrop/pkg/api"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/util"
	"gopkg.in/yaml.v3"
)

func (r *Raindrop) ListBookmarks() (bookmarks map[uint64]*data.Bookmark, err error) {
	bookmarks, err = r.getAllBookmarks()
	if err != nil {
		return bookmarks, err
	}

	return bookmarks, nil
}

func (r *Raindrop) LoadBookmarks() (err error) {
	bookmarks := make([]*data.Bookmark, 0)

	if util.PathExists(r.BookmarksFile) {
		contents, err := os.ReadFile(r.BookmarksFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(contents, &bookmarks)
		if err != nil {
			return err
		}
	}

	for _, bookmark := range bookmarks {
		r.Bookmarks[bookmark.Id] = bookmark
	}

	return nil
}

func (r *Raindrop) GetChanges() (new []*data.Bookmark, changed []*data.Bookmark, removed []uint64, err error) {
	bookmarks, err := r.getAllBookmarks()
	if err != nil {
		return new, changed, removed, err
	}

	// Find new and changed bookmarks
	for _, bookmark := range bookmarks {
		storedBookmark, exists := r.Bookmarks[bookmark.Id]
		if !exists {
			new = append(new, bookmark)
		} else if bookmark.LastUpdate.After(storedBookmark.LastUpdate) {
			changed = append(changed, bookmark)
		}
	}

	// See if any need deleting
	for _, bookmark := range r.Bookmarks {
		_, exists := bookmarks[bookmark.Id]
		if !exists {
			removed = append(removed, bookmark.Id)
		}
	}

	return new, changed, removed, nil
}

func (r *Raindrop) getAllBookmarks() (bookmarks map[uint64]*data.Bookmark, err error) {
	bookmarks = make(map[uint64]*data.Bookmark)
	page := 0

	for {
		listBookmarksResult, err := r.API.ListBookmarks(page)
		if err != nil {
			return bookmarks, err
		}

		for _, bookmark := range listBookmarksResult.Items {
			bookmarks[bookmark.Id] = bookmark
		}

		over := len(listBookmarksResult.Items) < api.PageSize

		if over {
			break
		}
		page += 1
	}

	return bookmarks, nil
}
