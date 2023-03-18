package server_common

import (
	"github.com/go-yaml/yaml"
	"goserver/common"
	"os"
	"path/filepath"
	"runtime"
)

type DataBaseTableCfgDef struct {
	TableName      string `yaml:"table_name"`
	TableProto     string `yaml:"table_proto"`
	QueryKeyProto  string `yaml:"query_key_proto"`
	UpdateKeyProto string `yaml:"update_key_proto"`
	DeleteKeyProto string `yaml:"delete_key_proto"`
}
type DataBaseCfgDef struct {
	DataBaseAppName string `yaml:"database_app_name"`
	DataBase        string `yaml:"database"`
	DataBaseUser    string `yaml:"database_user"`
	DataBasePWD     string `yaml:"database_pwd"`
	DataBaseIP      string `yaml:"database_ip"`
	DataBasePort    int    `yaml:"database_port"`

	Tables []DataBaseTableCfgDef `yaml:"tables"`
}

type DataBaseCfg struct {
	DataBases []DataBaseCfgDef `yaml:"databases"`
}

var Global_DBCfg DataBaseCfg

func ReadDBCfg(pathCfg string) {
	if len(pathCfg) == 0 {
		_,filename,_,_ := runtime.Caller(0)
		pathCfg = filepath.Dir(filename) + "..\\cfg\\dbfg.yml"
	}

	yamlFile, err := os.ReadFile(pathCfg)
	if err != nil {
		common.LogErr(err)
		return
	}

	err = yaml.Unmarshal(yamlFile, &Global_DBCfg)
	if err != nil {
		common.LogErr(err)
		return
	}
	common.LogDebug(&Global_DBCfg)
}

func GetAllDataBaseCfg() []DataBaseCfgDef {
	return Global_DBCfg.DataBases
}
