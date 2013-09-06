4.3.2: Example - Genomic Sequence Generator

Given the pattern *A(CG)G(AT)*, the *Gene* iterator will produce the sequence:

    ACGA
    AGGA
    ACGT
    AGGT

Strings within *()* represent possible permutations while strings outside *()* are static.

    type gene struct {
        atEnd bool
        data  []interface{}
        cur   string
    }
    
    func Gene(pattern string) i.Forward {
        r := regexp.MustCompile("[()]")
        tokens := r.Split(pattern, -1)
        g := gene{atEnd: false, data: make([]interface{}, 0, len(tokens))}
        for i, v := range tokens {
            if math.Mod(float64(i), 2) == 0 {
                // constant
                if v != "" {
                    g.data = append(g.data, v)
                }
            } else {
                // option
                if v != "" {
                    option := make([]interface{}, 0, len(v)+1)
                    option = append(option, 0)
                    for _, c := range v {
                        option = append(option, string(c))
                    }
                    g.data = append(g.data, option)
                }
            }
        }
        g.Next()
        return &g
    }

The state of the iteration is held in the *data* variable. The constructor accepts a pattern string a loops through it creating a datastructure such as:

    A(GA)C => ['A',['G','A'],'C']

A string element represent a static element while a slice element represents the possible permuations that we need to loop through. Calling *Next()* creates the first permuation.

    func (g *gene) Error() error {
        return nil
    }
    
    func (g *gene) Value() interface{} {
        return g.cur
    }
    
    func (g *gene) AtEnd() bool {
        return g.cur == ""
    }

The only thing we need to talk about here is that the *atEnd* variable realy means that there are no more permuations to generate. But there is still a value to return as log as *cur* is not empty.

    func (g *gene) Next() error {
        if g.atEnd {
            g.cur = ""
            return nil
        }
        finishedIncr := false
        var result bytes.Buffer
        for _, t := range g.data {
            if v, ok := t.(string); ok {
                result.WriteString(v)
            } else if v, ok := t.([]interface{}); ok {
                n := v[0].(int)
                result.WriteString(v[n+1].(string))
                if !finishedIncr {
                    if n == len(v)-2 {
                        v[0] = 0
                    } else {
                        v[0] = v[0].(int) + 1
                        finishedIncr = true
                    }
                }
            }
        }
        g.atEnd = !finishedIncr
        g.cur = result.String()
        return nil
    }

*Next()* starts by checking the *atEnd* variable, clearing *cur* and returning if it is true. *AtEnd()* will now return true.

Otherwise we loop through the *data* variable, generating the next value for the sequence. 

    i.Each(
        Gene(os.Args[1]),
        func(itr i.Iterator) bool {
            fmt.Println(itr.Value())
            return true
        })

The main loop is the same as before.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/genes.go).