package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/olebedev/config"

	"github.com/gocolly/colly"
)

const (
	baseURL  = "http://room.shmtu.edu.cn:8080"
	loginURL = "https://cas.shmtu.edu.cn/cas/login?service=http://room.shmtu.edu.cn:8080/CAS/docs/examples/cas_simple_login.php"
)

var (
	lt        string
	execution string
)

func main() {
	conf, err := config.ParseYamlFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	username, _ := conf.String("username")
	password, _ := conf.String("password")

	c := colly.NewCollector(
		colly.UserAgent(" Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"))
	loginCollector := c.Clone()
	userCheckCollector := c.Clone()

	c.OnHTML("input[name]", func(e *colly.HTMLElement) {
		switch e.Attr("name") {
		case "lt":
			lt = e.Attr("value")
		case "execution":
			execution = e.Attr("value")
		default:
			return
		}
		return
	})

	c.Visit(loginURL)

	loginCollector.OnResponse(func(r *colly.Response) {
		regP, _ := regexp.Compile(`p=\w+`)
		bodyStr := string(r.Body[:])
		log.Println(bodyStr)
		foundP := regP.FindStringSubmatch(bodyStr)
		if foundP == nil {
			log.Fatal("p not found")
		}
		p := foundP[0][2:]
		userCheckCollector.Visit(
			fmt.Sprintf("%s/Api/auto_user_check?user=%s&p=%s", baseURL, username, p))
	})

	userCheckCollector.OnResponse(func(r *colly.Response) {
		log.Println(userCheckCollector.Cookies(baseURL))
	})

	err = loginCollector.Post(loginURL, map[string]string{
		"username":  username,
		"password":  password,
		"_eventId":  "submit",
		"signin":    "登录",
		"lt":        lt,
		"execution": execution,
	})
	if err != nil {
		log.Fatal(err)
	}

}
