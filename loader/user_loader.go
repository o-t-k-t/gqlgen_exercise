package loader

import (
	"context"
	"log"

	"github.com/graph-gophers/dataloader"
	"github.com/o-t-k-t/gqlgen_exercise/graph/model"
)

func GetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	log.Println("UserReader) GetUsers called", keys)

	output := make([]*dataloader.Result, len(keys))
	for ix, key := range keys {
		output[ix] = &dataloader.Result{
			Data:  &model.User{ID: key.String()},
			Error: nil}
	}
	return output
}

// GetUser wraps the User dataloader for efficient retrieval by user ID
func GetUser(ctx context.Context, userID string) (*model.User, error) {
	log.Println("GetUser called")
	loaders := For(ctx)

	log.Printf("loaders: %+v", loaders)

	thunk := loaders.UserLoader.Load(ctx, dataloader.StringKey(userID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	return result.(*model.User), nil
}
