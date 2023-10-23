package util  
  
import (  
	"io"  
	"os"  
	"github.com/sirupsen/logrus"  
	"github.com/natefinch/lumberjack"  
	"sqlrest/internal/config"  
)  
  
var Logger = logrus.New()  
  
func LogInit() {  
	// Configure lumberjack as the output for logrus  
	logFile := &lumberjack.Logger{  
		Filename:   config.LoggerConfig.LogFile,  
		MaxSize:    config.LoggerConfig.MaxSize,  
		MaxBackups: config.LoggerConfig.MaxBackups,  
		MaxAge:     config.LoggerConfig.MaxAge,  
		Compress:   config.LoggerConfig.Compress,  
	}  
  
	Logger.SetOutput(io.MultiWriter(logFile, os.Stdout))  
	Logger.SetLevel(config.LoggerConfig.LogLevel)  
}  
