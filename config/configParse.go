package config

import (

	//"gopkg.in/yaml.v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

func Initialize( ) GlobalConfig {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	//var config backend.RedisConfig
	var config GlobalConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
		panic(err)
	}
	return config
}
