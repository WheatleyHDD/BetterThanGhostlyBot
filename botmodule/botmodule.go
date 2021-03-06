package botmodule

import (
    "log"
    "os"

    "github.com/yuin/gopher-lua"
    "github.com/WheatleyHDD/BetterThanGhostlyBot/globals"
    "github.com/imroc/req"
    "layeh.com/gopher-luar"
)

var (
    photoEx = [3]string{
        "png",
        "jpg",
        "jpeg",
    }
)

func Loader(L *lua.LState) int {
    // register functions to the table
    mod := L.SetFuncs(L.NewTable(), exports)

    // returns the module
    L.Push(mod)
    return 1
}

var exports = map[string]lua.LGFunction{
    "AddCommand": AddCommand,
    "Inspect": Inspect,
    "Method": Method,
    "UploadFile": UploadFile,
    "DownloadFile": DownloadFile,
    "GetAllCommands": GetAllCommands,
}

func Inspect(L *lua.LState) int {
    ud := L.ToUserData(1)

    log.Println(ud.Value)

    return 0
}

func DownloadFile(L *lua.LState) int {
    filename := L.ToString(2)
    url := L.ToString(1)

    r, _ := req.Get(url)
    r.ToFile(filename)

    return 0
}

func UploadFile(L *lua.LState) int {
    filename := L.ToString(2)
    url := L.ToString(1)

    file, _ := os.Open(filename)
    r, _ := req.Post(url, req.FileUpload{
    	File:      file,
    	FieldName: "file",       // FieldName is form field name
    	FileName:  filename, //Filename is the name of the file that you wish to upload. We use this to guess the mimetype as well as pass it onto the server
    })
    body, err := r.ToString()
    if err != nil {
        log.Println(err)
    }

    L.Push(lua.LString(body))
    return 1
}

func Method(L *lua.LState) int {
    method := L.ToString(1)
    params := L.ToTable(2)

    url := "https://api.vk.com/method/"+method
    formParams := "access_token="+globals.AccessToken

    params.ForEach(func(i, v lua.LValue) {
        formParams = formParams + "&" + i.String() + "=" + v.String()
    })

    formParams = formParams + "&v=5.130"

    r, _ := req.Post(url, formParams)
    body, err := r.ToString()
    if err != nil {
        log.Println(err)
    }

    L.Push(lua.LString(body))
    return 1
}

func AddCommand(L *lua.LState) int {
    cmd := L.ToString(1)
    fn := L.ToFunction(2)
    desc := L.ToString(3)
    icon := L.ToString(4)

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
            Description: desc,
            Icon: icon,
            Fn: fn,
            Module: L,
        })
    } else {
        log.Println("?????????????? \"" + cmd + "\" ?????? ????????????????????" )
    }
    L.Push(lua.LBool(has))
    return 1
}

func GetAllCommands(L *lua.LState) int {
    t := L.NewTable()
    for _, cmd := range globals.AllCmds {
        t.Append(luar.New(L, cmd))
    }
    L.Push(t)
    return 1
}
