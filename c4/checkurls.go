package main

import (
	"./spider"
	"fmt"
	"github.com/mg/i"
	"github.com/mg/i/hoi"
	"os"
)

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
