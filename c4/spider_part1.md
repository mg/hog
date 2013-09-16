4.7: An Extended Example: Web Spiders - Part 1

The spider example is, like the [Flat File Database](http://higherordergo.blogspot.com/2013/09/433-example-flat-file-database-part-1.html), long enough to need more than one post. And as with *Ffdb* I will attempt to explain the code from the bottom up. 

But first the birds view of the system. The *Node* iterator gives us access to the *html* document. We then use various iterators to transform that stream of nodes into a stream of links. The *fetcher*'s job is to retrieve documents from the web, and the *robot*'s job is to parse the *robot.txt* file so we can act as good citizens when hitting some poor server. The *spider* runs the whole show and maintains the state of the crawling. One crucial change I made from the example in *HOP*: the fetching operation is now a concurrent operation. This is a blog about Go and its characteristics after all.

PIC

The *Node* iterators job is to take a document tree of html nodes and turn it into a one dimensional stream of nodes. The iterator offers two strategies to do this, a *Depth-First* and a *Breath-First*.

    type SearchTactic int
    
    const (
        DepthFirst SearchTactic = iota
        BreathFirst
    )
    
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

The *nodeitr* maintains a reference to the current node, a queue for the *Breath-First* search and function pointer to the search algorithm. The constructor retrieves the root node and assigns the correct search function to *next* according to the *SearchTactic* value.

    func (i *nodeitr) Value() interface{} {
        return i.cur
    }
    
    func (i *nodeitr) Error() error {
        return i.err
    }
    
    func (i *nodeitr) AtEnd() bool {
        return i.cur == nil
    }
    
    func (i *nodeitr) Next() error {
        return i.next()
    }

The only thing to node here is that the only thing *Next()* does is to forward the operation to whatever function is in the *next* variable, be that the *Depth-First* or the *Breath-First* search.

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

The *depthFirst()* simply uses the *FirstChild*, *NextSibling* and *Parent* references that every node has to traverse the tree in a *Depth-First* fashion. If the current node has a child we go there. Else if it has a sibling we go there. If both fails we start looping up the parent references, stopping on the first one that has a sibling. If that fails we are back at the root node and the iteration is finished.

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

Our *breathFirst()* function uses a queue to maintain a list of nodes that we are yet to visit. If the queue is not empty we pop the first node of the queue; that is our current position. Then we check if this node has any child nodes; if so they get appended to the queue. If the queue is empty we are finished.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/spider/nodeitr.go).