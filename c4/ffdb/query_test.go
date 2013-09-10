package ffdb

import (
	"testing"
)

func TestQueryGreater(t *testing.T) {
	db, err := NewFfdbHeader("test.db", ":")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for itr := db.QueryGreater("owes", 100, Forward); !itr.AtEnd(); itr.Next() {
		t.Log(itr.Value())
	}
}

func TestQueryFieldRx(t *testing.T) {
	db, err := NewFfdbHeader("test.db", ":")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	itrNY := db.QueryFieldRx("state", "NY", Reverse)
	itrMA := db.QueryFieldRx("state", "MA", Forward)

	for !itrNY.AtEnd() && !itrMA.AtEnd() {
		t.Log(itrNY.Value())
		t.Log(itrMA.Value())
		itrNY.Next()
		itrMA.Next()
	}
}
