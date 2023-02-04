package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// Config 定义全局配置
var Config *Conf

func Init(path string) {
	Config = LoadConfig(path)
}

// Conf 定义全局配置变量
type Conf struct {
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	MYSQL struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Addr     string `yaml:"addr"`
		Database string `yaml:"database"`
	}
	Jwt struct {
		TokenExpireDuration int    `yaml:"token_expire_duration"` //小时为单位
		Secret              string `yaml:"secret"`
	}
	Server struct {
		Port string `yaml:"port"`
	}
	Oss struct {
		Region          string `yaml:"region"`
		ProviderType    string `yaml:"providerType"`
		Endpoint        string `yaml:"endpoint"`
		AccessKeyId     string `yaml:"access_key_id"`
		AccessKeySecret string `yaml:"access_key_secret"`
		BucketName      string `yaml:"bucket_name"`
	}
	CreateDatabase bool
}

// LoadConfig 获取配置
func LoadConfig(ConfigPath string) *Conf {
	var c = Conf{}
	yamlFile, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Println(err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Println(err.Error())
	}
	return &c
}
