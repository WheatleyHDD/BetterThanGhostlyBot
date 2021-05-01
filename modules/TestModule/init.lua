function onLoaded()
    print("Тестовый модуль запущен")
    botAddCommand("старт", startCommand)
end

function startCommand(args, peer_id)
    print(args)
    botSendMessage(peer_id, "хуярт", {})
end

function onClose()
    
end