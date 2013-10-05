package main

import (
	"bytes"
	"fmt"
	"math"
	"time"
)

type set map[string]string
type powerset []set

func keysAndValues(s set) ([]string, []string) {
	var ks []string
	var vs []string
	for k := range s {
		ks = append(ks, k)
		vs = append(vs, s[k])
	}
	return ks, vs
}

func powerset_recurse(s set) powerset {
	var f func(s set, p powerset, keys, values []string, n, i int) powerset
	f = func(s set, p powerset, keys, values []string, n, i int) powerset {
		if p == nil {
			keys, values = keysAndValues(s)
			n = len(keys)
			p = make(powerset, int(math.Pow(2, float64(n)))-1)
			i = 1
		}
		if i-1 == n {
			return p
		}

		c := int(math.Pow(2, float64(i-1)))
		for j := 0; j < c; j++ {
			ss := make(set, i)
			for k := 0; k < i; k++ {
				flag := 1 << uint(k)
				if (c+j)&flag == flag {
					ss[keys[k]] = values[k]
				}
			}
			p[c-1+j] = ss
		}
		return f(s, p, keys, values, n, i+1)
	}
	return f(s, nil, nil, nil, 0, 0)
}

func powerset_loop(s set) powerset {
	keys, values := keysAndValues(s)
	n := len(keys)
	p := make(powerset, int(math.Pow(2, float64(n)))-1)

	for i := 1; i <= n; i++ {
		c := int(math.Pow(2, float64(i-1)))
		for j := 0; j < c; j++ {
			ss := make(set, i)
			for k := 0; k < i; k++ {
				flag := 1 << uint(k)
				if (c+j)&flag == flag {
					ss[keys[k]] = values[k]
				}
			}
			p[c-1+j] = ss
		}
	}
	return p
}

func main() {
	s := set{
		"apple":  "red",
		"banana": "yellow",
		"grape":  "purple",
	}
	d, ps := timeThis(powerset_recurse, s)
	fmt.Printf("Powerset recursive at %d: %v\n", d, ps)
	d, ps = timeThis(powerset_loop, s)
	fmt.Printf("Powerset looping at %d: %v\n", d, ps)
}

func (p powerset) String() string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	for i, _ := range p {
		buf.WriteString(fmt.Sprint(p[i]))
		buf.WriteString(",\n")
	}
	buf.WriteString("}\n")
	return buf.String()
}

func timeThis(f func(set) powerset, s set) (time.Duration, powerset) {
	start := time.Now()
	res := f(s)
	return time.Since(start), res
}
