package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/line/line-bot-sdk-go/linebot"
)

var messagingSettings *LineMessagingSettings
var settingsOnce sync.Once

type LineMessagingSettings struct {
	APISecret          string `json:"APISecret"`
	ChannelAccessToken string `json:"channelAccessToken"`
	bot                *linebot.Client
}

func NewLineMessagingSetting(settingsFilePath string) *LineMessagingSettings {
	var settings LineMessagingSettings
	fileData, err := ioutil.ReadFile(settingsFilePath)
	panicError(err)
	err = json.Unmarshal(fileData, &settings)
	panicError(err)
	settings.bot, err = linebot.New(settings.APISecret, settings.ChannelAccessToken)
	panicError(err)
	return &settings
}

func GetLineMessagingSettings() *LineMessagingSettings {
	settingsOnce.Do(func() {
		messagingSettings = NewLineMessagingSetting(settingsFilePath)
	})
	return messagingSettings
}
