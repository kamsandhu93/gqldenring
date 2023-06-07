package graph

import (
	"context"

	"github.com/kamsandhu93/gqldenring/internal/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type DB interface {
	NewWeapon(ctx context.Context, weapon *model.NewWeapon) (*model.Weapon, error)
	UpdateWeapon(ctx context.Context, id string, weapon *model.NewWeapon) (*model.Weapon, error)
	DeleteWeapon(ctx context.Context, id string) (*model.Weapon, error)
	AllWeapons(ctx context.Context) ([]*model.Weapon, error)
}

type Resolver struct {
	db DB
}

func NewResolver(db DB) *Resolver {
	return &Resolver{
		db: db,
	}
}
