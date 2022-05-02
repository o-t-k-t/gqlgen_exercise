package storage

// import graph gophers with your other imports
import (
	"context"
	"log"
	"net/http"

	"github.com/graph-gophers/dataloader"
	"github.com/o-t-k-t/gqlgen_exercise/graph/model"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// UserReader reads Users from a database
type UserReader struct{}

// GetUsers implements a batch function that can retrieve many users by ID,
// for use in a dataloader
func (u *UserReader) GetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	log.Println("UserReader) GetUsers called", keys)
	output := make([]*dataloader.Result, len(keys))
	for ix, key := range keys {
		output[ix] = &dataloader.Result{
			Data:  &model.User{ID: key.String()},
			Error: nil}
	}
	return output
}

// Dataloader wrap your data loaders to inject via middleware
type Dataloader struct {
	UserLoader *dataloader.Loader
}

// NewDataloader instantiates data loaders for the middleware
func NewDataloader() *Dataloader {
	// define the data loader
	userReader := &UserReader{}
	loaders := &Dataloader{
		UserLoader: dataloader.NewBatchedLoader(userReader.GetUsers),
	}
	return loaders
}

// Middleware injects data loaders into the context
func Middleware(loaders *Dataloader, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Dataloader {
	return ctx.Value(loadersKey).(*Dataloader)
}

// GetUser wraps the User dataloader for efficient retrieval by user ID
func GetUser(ctx context.Context, userID string) (*model.User, error) {
	log.Println("GetUser called")
	loaders := For(ctx)
	thunk := loaders.UserLoader.Load(ctx, dataloader.StringKey(userID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}
