package main

import (
	"sonarqube-ouath-async/flag"
	"sonarqube-ouath-async/global"
	"sonarqube-ouath-async/http"
	"sonarqube-ouath-async/log"
	"sonarqube-ouath-async/schedule"
)

func main() {
	flag.Flag()
	err := global.Initialization()
	if err != nil {
		log.Logger.Error(err)
		return
	}

	log.Logger.Info("Start scheduled task")
	schedule.Run()

	log.Logger.Info("Start http service")
	http.Run()
}
