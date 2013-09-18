package spider

import (
	"container/list"
	"net/http"
	"strings"
	"sync"
	"time"
)

type urlpair struct {
	url, ref string
}

type page struct {
	url, ref string
	err      error
	response *http.Response
}

type Fetcher struct {
	pages  chan *page
	done   chan bool
	queue  list.List
	lockq  sync.Mutex
	client http.Client
}

func NewFetcher() *Fetcher {
	f := Fetcher{}
	f.pages = make(chan *page, 5)
	f.done = make(chan bool)
	return &f
}

func (f *Fetcher) Close() {
	f.done <- true
	<-f.done
}

func (f *Fetcher) Pages() <-chan *page {
	return f.pages
}

func (f *Fetcher) Queue(url, ref string) {
	f.lockq.Lock()
	defer f.lockq.Unlock()
	f.queue.PushBack(&urlpair{url: url, ref: ref})
}

func (f *Fetcher) Run() {
	go func() {
		for {
			f.lockq.Lock()
			if f.queue.Len() > 0 {
				e := f.queue.Front()
				urlpair, _ := e.Value.(*urlpair)
				f.queue.Remove(e)
				f.lockq.Unlock()
				headResp, err := f.client.Head(urlpair.url)
				var p *page
				if err == nil {
					content := headResp.Header.Get("Content-Type")
					if !strings.HasPrefix(content, "text/html") || !strings.HasPrefix(content, "text/xhtml") {
						headResp.Body.Close()
						getResp, err := f.client.Get(urlpair.url)
						if err == nil {
							p = &page{url: urlpair.url, response: getResp, ref: urlpair.ref}
						} else {
							p = &page{url: urlpair.url, ref: urlpair.ref, err: err}
						}
					} else {
						p = &page{url: urlpair.url, ref: urlpair.ref, response: headResp}
					}
				} else {
					p = &page{url: urlpair.url, ref: urlpair.ref, err: err}
				}
				select {
				case f.pages <- p:
				case <-f.done:
					p.response.Body.Close()
					for {
						select {
						case sentpage := <-f.pages:
							sentpage.response.Body.Close()
						default:
							close(f.pages)
							close(f.done)
							return
						}
					}
				}
			} else {
				f.lockq.Unlock()
				time.Sleep(1)
			}
		}
	}()
}
