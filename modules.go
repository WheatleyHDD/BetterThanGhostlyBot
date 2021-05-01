package main

import (
    "io/ioutil"
    "os"
    "log"
    
    "github.com/WheatleyHDD/BetterThanGhostlyBot/botmodule"
    "github.com/WheatleyHDD/BetterThanGhostlyBot/globals"
    
    "github.com/yuin/gopher-lua"
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
            
            L.PreloadModule("bot", botmodule.Loader)
            
            if err := L.DoFile("modules/" + mod.Name() + "/init.lua"); err != nil {
                panic(err)
            }
            globals.LoadedModules = append(globals.LoadedModules, L)
        }
        // Запускаем инициализацию
        startInitFuncs()
    }
}

func CloseModules() {
    for _, L := range globals.LoadedModules {
        // Закрываем каждый модуль
        L.Close()
    }
}

func startInitFuncs() {
    for _, L := range globals.LoadedModules {
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