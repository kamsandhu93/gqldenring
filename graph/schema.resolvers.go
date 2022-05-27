package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/kamsandhu93/gqldenring/db"
	"github.com/kamsandhu93/gqldenring/graph/generated"
	"github.com/kamsandhu93/gqldenring/graph/model"
)

func (r *queryResolver) Weapons(ctx context.Context) ([]*model.Weapon, error) {
	weapons := db.Database()
	return weapons, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
