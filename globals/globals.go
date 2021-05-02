package globals

import (
    "github.com/SevereCloud/vksdk/v2/api"
    
    "github.com/yuin/gopher-lua"
)

type Command struct {
    Cmd string
    Fn *lua.LFunction
    Module *lua.LState
}

var(
    AccessToken string
    VK *api.VK
    LoadedModules []*lua.LState
    AllCmds []*Command
)