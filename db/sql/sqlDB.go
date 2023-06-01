package sqlDB

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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

func (db *db) AllWeapons(ctx context.Context) ([]*model.Weapon, error) {
	return dbQuery(db.db, "SELECT * FROM weapons WHERE Deleted=FALSE")
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

	weapons, err := dbQuery(db.db, "SELECT * FROM weapons WHERE id = ?", lastId)
	if err != nil {
		log.Printf("[ERROR] Retrieving weapons details %v", err)
		return nil, err
	}

	return weapons[0], nil
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
	res, err := db.db.Exec("UPDATE weapons SET Deleted=TRUE WHERE id=?", id)

	if err != nil {
		log.Printf("[ERROR] soft deleting row %v", err)
		return nil, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Printf("[ERROR] getting row count  %v", err)
		return nil, err
	}
	log.Printf("[INFO] db delete ID = %s, affected = %d\n", id, rowCnt)

	weapons, err := dbQuery(db.db, "SELECT * FROM weapons WHERE id = ?", id)
	if err != nil {
		log.Printf("[ERROR] Retrieving weapons details %v", err)
		return nil, err
	}

	return weapons[0], nil
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
		var deleted bool
		err := rows.Scan(&weapon.Name, &weapon.Type, &weapon.Phy, &weapon.Mag, &weapon.Fir, &weapon.Lit, &weapon.Hol,
			&weapon.Cri, &weapon.Sta, &weapon.Str, &weapon.Dex, &weapon.Int, &weapon.Fai,
			&weapon.Arc, &weapon.Any, &weapon.Phyb, &weapon.Magb, &weapon.Firb, &weapon.Litb, &weapon.Holb,
			&weapon.Bst, &weapon.Rst, &weapon.Wgt, &weapon.Upgrade, &weapon.Custom, &weapon.ID,
			&weapon.LastUpdated, &deleted)

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
