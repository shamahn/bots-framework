package gae_host

import (
	"github.com/qedus/nds"
	"github.com/strongo/bots-framework/core"
	"google.golang.org/appengine/datastore"
)

type EntityTypeValidator interface {
}

// Persist chat to GAE datastore
type GaeBotChatStore struct {
	GaeBaseStore
	botChats                  map[interface{}]bots.BotChat
	botChatKey                func(botChatId interface{}) *datastore.Key
	validateBotChatEntityType func(entity bots.BotChat)
	newBotChatEntity          func() bots.BotChat
}

var _ bots.BotChatStore = (*GaeBotChatStore)(nil) // Check for interface implementation at compile time

// ************************** Implementations of  bots.ChatStore **************************
func (s *GaeBotChatStore) GetBotChatEntityById(botChatId interface{}) (bots.BotChat, error) { // Former LoadBotChatEntity
	//s.logger.Debugf(s.Context(), "GaeBotChatStore.GetBotChatEntityById(%v)", botChatId)
	if s.botChats == nil {
		s.botChats = make(map[interface{}]bots.BotChat, 1)
	}
	botChatEntity := s.newBotChatEntity()
	err := nds.Get(s.Context(), s.botChatKey(botChatId), botChatEntity)
	if err != nil {
		s.logger.Infof(s.Context(), "Failed to get bot chat entity by ID: %v - %T(%v)", botChatId, err, err)
		if err == datastore.ErrNoSuchEntity {
			return nil, bots.ErrEntityNotFound
		}
	} else {
		s.botChats[botChatId] = botChatEntity
	}
	return botChatEntity, err
}

func (s *GaeBotChatStore) SaveBotChat(chatId interface{}, chatEntity bots.BotChat) error { // Former SaveBotChatEntity
	s.validateBotChatEntityType(chatEntity)
	chatEntity.SetDtUpdatedToNow()
	_, err := nds.Put(s.Context(), s.botChatKey(chatId), chatEntity)
	return err
}

func (s *GaeBotChatStore) NewBotChatEntity(botID string, botChatId interface{}, appUserID int64, botUserID interface{}, isAccessGranted bool) bots.BotChat {
	s.logger.Debugf(s.Context(), "NewBotChatEntity(botID=%v, botChatId=%v, appUserID=%v, botUserID=%v, isAccessGranted=%v)", botID, botChatId, appUserID, botUserID, isAccessGranted)
	botChat := s.newBotChatEntity()
	botChat.SetAppUserIntID(appUserID)
	botChat.SetBotUserID(botUserID)
	botChat.SetAccessGranted(isAccessGranted)
	botChat.SetBotID(botID)
	s.botChats[botChatId] = botChat
	return botChat
}

func (s *GaeBotChatStore) Close() error { // Former SaveBotChatEntity
	if len(s.botChats) == 0 {
		//s.logger.Debugf(s.Context(), "GaeBotChatStore.Close(): Nothing to save")
		return nil
	}
	//s.logger.Debugf(s.Context(), "GaeBotChatStore.Close(): %v entities to save", len(s.botChats))
	var chatKeys []*datastore.Key
	var chatEntities []bots.BotChat
	for chatId, chatEntity := range s.botChats {
		s.validateBotChatEntityType(chatEntity)
		chatEntity.SetDtUpdatedToNow()
		chatEntity.SetDtLastInteractionToNow()
		chatKeys = append(chatKeys, s.botChatKey(chatId))
		chatEntities = append(chatEntities, chatEntity)
	}
	_, err := nds.PutMulti(s.Context(), chatKeys, chatEntities)
	if err == nil {
		s.logger.Infof(s.Context(), "Succesfully saved %v BotChat entities with keys: %v", len(chatKeys), chatKeys)
		s.botChats = nil
	} else {
		s.logger.Errorf(s.Context(), "Failed to save %v BotChat entities: %v", len(chatKeys), err)
	}
	return err
}
