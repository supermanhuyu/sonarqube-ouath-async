package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = New()

func New() *logrus.Logger {
	log := logrus.New()
	/*file, err := os.OpenFile("sonarqube-ouath-async.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}*/
	log.Out = os.Stdout
	return log
}
