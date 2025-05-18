package stores

import (
	"context"
	"time"

	"github.com/akatranlp/sentinel/token"
	tokenbasestore "github.com/akatranlp/sentinel/token/base_store"
	"gitlab.hs-flensburg.de/gitlab-classroom/identity_provider/stores/models"
	"gitlab.hs-flensburg.de/gitlab-classroom/identity_provider/stores/query"
	"gorm.io/gorm"
)

type BasicGORMTokenStore struct {
	db *gorm.DB
	q  *query.Query
}

func NewBasicGORMTokenStore(db *gorm.DB) (token.TokenStore, error) {
	repo := &BasicGORMTokenStore{
		db: db,
		q:  query.Use(db),
	}
	return tokenbasestore.NewBaseTokenStore(repo)
}

func (s *BasicGORMTokenStore) GetSessionByID(ctx context.Context, sessionID string) (token.Session, error) {
	querySession := s.q.TokenSession
	sess, err := querySession.WithContext(ctx).
		Where(querySession.SessionID.Eq(sessionID)).
		First()
	if err != nil {
		return token.Session{}, err
	}
	return token.Session{
		SessionID:  sess.SessionID,
		Expiry:     sess.Expiry,
		RefreshJTI: sess.RefreshJTI,
	}, nil
}

func (s *BasicGORMTokenStore) UpdateSession(ctx context.Context, sess token.Session) error {
	querySession := s.q.TokenSession
	_, err := querySession.WithContext(ctx).
		Where(querySession.SessionID.Eq(sess.SessionID)).
		Updates(models.TokenSession{
			Expiry:     sess.Expiry,
			RefreshJTI: sess.RefreshJTI,
		})
	return err
}

func (s *BasicGORMTokenStore) DeleteSessionByID(ctx context.Context, sessionID string) error {
	querySession := s.q.TokenSession
	_, err := querySession.WithContext(ctx).
		Where(querySession.SessionID.Eq(sessionID)).
		Delete()
	return err
}

func (s *BasicGORMTokenStore) SetSession(ctx context.Context, sess token.Session) error {
	querySession := s.q.TokenSession
	err := querySession.WithContext(ctx).
		Save(&models.TokenSession{
			SessionID:  sess.SessionID,
			Expiry:     sess.Expiry,
			RefreshJTI: sess.RefreshJTI,
		})
	return err
}

func (s *BasicGORMTokenStore) DeleteSessionsAfterExpiry(ctx context.Context) error {
	querySession := s.q.TokenSession
	_, err := querySession.WithContext(ctx).
		Where(querySession.Expiry.Lt(time.Now())).
		Delete()
	return err
}
