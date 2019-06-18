package LineMessagingSettings

import (
	"CexOutOfStock/crash"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/line/line-bot-sdk-go/linebot"
)

const settingsFilePath = "./settings.json"

var messagingSettings *LineMessagingSettings
var settingsOnce sync.Once

type LineMessagingSettings struct {
	APISecret          string `json:"APISecret"`
	ChannelAccessToken string `json:"channelAccessToken"`
	Port               int    `json:"port"`
	CertFile           string `json:"certFile"`
	KeyFile            string `json:"keyFile"`
	IsTLS              bool
	Bot                *linebot.Client
}

func Settings(settingsFilePath string) *LineMessagingSettings {
	var settings LineMessagingSettings
	fileData, err := ioutil.ReadFile(settingsFilePath)
	crash.PanicError(err)
	err = json.Unmarshal(fileData, &settings)
	crash.PanicError(err)
	settings.Bot, err = linebot.New(settings.APISecret, settings.ChannelAccessToken)
	crash.PanicError(err)
	return &settings
}

func SettingsEnv() *LineMessagingSettings {
	var settings LineMessagingSettings
	settings.ChannelAccessToken = checkIfblank("channel")
	settings.APISecret = checkIfblank("secret")
	portString := checkIfblank("PORT")
	var err error
	if settings.Port, err = strconv.Atoi(portString); err != nil {
		crash.PanicError(err)
	}
	settings.Bot, err = linebot.New(settings.APISecret, settings.ChannelAccessToken)
	crash.PanicError(err)
	return &settings
}

func GetLineMessagingSettings() *LineMessagingSettings {
	settingsOnce.Do(func() {
		if _, err := os.Stat(settingsFilePath); err == nil {
			messagingSettings = Settings(settingsFilePath)
		} else {
			messagingSettings = SettingsEnv()
		}
		messagingSettings.IsTLS = messagingSettings.CertFile != ""
	})
	return messagingSettings
}

func checkIfblank(env string) string {
	envValue := ""
	if envValue = os.Getenv(env); envValue == "" {
		crash.PanicError(errors.New(env + " not set as enviroment variable"))
	}
	return envValue
}
