4.3.3: Example - A Flat-File Database - part 3

We conclude this example by going over the *Query* package.

    type (
        FieldQueryFunc func(string) bool
    )

A *FieldQueryFunc* is similar to a *i.FilterFunc* with the exception that it is specialized to work on *string* values rather than *i.Iterator* values.

    func (db *Ffdb) Query(f i.FilterFunc, d Direction) i.Forward {
        return i.Filter(f, RecordItr(db, d))
    }
    
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

The *Query* package contains two fundamental methods that serve as building blocks for all other queries. The first, *Query()*, is an iterator that uses the *i.Filter* iterator to construct a *i.Forward* iterator that will filter a *Record* stream in the direction of *Direction* with the filter *f*.

The *QueryField()* constructs a *i.FilterFunc* and binds it to a field name and a *FieldQueryFunc*. We can now iterate over the stream of *Records* and ask specific questions about individual field values. Two examples are provided in the *Query* package.

    func (db *Ffdb) QueryFieldRx(fieldname, rxs string, d Direction) i.Forward {
        rx := regexp.MustCompile(rxs)
        return db.QueryField(fieldname, func(val string) bool {
            return rx.MatchString(val)
        }, d)
    }
    
    func (db *Ffdb) QueryGreater(fieldname string, check float64, d Direction) i.Forward {
        return db.QueryField(fieldname, func(fieldvalue string) bool {
            value, _ := strconv.ParseFloat(fieldvalue, 64)
            return value > check
        }, d)
    }

*QueryFeildRx* is a query that checks field named *fieldname* against a regular expression *rsx*. *QueryGreater* is a query that returns records where the float64 value of field *fieldname* is greater than some value suppled to the query.

An entire pipeline of iterators, from query to text file using e.g. the *QueryFieldRx*, would look something like this:

PIC

The blue lines represent how data moves through the construction of the pipeline. *FieldName*, *Rx* and *Direction* are send to the *QueryFieldRx* which constructs a *FieldQueryFunc* and sends it along with the *FieldName* and *Direction* to the *QueryField* object. It constructs a *i.FilterFunc* and sends it with the *Direction* to the *Query* object. The *Query* constructs the *Record* iterator, which in turns constructs the *Reader* from the *Direction*.

Now that the pipeline is ready, we can read bytes from the file. The *Reader* turns those into  a string for each line in the file and hands those over to the *Record* iterator. The *Record* iterator, being a *i.Map* iterator, maps the lines to a *Record* which the *Query* then filters and returns to the user.

Following is an example of how to use this package. 

        db, _ := ffdb.NewFfdbHeader(os.Args[1], "[:]")
        defer db.Close()
    
        fmt.Println("From state: MA & NY")
        itrNY := db.QueryFieldRx("state", "NY", ffdb.Reverse)
        itrMA := db.QueryFieldRx("state", "MA", ffdb.Forward)
    
        for !itrNY.AtEnd() || !itrMA.AtEnd() {
            if !itrNY.AtEnd() {
                fmt.Println(itrNY.Value())
                itrNY.Next()
            }
            if !itrMA.AtEnd() {
                fmt.Println(itrMA.Value())
                itrMA.Next()
            }
        }
    
        fmt.Println("\nOwes more than 100")
        for itr := db.QueryGreater("owes", 100, ffdb.Forward); !itr.AtEnd(); itr.Next() {
            fmt.Println(itr.Value())
        }

We open the database, build some query iterators and run through them printing out the result.

The two source files discussed here are available at GitHub: [Query](https://github.com/mg/hog/blob/master/c4/ffdb/query.go) and [ffdb](https://github.com/mg/hog/blob/master/c4/ffdb.go)
