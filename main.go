package main

import (
	"log"

	"github.com/olebedev/config"
)

func main() {
	conf, err := config.ParseYamlFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	username, _ := conf.String("username")
	password, _ := conf.String("password")

	cc := checkedCollector{username: username, password: password}
	_, err = cc.newCheckedCollector()
	if err != nil {
		log.Fatal(err)
	}
	err = cc.bookRoom("21:00", "21:30", "2019-05-16", "立项相关内容", "立项相关内容立项相关内容立项相关内容", []string{"201610733001", "201610733002", "201610733003", "201610733004", "201610733005", "201610733006"}, "15021617205")
	if err != nil {
		log.Fatal(err)
	}
}
