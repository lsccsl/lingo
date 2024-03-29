package server_common

import (
	"github.com/go-yaml/yaml"
	"goserver/common"
	"os"
	"path/filepath"
	"runtime"
)


type ServerCfg struct {
	MsgQueCent      MsgQueCenterCfg         `yaml:"msgquecenter"`
	MapMsgQueServer map[string]MsgQueSrvCfg `yaml:"msgqueserver"`
	MsgDef          string                  `yaml:"msgdef"`
	MapGSCfg        map[string]SrvCfg       `yaml:"gameserver"`
	MapLogonCfg     map[string]SrvCfg       `yaml:"logonserver"`
	MapDBCfg        map[string]DBSrvCfg     `yaml:"dbserver"`

}

type MsgQueCenterCfg struct {
	BindAddr string `yaml:"bind_addr"`
	OutAddr string `yaml:"out_addr"`
}

type MsgQueSrvCfg struct {
	BindAddr string `yaml:"bind_addr"`
	OutAddr string `yaml:"out_addr"`
}

type SrvCfg struct {
	BindAddr string `yaml:"bind_addr"`
	OutAddr string `yaml:"out_addr"`
}

type DBSrvCfg struct {
	DBUser string `yaml:"db_user"`
	DBPwd string `yaml:"db_pwd"`
	DBIp string `yaml:"db_ip"`
	DBPort int `yaml:"db_port"`
}


var Global_ServerCfg ServerCfg

func ReadCfg(pathCfg string) {
	if len(pathCfg) == 0 {
		_,filename,_,_ := runtime.Caller(0)
		pathCfg = filepath.Dir(filename) + "..\\cfg\\srvcfg.yml"
	}

	yamlFile, err := os.ReadFile(pathCfg)
	if err != nil {
		common.LogErr(err)
		return
	}

	err = yaml.Unmarshal(yamlFile, &Global_ServerCfg)
	if err != nil {
		common.LogErr(err)
		return
	}
	common.LogDebug(&Global_ServerCfg)
}

func GetMsgQueSrvCfg(id string)*MsgQueSrvCfg {
	qCfg, ok := Global_ServerCfg.MapMsgQueServer[id]
	if !ok {
		return nil
	}

	return &qCfg
}

func GetGameSrvCfg(id string)*SrvCfg {
	gsCfg, ok := Global_ServerCfg.MapGSCfg[id]
	if !ok {
		return nil
	}

	return &gsCfg
}

func GetLogonsrvCfg(id string)*SrvCfg {
	logonCfg, ok := Global_ServerCfg.MapLogonCfg[id]
	if !ok {
		return nil
	}

	return &logonCfg
}

func GetDBSrvCfg(id string)*DBSrvCfg {
	dbCfg, ok := Global_ServerCfg.MapDBCfg[id]
	if !ok {
		return nil
	}

	return &dbCfg
}
