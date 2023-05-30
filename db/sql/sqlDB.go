package sqlDB

import (
	"context"

	"github.com/kamsandhu93/gqldenring/model"
)

type db struct {
	sqlConn string
}

func NewDB(sqlConn string) *db {
	return &db{
		sqlConn: sqlConn,
	}
}

func (db *db) Database(ctx context.Context) []*model.Weapon {
	return []*model.Weapon{}
}

func (db *db) NewWeapon(ctx context.Context, weapon *model.NewWeapon) (*model.Weapon, error) {
	return &model.Weapon{}, nil
}

func (db *db) UpdateWeapon(ctx context.Context, id string, weapon *model.NewWeapon) (*model.Weapon, error) {
	return &model.Weapon{}, nil
}

func (db *db) DeleteWeapon(ctx context.Context, id string) (*model.Weapon, error) {
	return &model.Weapon{}, nil
}
