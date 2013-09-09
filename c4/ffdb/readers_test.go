package ffdb

import (
	"testing"
)

func TestReader(t *testing.T) {
	db, err := NewFfdbHeader("test.db", ":")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for itr := Forward(db); !itr.AtEnd(); itr.Next() {
		t.Log(itr.Value())
	}
}

func TestReversedReader(t *testing.T) {
	db, err := NewFfdbHeader("test.db", ":")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for itr := Reverse(db); !itr.AtEnd(); itr.Next() {
		t.Log(itr.Value())
	}
}
