package logging

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

func createLogFolder() {
	if _, err := os.Stat(viper.GetString("logFile")); os.IsNotExist(err) {
		err = os.MkdirAll(viper.GetString("logFile"), 0755)
		if err != nil {
			panic(err)
		}
	}
}

func SetupLogFileAndFolder(fileName string) *os.File {
	createLogFolder()
	logFileName := time.Now().Format("2006-01-02") + "-" + fileName + ".log"
	logFile, err := os.OpenFile(viper.GetString("logFile")+logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	return logFile
}
