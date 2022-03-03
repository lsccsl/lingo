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
	BindAddr string `yaml:"bind_addr"`
	SrvID int `yaml:"srv_id"`
	AliasName string `yaml:"alias_name"`
	HttpAddr string `yaml:"http_addr"`
	MapCluster map[int]string `yaml:"map_cluster"`
}

var Global_ServerCfg ServerCfg

func ReadCfg(pathCfg string) {
	if len(pathCfg) == 0 {
		_,filename,_,_ := runtime.Caller(0)
		pathCfg = filepath.Dir(filename) + "\\..\\..\\..\\cfg\\cfg.yml"
	}

	yamlFile, err := ioutil.ReadFile(pathCfg)
	if err != nil {
		log.LogErr(err)
		return
	}

	err = yaml.Unmarshal(yamlFile, &Global_ServerCfg)
	if err != nil {
		log.LogErr(err)
		return
	}
	log.LogDebug(&Global_ServerCfg)
}