local bot = require("bot")

function onLoaded()
    print("Тестовый модуль запущен")
    bot.AddCommand("старт", startCommand)
end

function startCommand(args, peer_id)
    print(args)
    bot.MessagesSend(peer_id, "хуярт", {})
end

function onClose()
    
end