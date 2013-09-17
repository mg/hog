4.7: An Extended Example: Web Spiders - Part 2

The job of the *linkitr.go* package is to operate on a stream of html nodes from the *NodeItr* iterator and transform it into a stream of strings.

The first iterator is the *LinkItr*.

    func links(itr i.Iterator) bool {
        n, _ := itr.Value().(*html.Node)
    
        if n.Type != html.ElementNode {
            return false
        }
        if n.Data != "a" && n.Data != "img" && n.Data != "link" && n.Data != "style" && n.Data != "script" {
            return false
        }
    
        if n.Data == "style" || n.Data == "script" {
            src := attr("src", n.Attr)
            return src != ""
        }
        return true
    }
    
    func LinkItr(itr i.Forward) i.Forward {
        return hoi.Filter(links, itr)
    }

The *LinkItr* is an *i.Filter* iterator that runs through the node stream and removes any nodes that are not *a*, *img*, *link* or *style* nodes.
    
    func geturl(itr i.Iterator) interface{} {
        n, _ := itr.Value().(*html.Node)
        var url string
        if n.Data == "a" {
            url = attr("href", n.Attr)
        } else if n.Data == "img" {
            url = attr("src", n.Attr)
        } else if n.Data == "link" {
            url = attr("href", n.Attr)
        } else if n.Data == "style" {
            url = attr("srr", n.Attr)
        } else if n.Data == "script" {
            url = attr("src", n.Attr)
        }
        return url
    }
    
    func UrlItr(itr i.Forward) i.Forward {
        return hoi.Map(geturl, itr)
    }

The *UrlItr* uses *hoi.Map* to transform the node stream into a string stream, from now on we are operating on a stream of urls.

    func attr(name string, attrs []html.Attribute) string {
        for _, a := range attrs {
            if a.Key == name {
                return a.Val
            }
        }
        return ""
    }

The *attr()* method is a helper that gets the value of attribute *name* from the node.

One improvement that could be made to this code would be to implement both the *links* method and the *geturl* method with a *DispatchTable* as discussed in chapter 2 of *HOP*.

After this we hand the stream over to the iterators in the *urlitr.go* package. It contains five iterators, *NormalizeItr*, *HostMapper*, *Referer*, *BindByRef* and *BindByHost*.

    func removeUrl(itr i.Iterator) bool {
        url, _ := itr.Value().(string)
        url = strings.TrimSpace(url)
        if url == "" {
            return false
        }
        if strings.HasPrefix(url, "tel:") || strings.HasPrefix(url, "mailto:") {
            return false
        }
        return true
    }
    
    func remHash(u string) string {
        idx := strings.Index(u, "#")
        if idx == -1 {
            return u
        }
        if idx == 0 {
            return ""
        }
        return u[0 : idx-1]
    }
    
    func normalizeUrl(itr i.Iterator) interface{} {
        url, _ := itr.Value().(string)
        url = strings.TrimSuffix(remHash(url), "/")
        return strings.ToLower(url)
    }
    
    func NormalizeItr(itr i.Forward) i.Forward {
        return hoi.Filter(removeUrl, hoi.Map(normalizeUrl, itr))
    }

The job of *NormalizeItr* is to trasform the url into a canalogical version of the url. The hash portion gets cut of, trailing slash chomped, the rest lower cased, emtpy strings excluded along with *mailto* and *tel* links. The iterator uses both *hoi.Map* and *hoi.Filter* to accomplish this.

    func hostFromBase(base string) string {
        slashslashidx := strings.Index(base, "//")
        idx := strings.Index(base[slashslashidx+2:], "/")
        if idx > 0 {
            return base[0 : slashslashidx+2+idx]
        }
        return base
    }
    
    func hostMapper(base string) i.MapFunc {
        host := hostFromBase(base)
        if !strings.HasSuffix(base, "/") {
            base = base + "/"
        }
        return func(itr i.Iterator) interface{} {
            url, _ := itr.Value().(string)
            if strings.Contains(url, "://") {
                return url
            }
            if strings.HasPrefix(url, "/") {
                return host + url
            }
            return base + url
        }
    }
    
    func HostMapper(base string, itr i.Forward) i.Forward {
        return hoi.Map(hostMapper(base), itr)
    }

The *HostMapper* runs over the url stream, turning any relative urls into absolute urls.

    func referer(ref string) hoi.MapFunc {
        return func(itr i.Iterator) interface{} {
            url, _ := itr.Value().(string)
            return []string{url, ref}
        }
    }
    
    func Referer(ref string, itr i.Forward) i.Forward {
        return hoi.Map(referer(ref), itr)
    }

The *Referer* transforms the url stream into a stream of twinlets, the url and the page it was found on.

    func removeIfNotReferedBy(ref string) hoi.FilterFunc {
        host := hostFromBase(ref)
        return func(itr i.Iterator) bool {
            ref, _ := itr.Value().([]string)
            return strings.HasPrefix(ref[1], host)
        }
    }
    
    func BindByRef(ref string, itr i.Forward) i.Forward {
        return hoi.Filter(removeIfNotReferedBy(ref), itr)
    }

The *BindByRef* removes any urls from the stream that are not refered by the reference url. This is to stop the crawler to run away into the distance, we are only checking the urls that are refered by the host of the starting url. Those urls might take us to another website but we wont go any further.

    func removeIfNotOnHost(url string) hoi.FilterFunc {
        host := hostFromBase(url)
        return func(itr i.Iterator) bool {
            u, _ := itr.Value().([]string)
            return strings.HasPrefix(u[0], host)
        }
    }
    
    func BindByHost(host string, itr i.Forward) i.Forward {
        return hoi.Filter(removeIfNotOnHost(host), itr)
    }

The *BindByHost* is an iterator that removes any urls from the stream that are not on the same host as the referencing url. If you use this iterator on the stream, you are only crawling urls that are one the same domain as the starting url.

Both of these source files are on GitHub at [hog/spider/linkitr](https://github.com/mg/hog/blob/master/c4/spider/linkitr.go) and [hog/spider/urlitr](https://github.com/mg/hog/blob/master/c4/spider/urlitr.go).
