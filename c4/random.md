4.3.6: Example - Random Number Generator

Our final example of an iterator is short and simple, it is a random number generator.

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

The structure contains the seed used to generate the sequence, the current value (a floating number in the [0..1] range) and a reference to Go's random number generator. 

The constructor simply seeds a random sequence from the current time and generates the first value.

    func (r *Rnd) First() error {
        r.gen = rand.New(rand.NewSource(r.seed))
        r.cur = r.gen.Float64()
        return nil
    }
    
    func (r *Rnd) Next() error {
        r.cur = r.gen.Float64()
        return nil
    }

The iterator is a *BoundedAtStart* iterator, i.e. it has a beginning and can easily be restored to it through the *First()* method. The *Next()* generates the next random value.

    func (r *Rnd) AtEnd() bool {
        return false
    }
    
    func (r *Rnd) Error() error {
        return nil
    }

This iterator represents an infinite stream of values on the form *[rval1, rval2, ...)*, therefore the only job of the *AtEnd()* method is to return *false*.

    func (r *Rnd) Value() interface{} {
        return r.cur
    }
    
    func (r *Rnd) Float64() float64 {
        return r.cur
    }

The iterator has a special type casting function so we can read *float64* values from it without have to do the type assertion ourselves.

    func (r *Rnd) SetSeed(seed int64) {
        r.seed = seed
        r.First()
    }
    
    func (r *Rnd) Seed() int64 {
        return r.seed
    }

Lastsly, the iterator provids methods to retrieve the seeding used to generate the sequence, as well as reset and iterator to a previously saved seeding value.

A sampel usage is as follows:

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

We generate three list of a quarted of a pair of random numbers, where the first sequence and the last sequence are the same.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/random.go). The *Random* iterator is also awailable in the [iterator library](https://github.com/mg/i/blob/master/igen/rand.go) as a *RandomAccess* iterator.