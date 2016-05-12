package gae_host

import (
	"github.com/qedus/nds"
	"google.golang.org/appengine/datastore"
	"github.com/strongo/bots-framework/core"
	"golang.org/x/net/context"
)

// Persist user to GAE datastore
type GaeBotUserStore struct {
	GaeBaseStore
	//botUsers 					  map[interface{}]bots.BotUser
	botUserKey                func(botUserID interface{}) *datastore.Key
	validateBotUserEntityType func(entity bots.BotUser)
	newBotUserEntity          func(apiUser bots.WebhookActor) bots.BotUser
	gaeAppUserStore           GaeAppUserStore
}
var _ bots.BotUserStore = (*GaeBotUserStore)(nil) // Check for interface implementation at compile time


// ************************** Implementations of  bots.BotUserStore **************************
func (s GaeBotUserStore) GetBotUserById(botUserId interface{}) (bots.BotUser, error) { // Former LoadBotUserEntity
	//if s.botUsers == nil {
	//	s.botUsers = make(map[int]bots.BotUser, 1)
	//}
	botUserEntity := s.newBotUserEntity(nil)
	ctx := s.Context()
	err := nds.Get(ctx, s.botUserKey(botUserId), botUserEntity)
	return botUserEntity, err
}

func (s GaeBotUserStore) SaveBotUser(botUserID interface{}, userEntity bots.BotUser) error { // Former SaveBotUserEntity
	s.validateBotUserEntityType(userEntity)
	userEntity.SetDtUpdatedToNow()
	_, err := nds.Put(s.Context(), s.botUserKey(botUserID), userEntity)
	return err
}

func (s GaeBotUserStore) CreateBotUser(apiUser bots.WebhookActor) (bots.BotUser, error) {
	s.log.Debugf("CreateBotUser() started...")
	botUserID := apiUser.GetID()
	botUserKey := s.botUserKey(botUserID)
	botUserEntity := s.newBotUserEntity(apiUser)

	err := datastore.RunInTransaction(s.Context(), func(ctx context.Context) error {
		err := nds.Get(ctx, botUserKey, botUserEntity)

		if err == datastore.ErrNoSuchEntity {
			appUserId, _, err := s.gaeAppUserStore.createAppUser(ctx, apiUser)
			if err != nil {
				return err
			}
			botUserEntity.SetAppUserID(appUserId)
			botUserEntity.SetDtUpdatedToNow()
			nds.Put(ctx, botUserKey, botUserEntity)
		} else if err != nil {
			return err
		}

		return nil
	}, &datastore.TransactionOptions{XG: true})

	if err != nil {
		return nil, err
	}
	return botUserEntity, nil
}