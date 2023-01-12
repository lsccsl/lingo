package server_common

import (
	"github.com/go-yaml/yaml"
	"lin/lin_common"
	"os"
	"path/filepath"
	"runtime"
)

type MsgQueCenterCfg struct {
	BindAddr string `yaml:"bind_addr"`
	OutAddr string `yaml:"out_addr"`
}

type MsgQueSrvCfg struct {
	BindAddr string `yaml:"bind_addr"`
	OutAddr string `yaml:"out_addr"`
}

type ServerCfg struct {
	MsgQueCent MsgQueCenterCfg `yaml:"msgquecenter"`
	MapMsgQueServer map[string]MsgQueSrvCfg `yaml:"msgqueserver"`
	MsgDef string `yaml:"msgdef"`
}

var Global_ServerCfg ServerCfg

func ReadCfg(pathCfg string) {
	if len(pathCfg) == 0 {
		_,filename,_,_ := runtime.Caller(0)
		pathCfg = filepath.Dir(filename) + "..\\cfg\\srvcfg.yml"
	}

	yamlFile, err := os.ReadFile(pathCfg)
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