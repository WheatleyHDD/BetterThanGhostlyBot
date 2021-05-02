package botmodule

import (
    "log"
    "strings"
    "strconv"
    "os"
    
    "github.com/yuin/gopher-lua"
    "github.com/SevereCloud/vksdk/v2/api"
    "github.com/WheatleyHDD/BetterThanGhostlyBot/globals"
    "github.com/imroc/req"
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
    "MessagesSend": MessagesSend,
    "AddCommand": AddCommand,
    "Inspect": Inspect,
    "Method": Method,
    "UploadFile": UploadFile,
    "DownloadFile": DownloadFile,
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

func MessagesSend(L *lua.LState) int {
    peerId := L.ToInt(1)
    text := L.ToString(2)
    attachments := L.ToTable(3)
    attch := ""
    
    MediaTypes := [7]string{
        "photo",
        "video",
        "audio",
        "doc",
        "wall",
        "market",
        "poll",
    }
    
    attachments.ForEach(func(i, v lua.LValue) {
        if strings.HasPrefix(v.String(), "https://") || strings.HasPrefix(v.String(), "http://") {
            isF, ft := getFileType(v.String())
            if isF {
                url := v.String()
                r, _ := req.Get(url)
                
                sf := strconv.Itoa(peerId) + url[len(url)-4:]
                
                log.Println(url)
                log.Println(sf)
                
                r.ToFile(sf)
                file, _ := os.Open(sf)
                
                switch ft {
                    case "photo":
                        photosPhoto, err := globals.VK.UploadMessagesPhoto(peerId, file)
                        if err != nil {
                            log.Println(err)
                        }
                        attch = attch + ft + strconv.Itoa(photosPhoto[0].OwnerID) + "_" + strconv.Itoa(photosPhoto[0].ID) + ","
                    case "doc":
                        photosPhoto, err := globals.VK.UploadMessagesDoc(peerId, "doc", "test", "", file)
                        if err != nil {
                            log.Println(err)
                        }
                        attch = attch + ft + strconv.Itoa(photosPhoto.Doc.OwnerID) + "_" + strconv.Itoa(photosPhoto.Doc.ID) + ","
                    case "video":
                        photosPhoto, err := globals.VK.UploadVideo(api.Params{
                            "is_private": true,
                        }, file)
                        if err != nil {
                            log.Println(err)
                        }
                        attch = attch + ft + strconv.Itoa(photosPhoto.OwnerID) + "_" + strconv.Itoa(photosPhoto.VideoID) + ","
                    case "audio":
                        photosPhoto, err := globals.VK.UploadMessagesDoc(peerId, "audio_message", "what?", "", file)
                        if err != nil {
                            log.Println(err)
                        }
                        attch = attch + "doc" + strconv.Itoa(photosPhoto.AudioMessage.OwnerID) + "_" + strconv.Itoa(photosPhoto.AudioMessage.ID) + ","
                }
                //_ = os.Remove(sf)
            } else {
                attch = attch + v.String() + ","
            }
            
        } else {
            already := false
            for _, mtype := range MediaTypes {
                if strings.HasPrefix(v.String(), mtype) {
                    already = true
                    attch = attch + v.String() + ","
                    break
                }
            }
            if !already {
                _, ft := getFileType(v.String())
                file, err := os.Open(v.String())
                if err != nil {
                    log.Println(err)
                }
                switch ft {
                    case "photo":
                        photosPhoto, err := globals.VK.UploadMessagesPhoto(peerId, file)
                        if err != nil {
                            log.Println(err)
                        }
                        attch = attch + ft + strconv.Itoa(photosPhoto[0].OwnerID) + "_" + strconv.Itoa(photosPhoto[0].ID) + ","
                    case "doc":
                        photosPhoto, err := globals.VK.UploadMessagesDoc(peerId, "doc", "test", "", file)
                        if err != nil {
                            log.Println(err)
                        }
                        attch = attch + ft + strconv.Itoa(photosPhoto.Doc.OwnerID) + "_" + strconv.Itoa(photosPhoto.Doc.ID) + ","
                    case "video":
                        photosPhoto, err := globals.VK.UploadVideo(api.Params{
                            "is_private": true,
                        }, file)
                        if err != nil {
                            log.Println(err)
                        }
                        attch = attch + ft + strconv.Itoa(photosPhoto.OwnerID) + "_" + strconv.Itoa(photosPhoto.VideoID) + ","
                    case "audio":
                        photosPhoto, err := globals.VK.UploadMessagesDoc(peerId, "audio_message", "what?", "", file)
                        if err != nil {
                            log.Println(err)
                        }
                        attch = attch + "doc" + strconv.Itoa(photosPhoto.AudioMessage.OwnerID) + "_" + strconv.Itoa(photosPhoto.AudioMessage.ID) + ","
                }
            }
        }
    })
    
    resp, _ := globals.VK.MessagesSend(api.Params{
		"peer_id":    peerId,
		"message":    text,
		"attachment": attch,
		"random_id":  0,
	})
    
    L.Push(lua.LNumber(resp))
    return 1
}

func getFileType(file string) (isFile bool, fileType string) {
    ft := ""
    isF := false
    for _, t := range photoEx {
        if strings.HasSuffix(file, t) {
            ft = "photo"
            isF = true
        }
    }
    if strings.HasSuffix(file, "mp4") {
        ft = "video"
        isF = true
    } else if strings.HasSuffix(file, "gif") {
        ft = "doc"
        isF = true
    } else if strings.HasSuffix(file, "ogg") {
        ft = "audio"
        isF = true
    }
    return isF, ft
    
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