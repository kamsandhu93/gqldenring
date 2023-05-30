package memDB

import "testing"

//func TestPrintdb(t *testing.T) {
//	Printdb()
//}

func TestIntegrity(t *testing.T) {
	db := NewDB()

	expected := 307
	actual := len(db.data)
	if len(db.data) != expected {
		t.Fatalf("DB has the wrong number of records, expected %d got %d", expected, actual)
	}
}
