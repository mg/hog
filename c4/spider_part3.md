4.7: An Extended Example: Web Spiders - Part 3

The *robots* package is very simple since [github.com/temoto/robotstxt-go](http://github.com/temoto/robotstxt-go) does the heavy lifting of parsing and querying the *robot.txt* file.

    type Robot interface{}
    
    func NewRobot(url string) (Robot, error) {
        src := hostFromBase(url) + "/robots.txt"
        resp, err := http.Get(src)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        robots, err := robot.FromResponse(resp)
        if err != nil {
            return nil, err
        }
        return robots.FindGroup("GoGoSpider"), nil
    }
    
    func RobotItr(rules Robot, itr i.Forward) i.Forward {
        rrules, _ := rules.(*robot.Group)
        filter := func(itr i.Iterator) bool {
            url, _ := itr.Value().(string)
            return rrules.Test(url)
        }
        return hoi.Filter(filter, itr)
    }

The *NewRobot()* constructor uses the host name of our root url to fetch the *robots.txt* file and returns the rules for the *GoGoSpider* group (thats us!). The iterator ueses these rules to filter the url stream, rejecting any urls that are not allowed according to the rules in the *robots.txt* file.

The *Fetcher* object is the first, and only object in this system that has nothing at all to do with iterators. It doesn't consume them, doesn't produce them. What it does do is fetch and produce web pages concurrently.

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

The *urlpair* is used to transport an url and its referer from the user of the *Fetcher* to the engine that fetches the pages. The *page* contains those url pairs along with the http response recieved. The *Fetcher* maintains page channel, a queue and a mutex to guard access to the queue. 

    func NewFetcher() *Fetcher {
        f := Fetcher{}
        f.pages = make(chan *page, 5)
        f.done = make(chan bool)
        return &f
    }
    
    func (f *Fetcher) Stop() {
        f.done <- true
        <-f.done
    }

The constructor creates buffered channels to use in the fetching operations. It can fetch 5 pages before it blocks and has to wait.

    func (f *Fetcher) Pages() <-chan *page {
        return f.pages
    }
    
    func (f *Fetcher) Queue(url, ref string) {
        f.lockq.Lock()
        defer f.lockq.Unlock()
        f.queue.PushBack(&urlpair{url: url, ref: ref})
    }
    

The *Queue()* method is where links enter the fetching system. They get stored on a queue so the caller won't block (for long, only if the mutix is locked). The *Pages()* method gives the user access to the page channel to retrieve them. Its where the results of the queueing operation come back out to the user.

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

The *Run()* method is the engine of the *Fetcher*. It starts a Go routine that fetches links of the queue. It performes a *HEAD* request on the link, checking the response and the content type of the response. If the content is an html document it performes a *GET* request, fetching the document. The response of the request that succeded (or the fact that it failed) is then packaged with the url and the referer variables and sent over the *pages* channel.

If the *done* channel is sending a value, the *fetcher* retrieves any documents from the *pages* channel to close them and free any resources. Once that is done, it quites the Go routine. 

PIC

The red line represents the boundaries between the two threads of execution. The queue and the channel cross those boundaries.

Both of these source files are on GitHub at [hog/spider/robots](https://github.com/mg/hog/blob/master/c4/spider/robots.go) and [hog/spider/fetcher](https://github.com/mg/hog/blob/master/c4/spider/fetcher.go).