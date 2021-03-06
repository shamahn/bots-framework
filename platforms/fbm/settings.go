package fbm_strongo_bot

import (
	"github.com/strongo/app"
	"github.com/strongo/bots-framework/core"
)

func NewFbmBot(mode bots.BotMode, code, token, verifyToken string, locale strongo.Locale) bots.BotSettings {
	botSettings := bots.NewBotSettings(mode, code, token, locale)
	botSettings.VerifyToken = verifyToken
	return botSettings
}
