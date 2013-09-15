package ffdb

import (
	"fmt"
	"github.com/mg/i"
	"github.com/mg/i/hoi"
	"regexp"
	"strconv"
	"strings"
)

type (
	FieldQueryFunc func(string) bool
)

// Basic query builder, accepts a function to process a record, returns an iterator.
// Maintains file position state
func (db *Ffdb) Query(f hoi.FilterFunc, d Direction) i.Forward {
	return hoi.Filter(f, RecordItr(db, d))
}

// Build a query iterator to process a single field
func (db *Ffdb) QueryField(fieldname string, f FieldQueryFunc, d Direction) i.Forward {
	fieldname = strings.ToUpper(fieldname)
	if _, ok := db.fieldnames[fieldname]; !ok {
		panic(fmt.Errorf("Unknown field: %q\n", fieldname))
	}
	return db.Query(func(itr i.Iterator) bool {
		r, _ := itr.Value().(*Record)
		return f(r.Value(fieldname))
	}, d)
}

// Query iterator that matches a field value to a regular expression
func (db *Ffdb) QueryFieldRx(fieldname, rxs string, d Direction) i.Forward {
	rx := regexp.MustCompile(rxs)
	return db.QueryField(fieldname, func(val string) bool {
		return rx.MatchString(val)
	}, d)
}

// Query iterator that checks if floating value of a field is greater than supplied value
func (db *Ffdb) QueryGreater(fieldname string, check float64, d Direction) i.Forward {
	return db.QueryField(fieldname, func(fieldvalue string) bool {
		value, _ := strconv.ParseFloat(fieldvalue, 64)
		return value > check
	}, d)
}
