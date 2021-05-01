package main

import (
    "io/ioutil"
    "os"
    "log"
    
    "github.com/yuin/gopher-lua"
	"github.com/SevereCloud/vksdk/v2/api"
)

type Command struct {
    Cmd string
    Fn *lua.LFunction
    Module *lua.LState
}

var (
    LoadedModules []*lua.LState
    AllCmds []*Command
)

func LoadModules() {
    // Загружаем список модулей в папке
    files := readDirUtil("modules")
    for _, mod := range files {
        // Работаем с каждой папкой по отдельности
        // Проверяем папку на наличие стартового файла
        if hasFile("modules/" + mod.Name(), "init.lua") {
            
            // Кешируем наш модуль
            L := lua.NewState()
            L.SetGlobal("botAddCommand", L.NewFunction(botAddCommand))
            L.SetGlobal("botSendMessage", L.NewFunction(botSendMessage))
            if err := L.DoFile("modules/" + mod.Name() + "/init.lua"); err != nil {
                panic(err)
            }
            LoadedModules = append(LoadedModules, L)
        }
        // Запускаем инициализацию
        startInitFuncs()
    }
}

func CloseModules() {
    for _, L := range LoadedModules {
        // Закрываем каждый модуль
        L.Close()
    }
}

func botSendMessage(L *lua.LState) int {
    peerId := L.ToInt(1)
    text := L.ToString(2)
    
    vk.MessagesSend(api.Params{
		"peer_id":    peerId,
		"message":    text,
		"attachment": "",
		"random_id":  0,
	})
    
    return 0
}

func botAddCommand(L *lua.LState) int {
    cmd := L.ToString(1)
    fn := L.ToFunction(2)
    has := false
    for _, cmmnd := range AllCmds {
        if cmmnd.Cmd == cmd {
            has = true
            break
        }
    }
    if !has {
        AllCmds = append(AllCmds, &Command{
            Cmd: cmd,
            Fn: fn,
            Module: L,
        })
    } else {
        log.Println("Команда \"" + cmd + "\" уже существует" )
    }
    L.Push(lua.LBool(has))
    return 1
}

func startInitFuncs() {
    for _, L := range LoadedModules {
        if err := L.CallByParam(lua.P{
            Fn: L.GetGlobal("onLoaded"),
            NRet: 0,
            Protect: true,
            }); err != nil {
            log.Println(err)
        }
    }
}

func hasFile(dir, file string) bool {
    had := false
    files := readDirUtil(dir)
    for _, f := range files {
        if f.Name() == file {
            had = true
            break
        }
    }
    return had
}

func readDirUtil(dir string) []os.FileInfo {
    files, err := ioutil.ReadDir(dir)
	if err != nil {
	    // Если нихуя, то кидаем ошиб очку
		log.Fatal(err)  // Нихуя
	}
	return files
}