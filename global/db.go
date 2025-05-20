package global

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"sonarqube-ouath-async/database"
	"sonarqube-ouath-async/flag"
)

var SonarDb *gorm.DB

func Initialization() (err error) {
	SonarDb, err = database.Setup(flag.Configuration.PostgreSQLDsn)
	if err != nil {
		log.Fatal(fmt.Sprintf("SonarDb setup fail: %s", err.Error()))
	}
	return
}
