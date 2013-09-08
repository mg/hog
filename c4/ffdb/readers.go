package ffdb

import (
	"bufio"
	"bytes"
	"github.com/mg/i"
	"io"
	"os"
	"strings"
)

type Direction func(*Ffdb) i.Forward

type reader struct {
	pos   int64
	line  string
	file  *os.File
	err   error
	atEnd bool
	bound Bound
}

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

const readbufsize = 4096

type reverse struct {
	reader
	readbufpos, readbuflen int64
	readbuf                []byte
	buf                    bytes.Buffer
}

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
