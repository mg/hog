4.3.1: Example - Permutations

The purpose of the *Permute* iterator is to accept a range such as *A..C* and write out all possible permutations:

    ABC
    ACB
    BAC
    BCA
    CAB
    CBA

As always, we start with capturing and storing the state of the iteration, i.e. the permuation.

    type permute struct {
        items []rune
        n     int
        atEnd bool
    }
    
    func Permute(items []rune) i.Forward {
        return &permute{items: items, n: 1, atEnd: false}
    }

*items[]* and *n* are the state of the current permuation while *atEnd* indicates whether we've listed all possible permuations.

    func (p *permute) Error() error {
        return nil
    }
    
    func (p *permute) Value() interface{} {
        return p.items
    }
    
    func (p *permute) AtEnd() bool {
        return p.atEnd
    }

These three functions are very simple, very little to do here apart from the obvious.

    func (per *permute) Next() error {
        var i int
        p := per.n
        for i = 1; i <= len(per.items) && math.Mod(float64(p), float64(i)) == 0; i++ {
            p /= i
        }
        d := int(math.Mod(float64(p), float64(i)))
        j := len(per.items) - i
        if j < 0 {
            per.atEnd = true
            return nil
        }
    
        copy(per.items[j+1:len(per.items)], reverse(per.items[j+1:len(per.items)]))
        per.items[j], per.items[j+d] = per.items[j+d], per.items[j]
        per.n++
        return nil
    }

*Next()* is where the generation of the next permutation occurs, and if the three before were to simple, this one is to clever by far. The idea is to go through the permuations from the last element to the first, such as:

    ABCD > ABDC > ACBD > ACDB > ADBC > ADCB > BADC ...

The first element is held fixed while we go through all permuations of the other elements. Once we are done, we swap the first and second element and go through the process again. If nothing else, it is a very efficient use of modular mathematics and vectors. 

    func reverse(in []rune) []rune {
        out := make([]rune, len(in))
        for i, v := range in {
            out[len(out)-i-1] = v
        }
        return out
    }

*reverse()* is a helper that *Next()* uses to copy one slice to another, reversing the order in the process.

    func generate(from, to rune) []rune {
        list := make([]rune, 0, to-from+1)
        for itr := iutil.Range(int(from), int(to)+1); !itr.AtEnd(); itr.Next() {
            list = append(list, rune(itr.Int()))
        }
        return list
    }

*generate()* accepts paramters such as *'A', 'D'* and returns a slice such as *'A','B','C','D'*. It utilizes *i.Range* to generate an integer interval such as *[63, 67)* that is used to generate the rune slice.

    i.Each(
        Permute(generate(from, to)),
        func(itr i.Iterator) bool {
            r, _ := itr.Value().([]rune)
            fmt.Println(string(r))
            return true
        })

The useage is the same as we've used to before, simply iterator through all possible permuations and print them out.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/permutations.go).