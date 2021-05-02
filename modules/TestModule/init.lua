local bot = require("bot")

function onLoaded()
    print("Тестовый модуль запущен")
    
    bot.AddCommand("тест", startCommand)
end

function startCommand(args, peer_id)
    bot.Method("messages.send", {
        peer_id = peer_id,
        random_id = 0,
        message = "Все зашибись!",
    })
end

function onClose()
    print("Тестовый модуль вырубается")
end