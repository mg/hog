4.7: An Extended Example: Web Spiders - Part 4

The *Spider* is the crawler. It constructs the iterator chain we use to parse the web pages and starts the *Fetcher* process. Its interface to the world is an *i.Forward* iterator that streams out urls and some information about those urls. It operates by sending the url to the *Fetcher* by putting it on the *Fetcher* queue. It waits then for the *Fetcher* to return the first page. The *Spider* then returns the result to the caller and hands the page over to the parsing routine that runs concurrently. The parser runs through the page, queueing any urls that come out of the iterator chain. The user only blocks if there are no pages waiting on the channel from the *Fetcher*.

    type Entry struct {
        Url, Referer string
        StatusCode   int
    }
    
    type Spider struct {
        starturl string
        entries  list.List
        queued   map[string]bool
        err      error
        fetcher  *Fetcher
        robot    Robot
        pages    chan *page
    }

The *Entry* struct is what we return to the user. It contains the url, the page that refered it, and the status code that resulted from attempting to fetch it.

The *Spider* struct contains the state of the operatation: *starturl* is the starting url, entries contains the data we send to the user, *queued* we use as a log of urls we've already worked on, *fetcher* and *robot* are the two main components of teh system, and *pages* is the channel we use to send pages to the parsing funciton.

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

The constructor creates the initial state of the spider and queues and fetches the starting url. It then waits for the first page and sends it to the parser before returning.

    func (s *Spider) Value() interface{} {
        return s.entries.Front().Value
    }
    
    func (s *Spider) Error() error {
        return s.err
    }
    
    func (s *Spider) AtEnd() bool {
        return s.entries.Len() == 0
    }
    
    func (s *Spider) Close() {
        close(s.pages)
        s.fetcher.Close()
    }

The value we return to the user is the head of the *entries* queue. Once its empty, we are done. To clean up resources, the user has to call *Close()* on the spider, it closes the channel to the parser and shuts down the *fetcher*.

    func (s *Spider) Next() error {
        s.entries.Remove(s.entries.Front())
        if s.AtEnd() {
            for {
                if p, ok := <-s.fetcher.Pages(); ok {
                    s.err = s.processEntry(p)
                    if !s.AtEnd() {
                        break
                    }
                }
            }
        }
        return s.err
    }

The *Next()* function pops the head from the queue. If it is empty it fetches the next page from the channel from the *Fetcher* and processes it.

    func (s *Spider) processEntry(p *page) error {
        if p.err != nil {
            return p.err
        }
    
        s.entries.PushBack(&Entry{Url: p.url, Referer: p.ref, StatusCode: p.response.StatusCode})
        s.pages <- p
        return nil
    }

The processing routine simply creates an *Entry* to return to the user and then sends the page off to the parsing function.

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

The parsing function loops on the channel of pages, blocking until a new one is available. When it comes through it starts a chain of iterators to filter the links on the page. Every link that comes through the chain is queued to be processed by the *Fetcher*.

PIC

The iterator channel built in the parsing routine is made of eight components that transform the html nodes to a pair of strings representing the url and its referer. Those processes are for the most part specialitations of *i.Map* and *i.Filter*.

The complete process of the spider looks like this:

PIC

The two red sections represent two independent threads of execution, the *Fetcher* and the *parse()* method.

Using the spider is a simple matter of iterating through the urls that the spider returns. In this case, we want to further filter them by only printin the urls that return a error code. This is accomplished by using a *hoi.FilterFunc* function that checks for the *StatusCode*. Now we have an url checker that checks the sites for urls that are not available. 

    func find4xx(itr i.Iterator) bool {
        e, _ := itr.Value().(*spider.Entry)
        return e.StatusCode >= 400 && e.StatusCode < 500
    }
    
    func main() {
        s := spider.NewSpider(os.Args[1])
        itr := hoi.Filter(find4xx, s)
        count := 0
        for ; !itr.AtEnd(); itr.Next() {
            e, _ := itr.Value().(*spider.Entry)
            count++
            fmt.Printf("%d: Url: %s, Code: %d, Referer: %s\n", count, e.Url, e.StatusCode, e.Referer)
        }
    
        if itr.Error() != nil {
            fmt.Println(itr.Error())
        }
        s.Close()
    
    }

 Both of these source files are on GitHub at [hog/spider/spider](https://github.com/mg/hog/blob/master/c4/spider/spider.go) and [hog/checkurls](https://github.com/mg/hog/blob/master/c4/checkurls.go).
