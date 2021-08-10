local bot = require("bot")
local http = require("http")
local json = require("json")

local req_url = "https://r34-json.herokuapp.com"

function onLoaded()
  print("R34Module запущен")

  bot.AddCommand("r34", r34Command)
end

function r34Command(args, peer_id)
  tags = ""
  if #args > 1 then
    for key, val in pairs(args) do
      if key ~= 1 then
        if key == 2 then
          tags = val
        else
          tags = tags .. "+" .. val
        end
      end
    end
  end

  print("[SUCCESS] R34Module: Получен запрос на картинку по тегам: " .. tags)

  local resp = http.get(req_url .. "/posts", {
    query = "limit=1&tags=" .. tags,
  })
  local resp_json = json.decode(resp.body)
  local count_on_tag = resp_json.count

  if tonumber(count_on_tag) > 2000 then count_on_tag = 2000 end
  rand = randNumber(count_on_tag)

  local resp = http.get(req_url .. "/posts", {
    query = "limit=1&pid=".. rand .."&tags=" .. tags,
  })
  local resp_json = json.decode(resp.body)

  local image = resp_json.posts[1]["file_url"]
  local file_tags = resp_json.posts[1].tags

  bot.DownloadFile(image, peer_id .. ".png")

  local gmus = bot.Method("photos.getMessagesUploadServer", {})
  local upload_server = json.decode(gmus).response.upload_url

  local upload_data = bot.UploadFile(upload_server, peer_id .. ".png")
  local ud = json.decode(upload_data)
  local saved = bot.Method("photos.saveMessagesPhoto", {
    photo = ud.photo,
    server = ud.server,
    hash = ud.hash,
  })

  os.remove(peer_id .. ".png")

  local saved_json = json.decode(saved)

  local attach = "photo" .. saved_json.response[1].owner_id .. "_" .. saved_json.response[1].id

  local author = "Автор не найден"

  for key, val in pairs(file_tags) do
    local resp = http.get(req_url .. "/tags", {
      query = "type=artist&name=" .. val,
    })
    local resp_json = json.decode(resp.body)

    if #resp_json ~= 0 then
      author = "Автор: " .. val
      break
    end
  end

  local vk_resp = bot.Method("messages.send", {
    peer_id = peer_id,
    random_id = 0,
    message = author,
    attachment = attach,
  })
  local vk_resp = json.decode(vk_resp)
  if vk_resp.error ~= nil then
    print("[ERROR] R34Module: " .. vk_resp.error.error_msg)
  else
    print("[SUCCESS] R34Module: Отправлено!")
  end
end

function onClose()
  print("R34Module вырубается")
end

function randNumber(n)
  math.randomseed(os.time())
  return math.random(n+1)-1
end
