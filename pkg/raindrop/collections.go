package raindrop

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/util"
	"gopkg.in/yaml.v3"
)

func (r *Raindrop) BackupCollections(flagPrune bool) (err error) {
	r.PruneOlder = flagPrune

	r.Logger.Info("Starting collections backup.")

	err = r.LoadCollections()
	if err != nil {
		return err
	}

	// Get updated and new bookmarks
	newCollections, changedCollections, removedCollections, err := r.GetCollectionChanges()
	if err != nil {
		return err
	}

	// Merge unchanged collections with changed/new
	keptIds := make(map[uint64]bool)
	r.UpdatedBookmarks = make([]*data.Bookmark, 0, len(r.Collections)+len(changedCollections))

	for _, collection := range newCollections {
		r.UpdatedCollections = append(r.UpdatedCollections, collection)
		keptIds[collection.Id] = true
	}

	for _, collection := range changedCollections {
		r.UpdatedCollections = append(r.UpdatedCollections, collection)
		keptIds[collection.Id] = true
	}

	for _, collection := range r.Collections {
		if _, exists := keptIds[collection.Id]; exists {
			continue
		}

		if !slices.Contains(removedCollections, collection.Id) {
			r.UpdatedCollections = append(r.UpdatedCollections, collection)
		}
	}

	r.PruneBackupFiles("collection")

	if len(newCollections) > 0 || len(changedCollections) > 0 || len(removedCollections) > 0 {
		err = r.SaveCollectionsBackupFile()
		if err != nil {
			return err
		}
	}

	// Report
	var changedString string = "collections"
	var newString string = "collections"
	var removedString = "collections"

	if len(changedCollections) == 1 {
		changedString = "collection"
	}

	if len(newCollections) == 1 {
		newString = "collection"
	}

	if len(removedCollections) == 1 {
		removedString = "collection"
	}
	r.Logger.Infof(
		"Finished. %d new %s; %d changed %s; %d removed %s.",
		len(newCollections),
		newString,
		len(changedCollections),
		changedString,
		len(removedCollections),
		removedString,
	)

	return nil
}

func (r *Raindrop) ListCollections() (collections map[uint64]*data.Collection, err error) {
	collections = make(map[uint64]*data.Collection)

	listCollectionsResult, err := r.API.ListCollections()
	if err != nil {
		return collections, err
	}

	for _, collection := range listCollectionsResult.Items {
		collections[collection.Id] = collection
	}

	return collections, nil
}

func (r *Raindrop) LoadCollections() (err error) {
	collections := make([]*data.Collection, 0)

	if util.PathExists(r.CollectionsFile) {
		contents, err := os.ReadFile(r.CollectionsFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(contents, &collections)
		if err != nil {
			return err
		}
	}

	for _, collection := range collections {
		r.Collections[collection.Id] = collection
	}

	return nil
}

func (r *Raindrop) GetCollectionChanges() (new []*data.Collection, changed []*data.Collection, removed []uint64, err error) {
	collections, err := r.getAllCollections()
	if err != nil {
		return new, changed, removed, err
	}

	// Find new and changed bookmarks
	for _, collection := range collections {
		storedCollection, exists := r.Collections[collection.Id]
		if !exists {
			new = append(new, collection)
		} else if collection.LastUpdate.After(storedCollection.LastUpdate) {
			changed = append(changed, collection)
		}
	}

	// See if any need deleting
	for _, collection := range r.Collections {
		_, exists := collections[collection.Id]
		if !exists {
			removed = append(removed, collection.Id)
		}
	}

	return new, changed, removed, nil
}

func (r *Raindrop) getAllCollections() (collections map[uint64]*data.Collection, err error) {
	collections = make(map[uint64]*data.Collection)

	listCollectionsResult, err := r.API.ListCollections()
	if err != nil {
		return collections, err
	}

	for _, collection := range listCollectionsResult.Items {
		collections[collection.Id] = collection
	}

	return collections, nil
}

func (r *Raindrop) SaveCollectionsBackupFile() (err error) {
	yamlCollections, err := yaml.Marshal(r.UpdatedCollections)
	if err != nil {
		return nil
	}

	if util.PathExists(r.CollectionsFile) {
		backupFilename := filepath.Join(r.RaindropRoot, fmt.Sprintf("collections-%d.yaml", time.Now().Unix()))

		r.Logger.Infof("Copying %s to %s.", r.CollectionsFile, backupFilename)
		r.Logger.Infof("Saving collections to %s.", r.CollectionsFile)

		err = os.Rename(r.CollectionsFile, backupFilename)
		if err != nil {
			return fmt.Errorf("failed to rename %s to %s: %s", r.CollectionsFile, backupFilename, err.Error())
		}
	}

	err = os.WriteFile(r.CollectionsFile, yamlCollections, 0600)
	if err != nil {
		return err
	}

	return nil
}
