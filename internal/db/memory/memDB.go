package memDB

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kamsandhu93/gqldenring/internal/logger"
	"github.com/kamsandhu93/gqldenring/internal/model"
)

var (
	//go:embed elden_ring_weapon.csv
	csv    string
	dbSeed []*model.Weapon
)

func init() {
	dbSeed = parseCsv()
}

type db struct {
	data    []*model.Weapon
	mu      sync.RWMutex
	counter *counter
}

func NewDB() *db {
	return &db{
		data:    dbSeed,
		mu:      sync.RWMutex{},
		counter: newCounter(len(dbSeed)),
	}
}

func (db *db) Printdb() {
	for _, weapon := range db.data {
		fmt.Println(*weapon)
	}
}

func (db *db) AllWeapons(ctx context.Context) ([]*model.Weapon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.data, nil
}

func (db *db) NewWeapon(ctx context.Context, weapon *model.NewWeapon) (*model.Weapon, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	newWeapon := &model.Weapon{
		Name:        weapon.Name,
		Custom:      true,
		ID:          strconv.Itoa(db.counter.increment()),
		LastUpdated: time.Now().Format(time.DateTime),
	}
	db.data = append(db.data, newWeapon)
	logger.LogID(ctx, "[INFO] Created weapon with ID %s", newWeapon.ID)
	return newWeapon, nil
}

func (db *db) UpdateWeapon(ctx context.Context, id string, weapon *model.NewWeapon) (*model.Weapon, error) {
	newWeapon := &model.Weapon{
		Name:        weapon.Name,
		Custom:      true,
		ID:          id,
		LastUpdated: time.Now().Format(time.DateTime),
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	for i, weapon := range db.data {
		if weapon.ID == id {
			db.data[i] = newWeapon
			break
		}
	}
	logger.LogID(ctx, "[INFO] Updated weapon with ID %s", newWeapon.ID)

	return newWeapon, nil
}

func (db *db) DeleteWeapon(ctx context.Context, id string) (*model.Weapon, error) {
	db.mu.Lock()
	var delwep *model.Weapon
	defer db.mu.Unlock()
	for i, weapon := range db.data {
		if weapon.ID == id {
			db.data[i] = db.data[len(db.data)-1]
			db.data = db.data[:len(db.data)-1]
			delwep = weapon
			break
		}
	}
	logger.LogID(ctx, "[INFO] Deleted weapon with ID %s", delwep.ID)

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
			Name:        fields[0],
			Type:        fields[1],
			Phy:         intFromStr(fields[2]),
			Mag:         intFromStr(fields[3]),
			Fir:         intFromStr(fields[4]),
			Lit:         intFromStr(fields[5]),
			Hol:         intFromStr(fields[6]),
			Cri:         intFromStr(fields[7]),
			Sta:         intFromStr(fields[8]),
			Str:         atrScaleFromStr(fields[9]),
			Dex:         atrScaleFromStr(fields[10]),
			Int:         atrScaleFromStr(fields[11]),
			Fai:         atrScaleFromStr(fields[12]),
			Arc:         atrScaleFromStr(fields[13]),
			Any:         fields[14],
			Phyb:        intFromStr(fields[15]),
			Magb:        intFromStr(fields[16]),
			Firb:        intFromStr(fields[17]),
			Litb:        intFromStr(fields[18]),
			Holb:        intFromStr(fields[19]),
			Bst:         fields[20],
			Rst:         fields[21],
			Wgt:         fields[22],
			Upgrade:     fields[23],
			Custom:      false,
			ID:          strconv.Itoa(i),
			LastUpdated: time.Now().Format(time.DateTime),
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
