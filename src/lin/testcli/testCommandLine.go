package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CmdFuncInfo struct{
	cmdFunc func(argStr []string)
	cmdHelp string
}

type CmdInfo struct {
	//mapCmd map[string]func(argStr []string, argCount int)
	mapCmd map[string]CmdFuncInfo
	cur_iggid int64
	cmdHelp string
}
var _cmd_info CmdInfo

func set_iggid(argStr []string){
	if len(argStr) < 1{
		return
	}
	_cmd_info.cur_iggid, _ = strconv.ParseInt(argStr[0], 10, 64)
	BeginClrPrint(9)
	fmt.Println("cur iggid:", _cmd_info.cur_iggid)
	EndClrPrint()
}

func InitCmd(){
	_cmd_info.mapCmd = make(map[string]CmdFuncInfo)
	_cmd_info.mapCmd["setiggid"] = CmdFuncInfo{set_iggid, "set cur iggid"}
}

func AddCmd(cmd_name string, cmd_help string, cmd_func func(argStr []string)){
	_cmd_info.mapCmd[cmd_name] = CmdFuncInfo{cmd_func, cmd_help}
}

func DumpAllCmd(){
	BeginClrPrint(9)
	for key, val := range _cmd_info.mapCmd{
		fmt.Println(key, ":", val.cmdHelp)
	}
	EndClrPrint()
}

func DoCmd(argStr []string, argCount int){
	if len(argStr) < 1{
		return
	}

	funcInfo, ok := _cmd_info.mapCmd[argStr[0]]
	if !ok{
		return
	}
	if len(argStr) >= 2{
		funcInfo.cmdFunc(argStr[1:])
	}else{
		funcInfo.cmdFunc(nil)
	}
}

func ParseCmd(){
	//command line
	input_scanner := bufio.NewScanner(os.Stdin)

	for input_scanner.Scan() {
		str := input_scanner.Text()
		fmt.Println("get input:", str)
		if len(str) != 0 {
			arrStr := strings.Fields(str)
			DoCmd(arrStr, len(arrStr))
		} else {
			break
		}
	}
}
