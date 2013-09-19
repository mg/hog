package spider

import (
	"container/list"
	"strings"
	"sync"
)

type Entry struct {
	Url, Referer string
	StatusCode   int
}

type Spider struct {
	starturl string
	entries  list.List
	elock    sync.Mutex
	queued   map[string]bool
	err      error
	fetcher  *Fetcher
	robot    Robot
	pages    chan *page
}

func NewSpider(url string) *Spider {
	url = strings.ToLower(url)
	var s Spider
	s.starturl = url
	s.robot, _ = NewRobot(url)
	s.queued = make(map[string]bool)
	s.queued[url] = true
	s.fetcher = NewFetcher()
	s.fetcher.Queue(url, url)
	s.fetcher.Run()
	s.pages = make(chan *page)
	s.parse()
	p := <-s.fetcher.Pages()
	s.err = s.processEntry(p)
	return &s
}

func (s *Spider) Value() interface{} {
	return s.entries.Front().Value
}

func (s *Spider) Error() error {
	return s.err
}

func (s *Spider) SetError(err error) {
	s.err = err
}

func (s *Spider) AtEnd() bool {
	return s.entries.Len() == 0
}

func (s *Spider) Close() {
	close(s.pages)
	s.fetcher.Close()
}

func (s *Spider) Next() error {
	s.entries.Remove(s.entries.Front())
	if s.AtEnd() {
		for {
			p := <-s.fetcher.Pages()
			s.err = s.processEntry(p)
			if !s.AtEnd() {
				break
			}
		}
	}
	return s.err
}

func (s *Spider) processEntry(p *page) error {
	if p.err != nil {
		return p.err
	}

	s.entries.PushBack(&Entry{Url: p.url, Referer: p.ref, StatusCode: p.response.StatusCode})
	s.pages <- p
	return nil
}

func (s *Spider) parse() {
	go func() {
		for p := range s.pages {
			if p.response.Body != nil {
				itr :=
					HostMapper(p.url,
						NormalizeItr(
							UrlItr(
								LinkItr(
									NodeItr(p.response.Body, DepthFirst)))))

				if s.robot != nil {
					itr = RobotItr(s.robot, itr)
				}
				itr = BindByRef(s.starturl, Referer(p.url, itr))

				for ; !itr.AtEnd(); itr.Next() {
					urlpair, _ := itr.Value().([]string)
					url := urlpair[0]
					if _, ok := s.queued[url]; !ok {
						s.queued[url] = true
						s.fetcher.Queue(url, urlpair[1])
					}
				}

				p.response.Body.Close()
			}
		}

	}()
}
