4.3.3: Example - A Flat-File Database - part 1

The flat file database is a text file in the form of e.g.:

    LASTNAME:FIRSTNAME:CITY:STATE:COUNTRY:OWES
    Adler:David:New York:NY:US:157.00
    Asthon:Elaine:Boston:MA:US:0.00
    Dominus:Mark:Philadelphia:PA:US:0.00
    Orwant:Jon:Cambridge:MA:US:26.30
    Schern:Michael:New York:NY:US:149658.23
    Wall:Larry:Mountain View:CA:US:-372.14
    Gylfason:Magnús:Reykjavík:NA:IS:20.00

The first line is an optional schema that names the individual columns, while ":" is the column separator. If no schema is in the file we have to supply a list of field names to the Ffdb constructor.

An overview over my implementation of this system is as such:

PIC

The *Readers* package provides iterators that allow us to iterate through text files line by line, the *Record* package provides an iterator that transforms a stream of lines into a stream of *Record* structures, the *Ffdb* provides the flat-file database object and lastly the *Query* package implements various examples of queries that are possible. 

This is a somewhat larger example than the previous ones, so we will go through this in more than one post.

I think it is best to go through this bottom up, so lets start down by the metal. The *Readers* package provides two *Forward* iterators, *Forward* and *Reverse*, that iterate through text files line by line. A future solution would be to provide one *BiDirectional* iterator that could fulfill both roles, but this was simpler and will have to do for now.

    type Direction func(*Ffdb) i.Forward

The *Direction* type is simply the type of the constructors of the readers, this allows us to specify when we create the query whether we want to search forwards or reversed simply by passing the name of the reader to the query.

    type reader struct {
        file  *os.File
        pos   int64
        line  string
        err   error
        atEnd bool
        bound Bound
    }

The *reader* structure is the state shared by both the *Forward* reader and the *Reverse* reader. It contains the file we are reading and our position in it, current line, last error, an *atEnd* indicator and the *Bound* structure which is used to bind the readers within a range of the file. This is so we can exclude the schema header if it is located within the file (usually the header). The readers are valid on the range *bound.start <= pos <= bound.end*.

    func (r *reader) Error() error {
        if r.err == io.EOF {
            return nil
        }
        return r.err
    }
    
    func (r *reader) Value() interface{} {
        return r.line
    }
    
    func (r *reader) AtEnd() bool {
        return r.atEnd
    }

These methods are shared between the readers and are therefore bound to the *reader* structure.

    type forward struct {
        reader
        in *bufio.Reader
    }
    
    func Forward(db *Ffdb) i.Forward {
        var itr forward
        if itr.pos, itr.err = db.file.Seek(db.bound.start, os.SEEK_SET); itr.err == nil {
            itr.file = db.file
            itr.bound = db.bound
            itr.in = bufio.NewReader(db.file)
            itr.Next()
        }
        return &itr
    }

The *forward* state is simply the *reader* state plus a *bufio.Reader*. This reader is very similar to the one we created in [4.3.3: Example - Filehandle Iterators]() with the exception of boundaries. The constructor seeks to the beginning position in the file as defined by *bound.Start*, creates the buffered reader from *bufio* and reads the first line.

    func (f *forward) Next() error {
        if f.err == io.EOF || (f.bound.end > 0 && f.pos >= f.bound.end) {
            f.atEnd = true
            return nil
        }
    
        f.file.Seek(f.pos, os.SEEK_SET)
        if f.line, f.err = f.in.ReadString('\n'); f.err != nil && f.err != io.EOF {
            return f.err
        }
        f.pos, _ = f.file.Seek(0, os.SEEK_CUR)
    
        // chomp
        f.line = strings.TrimSuffix(f.line, "\n")
     
        return nil
    }

The *Next()* method seeks to the last position in the file. This is so we can have many iterators on the same database file active at the same time, each maintaining its own position. It reads the line from the file, saves the new position in the file and chomps of the trailing *NL* character. A io.EOF or a position beyond the boundaries indicates that we've reached the end of the database for this iterator.

    const readbufsize = 4096
    
    type reverse struct {
        reader
        readbufpos, readbuflen int64
        readbuf                []byte
        buf                    bytes.Buffer
    }

The *Reverse* reader is a bit more complicated. It inherits from the *reader* and adds various buffers and reading positions to it. The strategy here is to read the file backwards in *readbufsize* chunks, searching through the chunks for *NL* characters. We need to handle both that each chunk might contain many lines and that any line could span across many chunks.

    func Reverse(db *Ffdb) i.Forward {
        whence := os.SEEK_SET
        pos := db.bound.end
        if pos == -1 {
            whence = os.SEEK_END
            pos = 0
        }
    
        var itr reverse
        if itr.pos, itr.err = db.file.Seek(pos, whence); itr.err == nil {
            itr.file = db.file
            itr.bound = db.bound
            itr.readbufpos = -1
            itr.readbuf = make([]byte, readbufsize)
            itr.Next()
        }
        return &itr
    }

The constructor must start by positioning itself at the end of the area that the iterator is valid on, be that the end of the file or the position in *bound.end*. When then proceed to read the first chunk into our buffer.

    func (r *reverse) writeBuffer() {
        var out bytes.Buffer
        out.Write(r.readbuf[r.readbufpos+1 : r.readbuflen])
        if r.buf.Len() > 0 {
            out.Write(r.buf.Bytes())
            r.buf.Reset()
        }
        r.readbuflen = r.readbufpos
        r.readbufpos--
        r.line = out.String()
    }

The *writeBuffer()* writes the current line that is in the chunk to the variable returned by the *Value()* method. It then appends 

    func (r *reverse) Next() error {
        for {
            if r.readbufpos < 0 {
                if r.pos == r.bound.start {
                    r.atEnd = true
                    return nil
                }
                r.readbuflen = readbufsize
                if r.pos < readbufsize {
                    r.readbuflen = r.pos
                }
                r.pos -= r.readbuflen
                r.file.Seek(r.pos, os.SEEK_SET)
                if _, r.err = r.file.Read(r.readbuf); r.err != nil && r.err != io.EOF {
                    return r.err
                }
                r.readbufpos = r.readbuflen - 1
            }
            if r.readbufpos < 0 {
                r.atEnd = true
                return nil
            }
            for ; r.readbufpos >= 0; r.readbufpos-- {
                if r.readbuf[r.readbufpos] == '\n' {
                    r.writeBuffer()
                    return nil
                }
            }
            if r.readbufpos < 0 && r.pos == r.bound.start {
                r.writeBuffer()
                return nil
            }
            var joinbuf bytes.Buffer
            joinbuf.Write(r.readbuf[0:r.readbuflen])
            if r.buf.Len() > 0 {
                joinbuf.Write(r.buf.Bytes())
            }
            r.buf = joinbuf
        }
        return nil
    }

The *Next()* method starts by checking the *readbufpos* variable, if it is below zero we must read a new chunk from the file into our buffer. Normally we try to read *readbufsize* number of characters but if we are closer than that to the start of the file we simply read the rest of the file.

If after reading *readbufpos* is still below zero we are finished. Otherwise we loop through the buffer searching for a *NL* character. If we find it we write out the buffer and return, otherwise we join the buffer with the previous buffer and iterate again. A special case is if we find no *NL* character but move beyond both the buffer and the reading boundary; then we've found the last (first) line in the file.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/ffdb/readers.go).