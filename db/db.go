package db

import (
	_ "embed"
	"fmt"
	"github.com/dgryski/trifles/uuid"
	"github.com/kamsandhu93/gqldenring/graph/model"
	"log"
	"strconv"
	"strings"
	"sync"
)

var (
	//go:embed elden_ring_weapon.csv
	csv string
	db  []*model.Weapon
	mu  sync.RWMutex
)

func init() {
	db = parseCsv()
}

func Printdb() {
	for _, weapon := range db {
		fmt.Println(*weapon)
	}
}

func Database() []*model.Weapon {
	mu.RLock()
	defer mu.RUnlock()
	return db
}

func NewWeapon(weapon *model.NewWeapon) (*model.Weapon, error) {
	newWeapon := &model.Weapon{
		Name:   weapon.Name,
		Custom: true,
		ID:     uuid.UUIDv4(),
	}
	mu.Lock()
	defer mu.Unlock()
	db = append(db, newWeapon)
	log.Printf("[INFO] Created weapon with ID %s", newWeapon.ID)
	return newWeapon, nil
}

func UpdateWeapon(id string, weapon *model.NewWeapon) (*model.Weapon, error) {
	newWeapon := &model.Weapon{
		Name:   weapon.Name,
		Custom: true,
		ID:     id,
	}
	mu.Lock()
	defer mu.Unlock()
	for i, weapon := range db {
		if weapon.ID == id {
			db[i] = newWeapon
			break
		}
	}
	log.Printf("[INFO] Updated weapon with ID %s", newWeapon.ID)

	return newWeapon, nil
}

func DeleteWeapon(id string) (*model.Weapon, error) {
	mu.Lock()
	var delwep *model.Weapon
	defer mu.Unlock()
	for i, weapon := range db {
		if weapon.ID == id {
			db[i] = db[len(db)-1]
			db = db[:len(db)-1]
			delwep = weapon
			break
		}
	}
	log.Printf("[INFO] Deleted weapon with ID %s", delwep.ID)

	return delwep, nil
}

func parseCsv() []*model.Weapon {
	records := strings.Split(csv, "\n")
	parsed := make([]*model.Weapon, 0, len(records))

	for i, record := range records {
		// Skip header
		if i == 0 {
			continue
		}

		if record == "" {
			continue
		}

		fields := strings.Split(record, ",")

		structure := model.Weapon{
			Name:    fields[0],
			Type:    fields[1],
			Phy:     intFromStr(fields[2]),
			Mag:     intFromStr(fields[3]),
			Fir:     intFromStr(fields[4]),
			Lit:     intFromStr(fields[5]),
			Hol:     intFromStr(fields[6]),
			Cri:     intFromStr(fields[7]),
			Sta:     intFromStr(fields[8]),
			Str:     atrScaleFromStr(fields[9]),
			Dex:     atrScaleFromStr(fields[10]),
			Int:     atrScaleFromStr(fields[11]),
			Fai:     atrScaleFromStr(fields[12]),
			Arc:     atrScaleFromStr(fields[13]),
			Any:     fields[14],
			Phyb:    intFromStr(fields[15]),
			Magb:    intFromStr(fields[16]),
			Firb:    intFromStr(fields[17]),
			Litb:    intFromStr(fields[18]),
			Holb:    intFromStr(fields[19]),
			Bst:     fields[20],
			Rst:     fields[21],
			Wgt:     fields[22],
			Upgrade: fields[23],
			Custom:  false,
			ID:      uuid.UUIDv4(),
		}

		parsed = append(parsed, &structure)

	}
	return parsed
}

func intFromStr(str string) int {
	if str == "-" {
		return 0
	}

	result, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return result
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
