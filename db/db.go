package db

import (
	_ "embed"
	"fmt"
	"github.com/dgryski/trifles/uuid"
	"github.com/kamsandhu93/gqldenring/graph/model"
	"strconv"
	"strings"
	"sync"
)

var (
	//go:embed elden_ring_weapon.csv
	csv     string
	db      []*model.Weapon
	dbMutex sync.Mutex
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
	dbMutex.Lock()
	defer dbMutex.Unlock()
	return db
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
			Str:     fields[9],
			Dex:     fields[10],
			Int:     fields[11],
			Fai:     fields[12],
			Arc:     fields[13],
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
