package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nandesh-dev/subtle/pkgs/config"
)

type l struct {
	path string
}

func (logger *l) Log(title string, content string) {
	file, err := os.OpenFile(logger.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Panicf("Error opening log file: %v", err)
	}

	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%-24s - [  Log  ] - ( %s ) - %s\n", time.Now().Format("2 Jan 2006 15:04:05 MST"), title, content)); err != nil {
		log.Panicf("Error writting to log file: %v", err)
	}
}

func (logger *l) Error(title string, content error) {
	file, err := os.OpenFile(logger.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Panicf("Error opening log file: %v", err)
	}

	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%-24s - [ Error ] - ( %s ) - %s\n", time.Now().Format("2 Jan 2006 15:04:05 MST"), title, content.Error())); err != nil {
		log.Panicf("Error writting to log file: %v", err)
	}
}

var logger *l

func Logger() *l {
	return logger
}

func Init() {
	logger = &l{
		path: config.Config().Logging.Path,
	}
}
