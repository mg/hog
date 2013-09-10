package main

import (
	"./ffdb"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage %s DBNAME\n", os.Args[0])
		os.Exit(0)
	}
	db, _ := ffdb.NewFfdbHeader(os.Args[1], "[:]")
	defer db.Close()

	fmt.Println("From state: MA & NY")
	itrNY := db.QueryFieldRx("state", "NY", ffdb.Reverse)
	itrMA := db.QueryFieldRx("state", "MA", ffdb.Forward)

	for !itrNY.AtEnd() || !itrMA.AtEnd() {
		if !itrNY.AtEnd() {
			fmt.Println(itrNY.Value())
			itrNY.Next()
		}
		if !itrMA.AtEnd() {
			fmt.Println(itrMA.Value())
			itrMA.Next()
		}
	}

	fmt.Println("\nOwes more than 100")
	for itr := db.QueryGreater("owes", 100, ffdb.Forward); !itr.AtEnd(); itr.Next() {
		fmt.Println(itr.Value())
	}

}
