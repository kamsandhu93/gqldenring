package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	"github.com/kamsandhu93/gqldenring/db"
	"github.com/kamsandhu93/gqldenring/graph/generated"
	"github.com/kamsandhu93/gqldenring/graph/model"
)

func (r *mutationResolver) CreateWeapon(ctx context.Context, input *model.NewWeapon) (*model.Weapon, error) {
	return db.NewWeapon(input)
}

func (r *mutationResolver) UpdateWeapon(ctx context.Context, id *string, input *model.NewWeapon) (*model.Weapon, error) {
	return db.UpdateWeapon(*id, input)
}

func (r *mutationResolver) DeleteWeapon(ctx context.Context, id *string) (*model.Weapon, error) {
	return db.DeleteWeapon(*id)
}

func (r *queryResolver) Weapons(ctx context.Context) ([]*model.Weapon, error) {
	weapons := db.Database()
	return weapons, nil
}

func (r *queryResolver) WeaponByName(ctx context.Context, name string) (*model.Weapon, error) {
	weapons := db.Database()
	for _, weapon := range weapons {
		if strings.ToLower(weapon.Name) == strings.ToLower(name) {
			return weapon, nil
		}
	}
	return nil, nil
}

func (r *queryResolver) WeaponsByAttributeScaling(ctx context.Context, attribute model.Attributes, scale model.AttributeScales) ([]*model.Weapon, error) {
	weapons := db.Database()
	results := []*model.Weapon{}
	var atrScale model.AttributeScales
	for _, weapon := range weapons {
		switch attribute {
		case model.AttributesStr:
			atrScale = weapon.Str
		case model.AttributesDex:
			atrScale = weapon.Dex
		case model.AttributesInt:
			atrScale = weapon.Int
		case model.AttributesFai:
			atrScale = weapon.Fai
		case model.AttributesArc:
			atrScale = weapon.Arc
		default:
			atrScale = "Unknown"
		}

		if strings.ToUpper(string(atrScale)) == strings.ToUpper(string(scale)) {
			results = append(results, weapon)
		}
	}
	return results, nil
}

func (r *queryResolver) WeaponsByCustom(ctx context.Context, custom bool) ([]*model.Weapon, error) {
	weapons := db.Database()
	results := []*model.Weapon{}
	for _, weapon := range weapons {
		if weapon.Custom == custom {
			results = append(results, weapon)
		}
	}
	return results, nil
}

func (r *queryResolver) WeaponByID(ctx context.Context, id string) (*model.Weapon, error) {
	weapons := db.Database()
	for _, weapon := range weapons {
		if strings.ToLower(weapon.ID) == strings.ToLower(id) {
			return weapon, nil
		}
	}
	return nil, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
