package main

import(
    "log"
    "strings"
    "time"

    "github.com/ztrue/shutdown"
    "os"

    // "github.com/WheatleyHDD/BetterThanGhostlyBot/utils"
    "github.com/WheatleyHDD/BetterThanGhostlyBot/globals"

    "github.com/yuin/gopher-lua"
    longpoll "github.com/SevereCloud/vksdk/v2/longpoll-user"
    "github.com/SevereCloud/vksdk/v2/longpoll-user/v3"
    "github.com/SevereCloud/vksdk/v2/api"

)

var (
    lp *longpoll.LongPoll
)

func main() {

    log.Println("Загрузка конфигов...")
    LoadConfig()
    globals.VK = api.NewVK(globals.AccessToken)

    log.Println("Загрузка модулей...")
    LoadModules()
    log.Println("Модули загружены")

    go GoToOnline()

    go StartLongPoll()

    shutdown.Listen()
}

func StartLongPoll() {
    mode := longpoll.ReceiveAttachments + longpoll.ExtendedEvents
    lp, err := longpoll.NewLongPoll(globals.VK, mode)
    if err != nil {
        panic(err)
    }


    shutdown.Add(func() {
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
      os.Exit(1)
  	})

    w := wrapper.NewWrapper(lp)

    // event with code 4
    w.OnNewMessage(OnMessage)

    log.Println("Лонгпул запущен")

    if err := lp.Run(); err != nil {
  		StartLongPoll()
  	}
}


func GoToOnline() {
    for {
        globals.VK.AccountSetOnline(api.Params{
            "voip": 0,
        })
        time.Sleep(time.Minute * 5)
    }
}

func OnMessage(m wrapper.NewMessage) {
    mText := strings.ToLower(m.Text)
    for _, a := range Appeals {
        if strings.HasPrefix(mText, strings.ToLower(a.(string))) {
           OnMessageToBot(m, strings.ToLower(a.(string)))
           break
        }
    }
}

func OnMessageToBot(m wrapper.NewMessage, appeal string) {
    rawText := strings.Replace(strings.ToLower(m.Text), appeal, "", 1)
    args := strings.Split(strings.ToLower(rawText), " ")
    //args = utils.RemoveItemString(args, 0)
    for _, cmd := range globals.AllCmds {
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
