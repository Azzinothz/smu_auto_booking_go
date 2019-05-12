package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gocolly/colly"
)

const (
	baseURL  = "http://room.shmtu.edu.cn:8080"
	loginURL = "https://cas.shmtu.edu.cn/cas/login?service=http://room.shmtu.edu.cn:8080/CAS/docs/examples/cas_simple_login.php"
)

type checkedCollector struct {
	username  string
	password  string
	lt        string
	execution string
	p         string
	collector *colly.Collector
}

func (cc *checkedCollector) newCheckedCollector() (newCollector *colly.Collector, err error) {
	cc.collector = colly.NewCollector(
		colly.UserAgent(" Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"))
	cc.getLtAndExecutionValue()
	err = cc.loginAndGetP()
	if err != nil {
		newCollector = nil
		return
	}

	newCollector = cc.collector.Clone()
	cc.collector = newCollector
	newCollector.Visit(
		fmt.Sprintf("%s/Api/auto_user_check?user=%s&p=%s", baseURL, cc.username, cc.p))
	return
}

func (cc *checkedCollector) getLtAndExecutionValue() {
	cc.collector.OnHTML("input[name]", func(e *colly.HTMLElement) {
		switch e.Attr("name") {
		case "lt":
			cc.lt = e.Attr("value")
		case "execution":
			cc.execution = e.Attr("value")
		default:
			return
		}
		return
	})
	cc.collector.Visit(loginURL)
}

func (cc *checkedCollector) loginAndGetP() (err error) {
	cc.collector.OnResponse(func(r *colly.Response) {
		regP, _ := regexp.Compile(`p=\w+`)
		bodyStr := string(r.Body[:])
		log.Println(bodyStr)
		foundP := regP.FindStringSubmatch(bodyStr)
		if foundP == nil {
			log.Fatal("p not found")
		}
		cc.p = foundP[0][2:]
	})
	err = cc.collector.Post(loginURL, map[string]string{
		"username":  cc.username,
		"password":  cc.password,
		"_eventId":  "submit",
		"signin":    "登录",
		"lt":        cc.lt,
		"execution": cc.execution,
	})
	return
}

func (cc *checkedCollector) bookRoom(startTime string, endTime string, day string, title string, application string, teamusers []string, mobile string) (err error) {
	c := cc.collector.Clone()
	requestDataStr := fmt.Sprintf(
		"startTime=%s&endTime=%s&day=%s&title=%s&application=%s&mobile=%s&userid=%s&type=%d&isPublic=%t",
		startTime,
		endTime,
		day,
		title,
		application,
		mobile,
		cc.username,
		2,
		false)
	for _, teamuser := range teamusers {
		requestDataStr += "&teamusers[]=" + teamuser
	}
	for _, cookie := range c.Cookies(baseURL) {
		if cookie.Name == "access_token" {
			requestDataStr += "&access_token=" + cookie.Value
			break
		}
	}
	requestData := []byte(requestDataStr)

	c.OnResponse(func(r *colly.Response) {
		log.Println(r.Body)
	})

	err = c.PostRaw(baseURL+"/api.php/spaces/3070/studybook", requestData)
	return
}
