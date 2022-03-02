package main

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"lin/log"
	"path/filepath"
	"runtime"
)

type ServerCfg struct {
	Cluster string `yaml:"cluster"`
	ClusterAdd []string `yaml:"cluster_addr"`
	BindAddr string `yaml:"bind_addr"`
	AliasName string `yaml:"alias_name"`
}

func ReadCfg() {
	_,filename,_,_ := runtime.Caller(0)
	pathCfg := filepath.Dir(filename) + "\\..\\..\\..\\cfg\\cfg.yml"

	yamlFile, err := ioutil.ReadFile(pathCfg)
	if err != nil {
		log.LogErr(err)
		return
	}

	conf := &ServerCfg{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.LogErr(err)
		return
	}
	log.LogDebug(conf)
}