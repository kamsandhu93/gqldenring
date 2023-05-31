package sqlDB

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kamsandhu93/gqldenring/model"
)

const cols = "`Name`, `Type`, `Phy`, `Mag`, `Fir`, `Lit`, `Hol`, `Cri`, `Sta`, `Str`, `Dex`, `Int`, `Fai`, `Arc`, `Any`, `PhyB`, `MagB`, `FirB`, `LitB`, `HolB`, `Bst`, `Rst`, `Wgt`, `Upgrade`, `Custom`"
const vals = "?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?"

type db struct {
	sqlConn string
	db      *sql.DB
}

func NewDB(sqlConn string) *db {
	sqlDB, err := sql.Open("mysql",
		sqlConn)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to database with conn string %s: %v", sqlConn, err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("[ERROR] Pinging database with conn string %s: %v", sqlConn, err)
	}

	return &db{
		sqlConn: sqlConn,
		db:      sqlDB,
	}
}

func (db *db) Database(ctx context.Context) []*model.Weapon {
	results, err := dbQuery(db.db, "select * from weapons")
	if err != nil {
		// needs to be handled properly
		log.Fatalf("Error %v", err)
	}

	return results
}

func (db *db) NewWeapon(ctx context.Context, newWeapon *model.NewWeapon) (*model.Weapon, error) {
	res, err := db.db.Exec(fmt.Sprintf("INSERT INTO weapons(%s) VALUES(%s)", cols, vals),

		newWeapon.Name, "unknown", "0", "0", "0", "0", "0", "0", "0", "-", "-", "-", "-", "-", "", "0", "0", "0", "0", "0", "unknown", "unknown", "unknown", "unknown", true,
	)

	if err != nil {
		log.Printf("[ERROR] inserting row %v", err)
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("[ERROR] getting last id  %v", err)
		return nil, err

	}

	return &model.Weapon{
		Name:   newWeapon.Name,
		ID:     strconv.FormatInt(lastId, 10),
		Custom: true,
	}, nil
}

func (db *db) UpdateWeapon(ctx context.Context, id string, newWeapon *model.NewWeapon) (*model.Weapon, error) {
	res, err := db.db.Exec("UPDATE weapons SET name=? WHERE id=?", newWeapon.Name, id)

	if err != nil {
		log.Printf("[ERROR] inserting row %v", err)
		return nil, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Printf("[ERROR] getting row count  %v", err)
		return nil, err
	}
	log.Printf("[INFO] db update ID = %s, affected = %d\n", id, rowCnt)

	weapons, err := dbQuery(db.db, "SELECT * FROM weapons WHERE id = ?", id)
	if err != nil {
		log.Printf("[ERROR] Retrieving weapons details %v", err)
		return nil, err
	}

	return weapons[0], nil

}

func (db *db) DeleteWeapon(ctx context.Context, id string) (*model.Weapon, error) {
	return &model.Weapon{}, nil
}

func dbQuery(db *sql.DB, query string, values ...any) ([]*model.Weapon, error) {
	rows, err := db.Query(query, values...)
	if err != nil {
		log.Printf("[ERROR] select * error %v", err)
		return nil, fmt.Errorf("select * from db error %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("[ERROR] closing db rows %v", err)
		}
	}(rows)

	var weapons []*model.Weapon

	for rows.Next() {
		var weapon model.Weapon
		var lastUpdated string
		err := rows.Scan(&weapon.Name, &weapon.Type, &weapon.Phy, &weapon.Mag, &weapon.Fir, &weapon.Lit, &weapon.Hol,
			&weapon.Cri, &weapon.Sta, &weapon.Str, &weapon.Dex, &weapon.Int, &weapon.Fai,
			&weapon.Arc, &weapon.Any, &weapon.Phyb, &weapon.Magb, &weapon.Firb, &weapon.Litb, &weapon.Holb,
			&weapon.Bst, &weapon.Rst, &weapon.Wgt, &weapon.Upgrade, &weapon.Custom, &weapon.ID,
			&lastUpdated)

		if err != nil {
			log.Printf("[ERROR] Error scanning rows %v", err)
			return nil, fmt.Errorf("error scanning rows %w", err)
		}

		weapons = append(weapons, &weapon)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("[ERROR] Error fetching all rows from sql db %v", err)
		return nil, fmt.Errorf("error fetching all rows from sql db %w", err)
	}

	return weapons, nil
}

func atrScaleFromStr(atrScale string) model.AttributeScales {
	switch atrScale {
	case "A":
		return model.AttributeScalesA
	case "B":
		return model.AttributeScalesB
	case "C":
		return model.AttributeScalesC
	case "D":
		return model.AttributeScalesD
	case "E":
		return model.AttributeScalesE
	default:
		return model.AttributeScales_
	}
}
