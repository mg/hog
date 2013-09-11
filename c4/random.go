package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Random iterator
type Rnd struct {
	seed int64
	cur  float64
	gen  *rand.Rand
}

func Rand() *Rnd {
	var r Rnd
	r.seed = time.Now().UnixNano()
	r.First()
	return &r
}

func (r *Rnd) First() error {
	r.gen = rand.New(rand.NewSource(r.seed))
	r.cur = r.gen.Float64()
	return nil
}

func (r *Rnd) Next() error {
	r.cur = r.gen.Float64()
	return nil
}

func (r *Rnd) AtEnd() bool {
	return false
}

func (r *Rnd) Error() error {
	return nil
}

func (r *Rnd) Value() interface{} {
	return r.cur
}

func (r *Rnd) Float64() float64 {
	return r.cur
}

func (r *Rnd) SetSeed(seed int64) {
	r.seed = seed
	r.First()
}

func (r *Rnd) Seed() int64 {
	return r.seed
}

func main() {

	fmt.Print("A quarted of random pairs: ")
	ritr1, ritr2 := Rand(), Rand()
	for i := 0; i < 4; i++ {
		r1 := ritr1.Float64()
		r2 := ritr2.Float64()
		fmt.Printf("(%f, %f), ", r1, r2)
		ritr1.Next()
		ritr2.Next()
	}
	seed1, seed2 := ritr1.Seed(), ritr2.Seed()
	fmt.Println("")

	fmt.Print("A quarted of another random pairs: ")
	ritr1, ritr2 = Rand(), Rand()
	for i := 0; i < 4; i++ {
		r1 := ritr1.Float64()
		r2 := ritr2.Float64()
		fmt.Printf("(%f, %f), ", r1, r2)
		ritr1.Next()
		ritr2.Next()
	}
	fmt.Println("")

	fmt.Print("The first quarted of random pairs: ")
	ritr1.SetSeed(seed1)
	ritr2.SetSeed(seed2)
	for i := 0; i < 4; i++ {
		r1 := ritr1.Float64()
		r2 := ritr2.Float64()
		fmt.Printf("(%f, %f), ", r1, r2)
		ritr1.Next()
		ritr2.Next()
	}
}
