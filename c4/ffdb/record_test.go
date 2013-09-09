package ffdb

import (
	"testing"
)

func TestRecord(t *testing.T) {
	db, err := NewFfdbHeader("test.db", ":")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for itr := RecordItr(db, Forward); !itr.AtEnd(); itr.Next() {
		t.Log(itr.Value())
	}
}
