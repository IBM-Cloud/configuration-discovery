package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

var (
	currentDir             string
	logDir                 string
	stateDir               string
	terraformerfWrapperDir string
)

// ServerConfiguration - holds the values used to configure the application
type ServerConfiguration struct {
	HTTPAddr               string `default:""`
	HTTPPort               int    `default:"8080"`
	ShutdownGraceTime      int    `default:"30"` // minutes
	ShutdownReportInterval int    `default:"1"`  // minutes
	MountDir               string `default:"/tmp" envconfig:"mount_dir"`
}

// MongoConfiguration - holds the values used to configure the application
// todo - Add DB and other details if any
type MongoConfiguration struct {
	Host     string `default:""`
	Port     int    `default:"27017"`
	UserName string `default:"root"`
	Password string `default:"example"`
}

// Configuration holds the values used to configure the application
type Configuration struct {
	Server ServerConfiguration
	Mongo  MongoConfiguration
}

// TheConfiguration :
var theConfiguration *Configuration

// GetConfiguration :
func GetConfiguration() *Configuration {
	return theConfiguration
}

// NewConfiguration :
func NewConfiguration(prefix string) (*Configuration, error) {
	config := new(Configuration)
	if err := envconfig.Process(prefix, &config.Server); err != nil {
		return nil, err
	}
	if err := envconfig.Process(prefix+"_MONGO", &config.Mongo); err != nil {
		return nil, err
	}
	theConfiguration = config
	return config, nil
}

func (conf *Configuration) String() string {
	var sanitizedConf = *conf
	sanitizedConf.Mongo.Password = "****"
	return fmt.Sprintf("{ServerConfig:%+v MongoConig:%+v}",
		sanitizedConf.Server,
		sanitizedConf.Mongo)
}

func init() {
	config, err := NewConfiguration("API")
	if err != nil {
		log.Fatalln("Could not read configuration", err)
	}

	currentDir = config.Server.MountDir
	if currentDir == "" {
		panic("MOUNT_DIR is not set. Please set MOUNT_DIR to continue")
	}

	logDir = currentDir + pathSep + "log"
	stateDir = currentDir + pathSep + "state"
	terraformerfWrapperDir = currentDir + pathSep + "terraformer_wrapper"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, os.ModePerm)
	}
	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		os.MkdirAll(stateDir, os.ModePerm)
	}

}

func SetGlobalDirs(current, logg, state, terraformerfWrapper string) {
	currentDir = current
	logDir = logg
	stateDir = state
	terraformerfWrapperDir = terraformerfWrapper
}
