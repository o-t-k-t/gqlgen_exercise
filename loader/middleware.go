package loader

// import graph gophers with your other imports
import (
	"context"
	"net/http"

	"github.com/graph-gophers/dataloader"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// Dataloader wrap your data loaders to inject via middleware
type Dataloader struct {
	UserLoader *dataloader.Loader
}

// NewDataloader instantiates data loaders for the middleware
func NewDataloader() *Dataloader {
	// define the data loader
	loaders := &Dataloader{
		UserLoader: dataloader.NewBatchedLoader(GetUsers),
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
