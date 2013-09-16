package spider

import (
	"code.google.com/p/go.net/html"
	"container/list"
	"github.com/mg/i"
	"io"
)

type SearchTactic int

const (
	DepthFirst SearchTactic = iota
	BreathFirst
)

// html node iterator
type nodeitr struct {
	err   error
	cur   *html.Node
	next  func() error
	queue list.List
}

func NodeItr(in io.Reader, stactic SearchTactic) i.Forward {
	n := nodeitr{}
	n.cur, n.err = html.Parse(in)
	if stactic == DepthFirst {
		n.next = n.depthFirst
	} else {
		n.queue.PushBack(n.cur)
		n.next = n.breathFirst
		n.next()
	}
	return &n
}

func (i *nodeitr) Value() interface{} {
	return i.cur
}

func (i *nodeitr) Error() error {
	return i.err
}

func (i *nodeitr) SetError(err error) {
	i.err = err
}

func (i *nodeitr) AtEnd() bool {
	return i.cur == nil
}

func (i *nodeitr) Next() error {
	return i.next()
}

func (i *nodeitr) depthFirst() error {
	if i.err != nil {
		return i.err
	}
	if i.cur.FirstChild != nil {
		i.cur = i.cur.FirstChild
	} else if i.cur.NextSibling != nil {
		i.cur = i.cur.NextSibling
	} else if i.cur.Parent != nil {
		for i.cur.Parent != nil {
			i.cur = i.cur.Parent
			if i.cur.NextSibling != nil {
				i.cur = i.cur.NextSibling
				return i.err
			}
		}
		i.cur = nil
	}
	return i.err
}

func (i *nodeitr) breathFirst() error {
	if i.err != nil {
		return i.err
	}
	i.cur = nil
	if i.queue.Len() > 0 {
		i.cur = i.queue.Front().Value.(*html.Node)
		i.queue.Remove(i.queue.Front())
		if i.cur.FirstChild != nil {
			for c := i.cur.FirstChild; c != nil; c = c.NextSibling {
				i.queue.PushBack(c)
			}
		}
	}
	return i.err
}
