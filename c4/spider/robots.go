package spider

import (
	"github.com/mg/i"
	"github.com/mg/i/hoi"
	robot "github.com/temoto/robotstxt-go"
	"net/http"
)

type Robot interface{}

func NewRobot(url string) (Robot, error) {
	src := hostFromBase(url) + "/robots.txt"
	resp, err := http.Get(src)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	robots, err := robot.FromResponse(resp)
	if err != nil {
		return nil, err
	}
	return robots.FindGroup("GoGoSpider"), nil
}

func RobotItr(rules Robot, itr i.Forward) i.Forward {
	rrules, _ := rules.(*robot.Group)
	filter := func(itr i.Iterator) bool {
		url, _ := itr.Value().(string)
		return rrules.Test(url)
	}
	return hoi.Filter(filter, itr)
}
