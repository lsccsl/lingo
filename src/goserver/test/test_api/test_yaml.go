package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type Yaml struct {
	App struct {
		Name string `yaml:"name",json:"name"`
	}
	Mysql struct {
		Host string `yaml:"host"`
		Port int32 `yaml:"port"`
		DbName string `yaml:"dbName"`
		User string `yaml:"user"`
		Password string `yaml:"password"`
	}
	Cache struct {
		Enable bool `yaml:"enable"`
		List []string `yaml:"list,flow"`
	}
}

func test_yaml()  {
	yamlFile, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		fmt.Println(err)
	}

	conf := new(Yaml)
	yaml.Unmarshal(yamlFile, conf)
	fmt.Println()
	fmt.Println(conf)
}

