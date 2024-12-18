package utils

import (
	"log"
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	Sugar *zap.SugaredLogger
}

func NewLogger(fileName string) Logger {

	file, err := os.OpenFile("./logs/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	// Use zap as logger
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{file.Name()}
	logger, err := config.Build()
	if err != nil {
		log.Fatalln(err)
	}

	return Logger{logger.Sugar()}
}
