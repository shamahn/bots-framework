package gae_host

import (
	"errors"
	"fmt"
	"github.com/qedus/nds"
	"github.com/strongo/app"
	"github.com/strongo/bots-framework/core"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"net/http"
	"reflect"
)

type GaeAppUserStore struct {
	appUserEntityKind string
	appUserEntityType reflect.Type
	newUserEntity     func() bots.BotAppUser
	GaeBaseStore
}

var _ bots.BotAppUserStore = (*GaeAppUserStore)(nil)

func NewGaeAppUserStore(log strongo.Logger, r *http.Request, appUserEntityKind string, appUserEntityType reflect.Type, newUserEntity func() bots.BotAppUser) GaeAppUserStore {
	return GaeAppUserStore{
		appUserEntityType: appUserEntityType,
		appUserEntityKind: appUserEntityKind,
		newUserEntity:     newUserEntity,
		GaeBaseStore:      NewGaeBaseStore(log, r, appUserEntityKind),
	}
}

// ************************** Helper functions **************************

func (s GaeAppUserStore) appUserKey(appUserId int64) *datastore.Key {
	return datastore.NewKey(s.Context(), s.appUserEntityKind, "", appUserId, nil)
}

// ************************** Implementations of  bots.AppUserStore **************************
func (s GaeAppUserStore) GetAppUserByID(appUserId int64, appUser bots.BotAppUser) error {
	return nds.Get(s.Context(), s.appUserKey(appUserId), appUser)
}

func (s GaeAppUserStore) CreateAppUser(actor bots.WebhookActor) (int64, bots.BotAppUser, error) {
	return s.createAppUser(s.Context(), actor)
}

func (s GaeAppUserStore) createAppUser(c context.Context, actor bots.WebhookActor) (int64, bots.BotAppUser, error) {
	appUserEntity := s.newUserEntity()
	appUserEntity.SetBotUserID(actor.Platform(), actor.GetID())
	appUserEntity.SetNames(actor.GetFirstName(), actor.GetLastName(), actor.GetUserName())
	key, err := nds.Put(c, s.appUserKey(0), appUserEntity)
	return key.IntID(), appUserEntity, err
}

func (s GaeAppUserStore) getAppUserIdByBotUserKey(c context.Context, botUserKey *datastore.Key) (int64, error) {
	query := datastore.NewQuery(s.appUserEntityKind).Filter("TelegramUserIDs =", botUserKey.IntID()).KeysOnly().Limit(2)
	//appUsers := reflect.MakeSlice(reflect.SliceOf(s.appUserEntityType), 0, 2)
	keys, err := query.GetAll(c, nil)
	if err != nil {
		s.logger.Errorf(c, "Failed to query app users by TelegramUserIDs: %v", err)
		return 0, err
	}
	switch len(keys) {
	case 0:
		return 0, nil
	case 1:
		return keys[0].IntID(), nil
	default:
		return 0, errors.New(fmt.Sprintf("Found few app users by %v", botUserKey))
	}
}

func (s GaeAppUserStore) SaveAppUser(appUserId int64, appUserEntity bots.BotAppUser) error {
	if appUserId == 0 {
		panic("appUserId == 0")
	}
	_, err := nds.Put(s.Context(), s.appUserKey(appUserId), appUserEntity)
	return err
}
