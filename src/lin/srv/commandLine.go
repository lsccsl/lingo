package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CmdFuncInfo struct{
	cmdFunc func(argStr []string)
	cmdHelp string
}

type CmdInfo struct {
	mapCmd map[string]CmdFuncInfo
}
var _cmd_info = CmdInfo{
	mapCmd:make(map[string]CmdFuncInfo),
}

func AddCmd(cmd_name string, cmd_help string, cmd_func func(argStr []string)){
	_cmd_info.mapCmd[cmd_name] = CmdFuncInfo{cmd_func, cmd_help}
}

func DumpAllCmd(argStr []string){
	for key, val := range _cmd_info.mapCmd{
		fmt.Println(key, ":", val.cmdHelp)
	}
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
		funcInfo.cmdFunc([]string{})
	}
}

func ParseCmd(){
	//command line
	input_scanner := bufio.NewScanner(os.Stdin)

	for input_scanner.Scan() {
		str := input_scanner.Text()
		fmt.Println("get input:", str)
		if str == "q" {
			break
		}
		if len(str) != 0 {
			arrStr := strings.Fields(str)
			DoCmd(arrStr, len(arrStr))
		}
	}
}
