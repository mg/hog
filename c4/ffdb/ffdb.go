package ffdb

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
)

type Bound struct {
	start, end int64
}

type Ffdb struct {
	file       *os.File
	fields     []string
	fieldnames map[string]int
	fieldsep   *regexp.Regexp
	bound      Bound
}

// Load new database
func NewFfdb(name, fieldsep string, schema []string) (*Ffdb, error) {
	var db Ffdb
	var err error
	db.fieldsep = regexp.MustCompile(fieldsep)

	if db.file, err = os.Open(name); err != nil {
		return nil, err
	}
	db.fields = schema
	db.bound.end = -1
	return &db, nil
}

func NewFfdbHeader(name, fieldsep string) (*Ffdb, error) {
	db, err := NewFfdb(name, fieldsep, nil)
	if err != nil {
		return nil, err
	}

	readbuf := make([]byte, 4096)
	var buf bytes.Buffer
	var pos int64
	found := false
	for !found {
		n, err := db.file.Read(readbuf)
		if err != nil && err != io.EOF {
			return nil, errors.New("Unable to read schema")
		}
		r := bytes.NewReader(readbuf)
		for r.Len() > 0 {
			ch, _, _ := r.ReadRune()
			if ch == '\n' {
				offset, err := r.Seek(0, os.SEEK_CUR)
				if err != nil {
					return nil, err
				}
				pos += offset
				buf.Write(readbuf[0 : offset-1])
				found = true
				break
			}
		}
		if !found {
			buf.Write(readbuf)
			pos += int64(n)
		}
	}

	db.fields = db.fieldsep.Split(buf.String(), -1)
	db.fieldnames = make(map[string]int, len(db.fields))
	for i, v := range db.fields {
		db.fieldnames[v] = i
	}

	db.bound.start = pos
	return db, nil
}

func (db *Ffdb) Close() {
	db.file.Close()
}
