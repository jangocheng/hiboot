// TODO: app config should be generic kit

package system

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

type Configuration struct {
	App App `mapstructure:"app"`
}

var conf Configuration

func Build() *Configuration {
	_, filename, _, _ := runtime.Caller(0)
	workdir := strings.Replace(filename, "boot/pkg/system/builder.go", "", -1)
	log.Debug(workdir)

	viper.SetDefault("app.project", "devops")
	viper.SetDefault("app.name", "hi")
	viper.SetDefault("app.server.port", 8080)
	viper.SetDefault("app.logging.level", "info")
	viper.AddConfigPath(workdir + "/config")
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error config file: %s \n", err))
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}
	appProfile := os.Getenv("APP_PROFILES_ACTIVE")
	if appProfile != "" {
		viper.SetConfigName("app-" + appProfile)
		viper.MergeInConfig()
		err = viper.Unmarshal(&conf)
		if err != nil {
			panic(fmt.Errorf("Unable to decode Config: %s \n", err))
		}
	}

	return &conf
}
