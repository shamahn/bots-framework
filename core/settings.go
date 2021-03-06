package bots

import (
	"github.com/strongo/app"
	"golang.org/x/net/context"
)

type BotMode int8

const (
	Production BotMode = iota
	Staging
	Development
	Local
)

type BotSettings struct {
	Mode        BotMode
	Code        string
	Token       string
	VerifyToken string // Used by Facebook
	Locale      strongo.Locale
}

func NewBotSettings(mode BotMode, code, token string, locale strongo.Locale) BotSettings {
	if code == "" {
		panic("Missing required parameter: code")
	}
	if token == "" {
		panic("Missing required parameter: token")
	}
	if locale.Code5 == "" {
		panic("Missing required parameter: locale.Code5")
	}
	return BotSettings{
		Code:   code,
		Mode:   mode,
		Token:  token,
		Locale: locale,
	}
}

type BotSettingsProvider func(c context.Context) BotSettingsBy

type BotSettingsBy struct { // TODO: Decide if it should have map[string]*BotSettings instead of map[string]BotSettings
	Code     map[string]BotSettings
	ApiToken map[string]BotSettings
	Locale   map[string]BotSettings
}

func NewBotSettingsBy(bots ...BotSettings) BotSettingsBy {
	count := len(bots)
	botsBy := BotSettingsBy{
		Code:     make(map[string]BotSettings, count),
		ApiToken: make(map[string]BotSettings, count),
		Locale:   make(map[string]BotSettings, count),
	}
	for _, bot := range bots {
		botsBy.Code[bot.Code] = bot
		botsBy.ApiToken[bot.Token] = bot
		botsBy.Locale[bot.Locale.Code5] = bot
	}
	return botsBy
}
