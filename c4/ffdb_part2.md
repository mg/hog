4.3.3: Example - A Flat-File Database - part 2

The *Record* package iterates over a stream of lines, parsing them into a *Record* structure and passing them on. It utilizes a new function in the [iterator library](https://github.com/mg/i/blob/master/hoi/map.go), *Map*. Lets start by running through it.

    type MapFunc func(Iterator) interface{}
    
    type imap struct {
        fmap MapFunc
        val  interface{}
        itr  Forward
    }
    
    func Map(fmap MapFunc, itr Forward) Forward {
        return &imap{fmap: fmap, itr: itr, val: nil}
    }
    
    func (i *imap) AtEnd() bool {
        return i.itr.AtEnd()
    }
    
    func (i *imap) Next() error {
        i.val = nil
        return i.itr.Next()
    }
    
    func (i *imap) Value() interface{} {
        if i.val == nil {
            i.val = i.fmap(i.itr)
        }
        return i.val
    }
    
    func (i *imap) Error() error {
        return i.itr.Error()
    }

The *Map* iterator accepts a *MapFunc* function and a *Forward* iterator. It then iterates over that iterator and uses the *MapFunc* to transform the value before returning it. You can see source at [GitHub](https://github.com/mg/i/blob/master/hoi/map.go).

    type (
        Record map[string]string
    )
    
    func RecordItr(db *Ffdb, d Direction) i.Forward {
        return i.Map(recordParser(db), d(db))
    }

The *Record* type is simply a map of strings keyed to a string. And the constructor creates a map from a stream of lines to a stream of records as defined by *recordParser*. The *Direction* is a function pointer to either the *Forward* or the *Reverse* reader.

    func recordParser(db *Ffdb) hoi.MapFunc {
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

The *recordParser* returns a function of type *hoi.MapFunc*. Its purpose is to accept a string, tokenize it with a field seperator defined in the *Ffdb* object, validate that the number of fields in the line matches with the number of fields in the *Ffdb* object, create a *Record* structure and fill it with the values from the line keyed to the fields defined in the *Ffdb* object.

    func (r *Record) String() string {
        m := (map[string]string)(*r)
        return fmt.Sprint(m)
    }
    
    func (r *Record) Value(key string) string {
        return (*r)[key]
    }

Finally we have some helper functions that allows to work with a *Record* structure.

The *Ffdb* package defines the flat file database object.

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

The *Bound* structure we've seen before, it defines the boundaries of the underlying text file for the readers. The *Ffdb* structure holds the neccessary state for the database: the name of the file, name of the fields, a mapping of field names to field numbers, the field seperator and the *Bound* structure.

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

We have to constructors for the *Ffdb* object. The former assumes that there is no schema in the file and therefore expects the caller to supplie a slice of strings containing the names of the fields (and indirectly the number of fields). It defines the boundaries as the entire file.

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

The second constructor assumes that the first line in the text file contains the schema of the data base. It attempts to read this line and parse it to fill both the *fields* and the *fieldnames* members. It also saves the ending position of the first line and uses that to define the boundaries for the readers.

    func (db *Ffdb) Close() {
        db.file.Close()
    }

The *Ffdb* is at heart a file resource, and must therefore be closed to avoid resource leaks.

The two source files discussed here are available at GitHub: [Record](https://github.com/mg/hog/blob/master/c4/ffdb/record.go) and [Ffdb](https://github.com/mg/hog/blob/master/c4/ffdb/ffdb.go)
