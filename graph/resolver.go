package graph

import (
	"context"

	"github.com/kamsandhu93/gqldenring/model"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type DB interface {
	NewWeapon(ctx context.Context, weapon *model.NewWeapon) (*model.Weapon, error)
	UpdateWeapon(ctx context.Context, id string, weapon *model.NewWeapon) (*model.Weapon, error)
	DeleteWeapon(ctx context.Context, id string) (*model.Weapon, error)
	Database(ctx context.Context) []*model.Weapon
}

type Resolver struct {
	db DB
}

func NewResolver(db DB) *Resolver {
	return &Resolver{
		db: db,
	}
}
