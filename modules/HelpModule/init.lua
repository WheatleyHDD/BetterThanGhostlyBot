local bot = require("bot")

function onLoaded()
  print("Модуль помощи запущен")

  bot.AddCommand("помощь", startCommand, "Все команды", "🚑")
end

function startCommand(args, peer_id)
  local allCommands = bot.GetAllCommands()

  local text = "Все команды:\n"

  for key, val in pairs(allCommands) do
    text = text .. val.Icon .. " " .. val.Cmd .. " - " .. val.Description .. "\n"
  end

  bot.Method("messages.send", {
    peer_id = peer_id,
    random_id = 0,
    message = text,
  })
end

function onClose()
  print("Модуль помощи вырубается")
end
