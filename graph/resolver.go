package graph

import "github.com/kamsandhu93/gqldenring/model"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type db interface {
	NewWeapon(weapon *model.NewWeapon) (*model.Weapon, error)
	UpdateWeapon(id string, weapon *model.NewWeapon) (*model.Weapon, error)
	DeleteWeapon(id string) (*model.Weapon, error)
	Database() []*model.Weapon
}

type Resolver struct {
	db db
}

func NewResolver(db db) *Resolver {
	return &Resolver{
		db: db,
	}
}
