package ffdb

import (
	"fmt"
	"github.com/mg/i"
)

type (
	Record map[string]string
)

func RecordItr(db *Ffdb, d Direction) i.Forward {
	return i.Map(recordParser(db), d(db))
}

// Create a Record from line of string
func recordParser(db *Ffdb) i.MapFunc {
	return func(itr i.Iterator) interface{} {
		line, _ := itr.Value().(string)
		fields := db.fieldsep.Split(line, -1)
		if len(fields) != len(db.fields) {
			panic(fmt.Errorf("Unexpected number of fields in record %q, expected %d, got %d\n", line, len(db.fields), len(fields)))
		}
		var record Record = make(map[string]string, len(fields))
		for i, v := range fields {
			record[db.fields[i]] = v
		}
		return &record
	}
}

// Dump record to string
func (r *Record) String() string {
	m := (map[string]string)(*r)
	return fmt.Sprint(m)
}

// Get value of field in record
func (r *Record) Value(key string) string {
	return (*r)[key]
}
