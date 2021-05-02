local bot = require("bot")
local http = require("http")
local scrape = require("scrape")
local json = require("json")

pageCount = 0

local typeEx = {
    photo = {
        ".png",
        ".jpg",
        "jpeg",
    },
    doc = {
        ".gif",
    },
}

function onLoaded()
    print("Модуль парсинга Rules34.online запущен")
    local html, err = http.get("https://rules34.online")
    local resp, err = scrape.find_attr_by_tag(html.body, "href", "a")
    pageCount = resp[#resp]:sub(16)
    
    bot.AddCommand("прон", pronCommand)
end

function pronCommand(args, peer_id)
    math.randomseed(os.clock())
    local page = math.random(pageCount-1)
    
    local html, err = http.get("https://rules34.online/?page="..page)
    local resp, err = scrape.find_attr_by_tag(html.body, "href", "a")
    
    math.randomseed(os.clock())
    rand = #resp-21 + math.random(15)
    local url = resp[rand]
    
    local html, err = http.get("https://rules34.online/"..url)
    local resp, err = scrape.find_attr_by_tag(html.body, "src", "img")
    
    local image = "https://rules34.online/" .. resp[3]
    
    attch = uploadImage(image, peer_id)
    bot.Method("messages.send", {
        peer_id = peer_id,
        random_id = 0,
        message = "Наслаждайся",
        attachment = attch,
    })
    
    --bot.MessagesSend(peer_id, "Наслаждайся!", {image})
end

function onClose()
    print("Модуль парсинга Rules34.online вырубается")
end

function uploadImage(image, peer_id)
    ex = image:sub(-4)
    ftype = "photo"
    result = ""
    
    bot.DownloadFile(image, peer_id..ex)
    for i, v in ipairs(typeEx.doc) do
        if v == ex then ftype = "doc" end
    end
    for i, v in ipairs(typeEx.photo) do
        if v == ex then ftype = "photo" end
    end
    
    if ftype == "doc" then
        resp = bot.Method("docs.getMessagesUploadServer", {
            peer_id = peer_id,
            type = "doc",
        })
        resp = json.decode(resp).response
        
        resp = bot.UploadFile(resp.upload_url, peer_id..ex)
        resp = json.decode(resp)
        
        resp = bot.Method("docs.save", {
            file = resp.file,
            title = "huyak",
            tags = "",
        })
        resp = json.decode(resp).response.doc
        result = "doc"..resp.owner_id.."_"..resp.id
    elseif ftype == "photo" then
        resp = bot.Method("photos.getMessagesUploadServer", {
            peer_id = peer_id,
        })
        resp = json.decode(resp).response
        
        resp = bot.UploadFile(resp.upload_url, peer_id..ex)
        resp = json.decode(resp)
        
        resp = bot.Method("photos.saveMessagesPhoto", {
            server = resp.server,
            photo = resp.photo,
            hash = resp.hash,
        })
        resp = json.decode(resp).response[1]
        result = "photo"..resp.owner_id.."_"..resp.id
    end
    os.remove(peer_id..ex)
    return result
end