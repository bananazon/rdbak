package raindrop

import "github.com/bananazon/raindrop/pkg/data"

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
