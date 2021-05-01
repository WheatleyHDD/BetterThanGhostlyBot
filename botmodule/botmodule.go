package botmodule

import (
    "log"
    
    "github.com/yuin/gopher-lua"
    
    "github.com/SevereCloud/vksdk/v2/api"
    
    "github.com/WheatleyHDD/BetterThanGhostlyBot/globals"
)

func Loader(L *lua.LState) int {
    // register functions to the table
    mod := L.SetFuncs(L.NewTable(), exports)

    // returns the module
    L.Push(mod)
    return 1
}

var exports = map[string]lua.LGFunction{
    "MessagesSend": MessagesSend,
    "AddCommand": AddCommand,
}

func MessagesSend(L *lua.LState) int {
    peerId := L.ToInt(1)
    text := L.ToString(2)
    
    resp, _ := globals.VK.MessagesSend(api.Params{
		"peer_id":    peerId,
		"message":    text,
		"attachment": "",
		"random_id":  0,
	})
    
    L.Push(lua.LNumber(resp))
    return 1
}

func AddCommand(L *lua.LState) int {
    cmd := L.ToString(1)
    fn := L.ToFunction(2)
    has := false
    for _, cmmnd := range globals.AllCmds {
        if cmmnd.Cmd == cmd {
            has = true
            break
        }
    }
    if !has {
        globals.AllCmds = append(globals.AllCmds, &globals.Command{
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