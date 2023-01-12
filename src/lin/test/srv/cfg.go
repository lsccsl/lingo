package main

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"lin/lin_common"
	"path/filepath"
	"runtime"
)

type ServerOneCfg struct {
	BindAddr string `yaml:"bind_addr"`
	SrvID int64 `yaml:"srv_id"`
	HttpAddr string `yaml:"http_addr"`
	AliasName string `yaml:"alias_name"`
	Cluster string `yaml:"cluster"`
	LogEnableConsolePrint bool `yaml:"log_enable_console_print"`
}

type ServerCfg struct {
	MapServer map[string]ServerOneCfg `yaml:"server"`
	Msgdef string                     `yaml:"msgdef"`
}

var Global_ServerCfg ServerCfg

func ReadCfg(pathCfg string) {
	if len(pathCfg) == 0 {
		_,filename,_,_ := runtime.Caller(0)
		pathCfg = filepath.Dir(filename) + "\\..\\..\\..\\cfg\\cfg.yml"
	}

	yamlFile, err := ioutil.ReadFile(pathCfg)
	if err != nil {
		lin_common.LogErr(err)
		return
	}

	err = yaml.Unmarshal(yamlFile, &Global_ServerCfg)
	if err != nil {
		lin_common.LogErr(err)
		return
	}
	lin_common.LogDebug(&Global_ServerCfg)
}

func GetSrvCfgByID(id string) *ServerOneCfg {
	val, ok := Global_ServerCfg.MapServer[id]
	if !ok {
		return nil
	}
	return &val
}