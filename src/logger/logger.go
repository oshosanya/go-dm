package logger

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var instance *logrus.Logger
var once sync.Once

//Get an instance of logger
func GetInstance() *logrus.Logger {
	once.Do(func() {
		instance = logrus.New()
		var filename = "logfile.log"
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			// Cannot open log file. Logging to stderr
			fmt.Println(err)
		} else {
			instance.SetOutput(f)
			instance.SetLevel(logrus.DebugLevel)
		}
	})
	return instance
}
