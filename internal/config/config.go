// internal/config/config.go
package config  
  
import (  
	"fmt"  
	"github.com/sirupsen/logrus"  
	"github.com/spf13/viper"  
)  
  
var (  
	DBConnectionString string  
	DebugMode          bool  
	LoggerConfig       LoggerConfiguration
	BaseUrl			   string
)  
  
type LoggerConfiguration struct {  
	LogFile    string  
	LogLevel   logrus.Level  
	MaxSize    int  
	MaxBackups int  
	MaxAge     int  
	Compress   bool  
}  
  
func LoadConfig(configFile string) {  
	viper.SetConfigFile(configFile)  
	err := viper.ReadInConfig()  
	if err != nil {  
		logrus.Warnf("Cannot read config file: %s", err)  
	}  
  
	dbUser := viper.GetString("database.user")  
	dbPassword := viper.GetString("database.password")  
	dbHost := viper.GetString("database.host")  
	dbPort := viper.GetInt("database.port")  
	dbName := viper.GetString("database.name")  

	BaseUrl = viper.GetString("base.url")  
	  
	DBConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)  
  
	DebugMode = viper.GetBool("debug_mode")  
  
	logLevel, _ := logrus.ParseLevel(viper.GetString("logger.log_level"))  
  
	LoggerConfig = LoggerConfiguration{  
		LogFile:    viper.GetString("logger.log_filename"),  
		LogLevel:   logLevel,  
		MaxSize:    viper.GetInt("logger.max_size"),  
		MaxBackups: viper.GetInt("logger.max_backups"),  
		MaxAge:     viper.GetInt("logger.max_age"),  
		Compress:   viper.GetBool("logger.compress"),  
	}  
}  
