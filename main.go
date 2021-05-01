package main

import(
    "log"
    "strings"
    
    "github.com/WheatleyHDD/BetterThanGhostlyBot/utils"
    
    "github.com/yuin/gopher-lua"
    longpoll "github.com/SevereCloud/vksdk/v2/longpoll-user"
    "github.com/SevereCloud/vksdk/v2/longpoll-user/v3"
    "github.com/SevereCloud/vksdk/v2/api"
    
	//"github.com/SevereCloud/vksdk/v2/api/params"
)

var (
    vk *api.VK
)

func main() {
    log.Println("Загрузка конфигов...")
    LoadConfig()
    log.Println("Загрузка модулей...")
    LoadModules()
    log.Println(LoadedModules)
    log.Println(AllCmds)
    log.Println("Модули загружены")
    
    vk = api.NewVK(AccessToken)
    
    mode := longpoll.ReceiveAttachments + longpoll.ExtendedEvents
    lp, err := longpoll.NewLongPoll(vk, mode)
    
    if err != nil {
        panic(err)
    }
    
    w := wrapper.NewWrapper(lp)

    // event with code 4
    w.OnNewMessage(OnMessage)
    
    if err := lp.Run(); err != nil {
		log.Fatal(err)
	}

	// Безопасное завершение
	// Ждет пока соединение закроется и события обработаются
	lp.Shutdown()

	// Закрыть соединение
	// Требует lp.Client.Transport = &http.Transport{DisableKeepAlives: true}
	lp.Client.CloseIdleConnections()
    
    log.Println("Отключение модулей...")
    CloseModules()
    log.Println("Модули отключены")
    log.Println("Пока :(")
}

func OnMessage(m wrapper.NewMessage) {
    mText := strings.ToLower(m.Text)
    for _, a := range Appeals {
        if strings.HasPrefix(mText, strings.ToLower(a.(string)) + ", ") {
           OnMessageToBot(m)
           break
        }
    }
}

func OnMessageToBot(m wrapper.NewMessage) {
    args := strings.Split(strings.ToLower(m.Text), " ")
    args = utils.RemoveItemString(args, 0)
    for _, cmd := range AllCmds {
        if args[0] == strings.ToLower(cmd.Cmd) {
            argTable := cmd.Module.NewTable()
            for _, arg := range args {
                argTable.Append(lua.LString(arg))
            }
            if err := cmd.Module.CallByParam(lua.P{
				Fn:      cmd.Fn,
				NRet:    0,
				Protect: true,
			}, argTable, lua.LNumber(m.PeerID)); err != nil {
				panic(err)
			}
			break
        }
    }
}