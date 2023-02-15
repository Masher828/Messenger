package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Address struct {
	Address string `yaml:"address"`
}
type Apps struct {
	Home    Address `yaml:"home"`
	Service Address `yaml:"service"`
	Auth    Address `yaml:"auth"`
}

type DatabaseDetails struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Prefix   string `yaml:"prefix"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Database struct {
	Redis    DatabaseDetails `yaml:"redis"`
	Mongo    DatabaseDetails `yaml:"mongo"`
	Postgres DatabaseDetails `yaml:"postgres"`
}

type IntegrationKeys struct {
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
}

type Integration struct {
	Google   IntegrationKeys `yaml:"google"`
	Razorpay IntegrationKeys `yaml:"razorpay"`
}

type Mail struct {
	From                   string `yaml:"from"`
	Password               string `yaml:"password"`
	Host                   string `yaml:"host"`
	Port                   string `yaml:"port"`
	SendActualMail         string `yaml:"sendActualMail"`
	SendActualMailForError string `yaml:"sendActualMailForError"`
}
type Config struct {
	DomainUrl         string      `yaml:"domainUrl"`
	Apps              Apps        `yaml:"apps"`
	Database          Database    `yaml:"database"`
	HtmlFilePath      string      `yaml:"htmlFilePath"`
	Integrations      Integration `yaml:"integrations"`
	SendActualMessage string      `yaml:"sendActualMessage"`
	Mail              Mail        `yaml:"mail"`
	AssetPath         string      `yaml:"assetPath"`
}

var (
	MessengerConfig = &Config{}
)

func LoadConfigFile() error {
	path := os.Getenv("MESSENGER_CONFIG_FILE")
	// path := "/opt/jc/conf"
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = viper.Unmarshal(MessengerConfig)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
