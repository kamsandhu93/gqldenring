package db

import "testing"

//func TestPrintdb(t *testing.T) {
//	Printdb()
//}

func TestIntegrity(t *testing.T) {
	expected := 307
	actual := len(db)
	if len(db) != expected {
		t.Fatalf("DB has the wrong number of records, expected %d got %d", expected, actual)
	}
}
