package stores

import (
	"context"
	"log"
	"time"

	"github.com/akatranlp/identity-provider/token"
	"gitlab.hs-flensburg.de/gitlab-classroom/identity_provider/stores/models"
	"gitlab.hs-flensburg.de/gitlab-classroom/identity_provider/stores/query"
	"gorm.io/gorm"
)

type GORMTokenStore struct {
	db *gorm.DB
	q  *query.Query
}

var _ token.TokenStore = (*GORMTokenStore)(nil)

func NewGORMTokenStore(db *gorm.DB) *GORMTokenStore {
	return NewGORMTokenStoreWithInterval(db, 5*time.Minute)
}

func NewGORMTokenStoreWithInterval(db *gorm.DB, interval time.Duration) *GORMTokenStore {
	store := &GORMTokenStore{db: db, q: query.Use(db)}
	go store.startCleanup(interval)
	return store
}

func (s *GORMTokenStore) SetSession(ctx context.Context, sid string, jti string, expiry time.Time) error {
	return s.q.TokenSession.WithContext(ctx).Save(&models.TokenSession{
		SessionID:  sid,
		RefreshJTI: jti,
		Expiry:     expiry,
	})
}

func (s *GORMTokenStore) GetSession(ctx context.Context, sid string) (token.Session, error) {
	querySession := s.q.TokenSession
	t := &models.TokenSession{}
	t, err := querySession.WithContext(ctx).
		Where(querySession.SessionID.Eq(sid)).
		First()
	if err != nil {
		log.Println(err)
		return token.Session{}, token.ErrSessionNotFound
	}
	if t.Expiry.Before(time.Now()) {
		_ = s.RevokeSession(ctx, sid)
		return token.Session{}, token.ErrSessionNotFound
	}

	return token.Session{
		SessionID:  t.SessionID,
		RefreshJTI: t.RefreshJTI,
		Expiry:     t.Expiry,
	}, nil
}

func (s *GORMTokenStore) RevokeSession(ctx context.Context, sid string) error {
	querySession := s.q.TokenSession
	_, err := querySession.WithContext(ctx).
		Where(querySession.SessionID.Eq(sid)).
		Delete()
	return err
}

func (s *GORMTokenStore) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		err := s.deleteExpired()
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *GORMTokenStore) deleteExpired() error {
	querySession := s.q.TokenSession
	_, err := querySession.WithContext(context.Background()).
		Where(querySession.Expiry.Lt(time.Now())).
		Delete()
	return err
}
