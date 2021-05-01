package main

import (
    "log"
    
    "github.com/pelletier/go-toml"
)

var(
    AccessToken string
    BotSettings *toml.Tree
    Appeals []interface{}
)

func LoadConfig() {
    
    config, err := toml.LoadFile("conf.toml")
	if err != nil {
	    log.Panic(err)
	}
	
	AccessToken = config.Get("account.access_token").(string)
	BotSettings = config.Get("bot_settings").(*toml.Tree)
	Appeals = config.Get("bot_settings.appeal").([]interface{})
}