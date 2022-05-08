package conf

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func LoadConfigFile() error {
	path := os.Getenv("SOCIAL_CONFIG_FILE")
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
