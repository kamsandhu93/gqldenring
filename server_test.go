package main

import (
	"os"
	"testing"

	"github.com/99designs/gqlgen/client"
	memDB "github.com/kamsandhu93/gqldenring/internal/db/memory"
	sqlDB "github.com/kamsandhu93/gqldenring/internal/db/sql"
	"github.com/kamsandhu93/gqldenring/internal/graph"
	"github.com/kamsandhu93/gqldenring/internal/model"
	"github.com/stretchr/testify/require"
)

func TestSrvInMemDB(t *testing.T) {

	db := memDB.NewDB()
	srv := newHandler(graph.NewResolver(db))
	c := client.New(srv)

	runCrudChecks(t, c)
}

func TestSrvSqlDB(t *testing.T) {
	if os.Getenv("SQL_TEST") == "" {
		t.Skip("Skipping testing in sql mode")
	}

	db := sqlDB.NewDB("root:qwerty@tcp(0.0.0.0:3306)/db")
	srv := newHandler(graph.NewResolver(db))
	c := client.New(srv)

	runCrudChecks(t, c)
}

func runCrudChecks(t *testing.T, c *client.Client) {
	t.Run("Get all Weapons check", func(t *testing.T) {
		var resp struct {
			Weapons []*model.Weapon
		}

		q := `{
			  weapons {
				id
				name
				custom
			  }
			}
			`

		c.MustPost(q, &resp)

		require.Equal(t, 307, len(resp.Weapons))
		require.Equal(t, "1", resp.Weapons[0].ID)
		require.Equal(t, "Academy Glintstone Staff", resp.Weapons[0].Name)
		require.Equal(t, false, resp.Weapons[0].Custom)
		require.Equal(t, "307", resp.Weapons[306].ID)
		require.Equal(t, "Zweihander", resp.Weapons[306].Name)
		require.Equal(t, false, resp.Weapons[306].Custom)
	})

	t.Run("Create weapon check", func(t *testing.T) {
		var resp struct {
			CreateWeapon *model.Weapon
		}

		q := `mutation {
			  createWeapon(input: {name:"testNew1"}) {
				id
				name
				custom
			  }
			}
			`

		c.MustPost(q, &resp)

		require.Equal(t, "308", resp.CreateWeapon.ID)
		require.Equal(t, "testNew1", resp.CreateWeapon.Name)
		require.Equal(t, true, resp.CreateWeapon.Custom)

		var resp2 struct {
			CreateWeapon *model.Weapon
		}

		q2 := `mutation {
			  createWeapon(input: {name:"testNew2"}) {
				id
				name
				custom
			  }
			}
			`

		c.MustPost(q2, &resp2)

		require.Equal(t, "309", resp2.CreateWeapon.ID)
		require.Equal(t, "testNew2", resp2.CreateWeapon.Name)
		require.Equal(t, true, resp2.CreateWeapon.Custom)
	})

	t.Run("Update weapon check", func(t *testing.T) {
		var resp struct {
			UpdateWeapon *model.Weapon
		}

		q := `mutation {
			  updateWeapon(id: "309", input: {name: "testUpdate2"}) {
				id
				name
				custom
			  }
			}
			`
		c.MustPost(q, &resp)

		require.Equal(t, "309", resp.UpdateWeapon.ID)
		require.Equal(t, "testUpdate2", resp.UpdateWeapon.Name)
		require.Equal(t, true, resp.UpdateWeapon.Custom)

		var resp2 struct {
			WeaponById *model.Weapon
		}
		q2 := `{
			  WeaponById(id: "309") {
				id
				name
				custom
			  }
			}
			`

		c.MustPost(q2, &resp2)

		require.Equal(t, "309", resp2.WeaponById.ID)
		require.Equal(t, "testUpdate2", resp2.WeaponById.Name)
		require.Equal(t, true, resp2.WeaponById.Custom)
	})
	t.Run("Delete weapon check", func(t *testing.T) {
		var resp struct {
			DeleteWeapon *model.Weapon
		}

		q := `mutation {
			  deleteWeapon(id: "309") {
				id
				name
				custom
			  }
			}
			`
		c.MustPost(q, &resp)

		var resp2 struct {
			Weapons []*model.Weapon
		}

		q2 := `{
			  weapons {
				id
				name
				custom
			  }
			}
			`

		c.MustPost(q2, &resp2)

		found := false
		for _, weapon := range resp2.Weapons {
			if weapon.ID == "309" {
				found = true
			}
		}

		if found {
			require.FailNow(t, "Expected weapon to be deleted but found it")
		}
	})
}
