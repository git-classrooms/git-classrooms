package stores

import (
	"context"
	"log"
	"slices"
	"time"

	"github.com/akatranlp/go-pkg/its"
	"github.com/akatranlp/identity-provider/account"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/identity_provider/stores/models"
	"gitlab.hs-flensburg.de/gitlab-classroom/identity_provider/stores/query"
	"gorm.io/gorm"
)

type GORMUserStore struct {
	db *gorm.DB
	q  *query.Query
}

var _ account.UserStore = (*GORMUserStore)(nil)

func NewGORMUserStore(db *gorm.DB) *GORMUserStore {
	return NewGORMUserStoreWithInterval(db, 5*time.Minute)
}

func NewGORMUserStoreWithInterval(db *gorm.DB, interval time.Duration) *GORMUserStore {
	return &GORMUserStore{db: db, q: query.Use(db)}
}

func (s *GORMUserStore) GetUserByID(ctx context.Context, id account.UserID) (account.User, error) {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return account.User{}, err
	}
	queryUser := s.q.AuthUser
	authUser, err := queryUser.WithContext(ctx).Where(queryUser.ID.Eq(userID)).First()
	if err != nil {
		log.Println("GetUserByID", err)
		return account.User{}, account.ErrUserNotFound
	}
	return authUser.MapToUser(), nil
}

func (s *GORMUserStore) GetUserByAccountID(ctx context.Context, id account.AccountID) (account.User, error) {
	queryAccount := s.q.AuthAccount
	authAccount, err := queryAccount.WithContext(ctx).
		Preload(queryAccount.User).
		Where(
			queryAccount.Provider.Eq(id.Provider),
			queryAccount.ProviderID.Eq(id.ProviderID),
		).
		First()
	if err != nil {
		log.Println("GetUserByAccountID", err)
		return account.User{}, account.ErrAccountNotFound
	}
	return authAccount.User.MapToUser(), nil
}

func (s *GORMUserStore) GetAccountsForUserID(ctx context.Context, id account.UserID) ([]account.Account, error) {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}
	queryUser := s.q.AuthUser
	authUser, err := queryUser.WithContext(ctx).
		Preload(queryUser.Accounts).
		Where(queryUser.ID.Eq(userID)).
		First()
	if err != nil {
		return nil, account.ErrAccountNotFound
	}
	if len(authUser.Accounts) == 0 {
		return nil, account.ErrAccountNotFound
	}
	return slices.Collect(its.Map(
		slices.Values(authUser.Accounts),
		func(a models.AuthAccount) account.Account { return a.MapToAccount() },
	)), nil
}

func (s *GORMUserStore) GetAccountByID(ctx context.Context, accID account.AccountID) (account.Account, error) {
	queryAccount := s.q.AuthAccount
	authAccount, err := queryAccount.WithContext(ctx).Where(
		queryAccount.Provider.Eq(accID.Provider),
		queryAccount.ProviderID.Eq(accID.ProviderID),
	).First()
	if err != nil {
		log.Println("GetAccountByID", err)
		return account.Account{}, account.ErrAccountNotFound
	}
	return authAccount.MapToAccount(), nil
}

func (s *GORMUserStore) GetAccountByProvider(ctx context.Context, id account.UserID, provider string) (account.Account, error) {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return account.Account{}, err
	}
	queryAccount := s.q.AuthAccount
	authAccount, err := queryAccount.WithContext(ctx).
		Preload(queryAccount.User).
		Where(
			queryAccount.Provider.Eq(provider),
			queryAccount.UserID.Eq(userID),
		).First()
	if err != nil {
		log.Println("GetAccountByID", err)
		return account.Account{}, account.ErrAccountNotFound
	}
	return authAccount.MapToAccount(), nil
}

func (s *GORMUserStore) GetOrCreateUserFromAccount(ctx context.Context, acc account.Account) (account.User, error) {
	q := s.q.Begin()
	defer q.Rollback()

	queryAccount := q.AuthAccount
	queryUser := q.AuthUser
	authAccount, err := queryAccount.WithContext(ctx).
		Preload(queryAccount.User).
		Where(
			queryAccount.Provider.Eq(acc.Provider),
			queryAccount.ProviderID.Eq(acc.ProviderID),
		).First()
	if err != nil {
		authUser := &models.AuthUser{
			Name:          acc.Name,
			Email:         acc.Email,
			EmailVerified: acc.EmailVerified,
			Username:      acc.PreferredUsername,
			Picture:       acc.Picture,
		}
		if err = queryUser.WithContext(ctx).Create(authUser); err != nil {
			return account.User{}, err
		}
		authAccount = &models.AuthAccount{
			Provider:          acc.Provider,
			ProviderID:        acc.ProviderID,
			Email:             acc.Email,
			EmailVerified:     acc.EmailVerified,
			AccessToken:       acc.AccessToken,
			Expiry:            acc.Expiry,
			RefreshToken:      acc.RefreshToken,
			RefreshExpiry:     acc.RefreshExpiry,
			TokenType:         acc.TokenType,
			IDToken:           acc.IDToken,
			Name:              acc.Name,
			PreferredUsername: acc.PreferredUsername,
			Nickname:          acc.Nickname,
			Picture:           acc.Picture,
			Profile:           acc.Profile,
			UserID:            authUser.ID,
		}
		if err = queryAccount.WithContext(ctx).Create(authAccount); err != nil {
			return account.User{}, err
		}
		if err = q.Commit(); err != nil {
			return account.User{}, err
		}
		return authUser.MapToUser(), nil
	}

	authAccount.Email = acc.Email
	authAccount.EmailVerified = acc.EmailVerified
	authAccount.AccessToken = acc.AccessToken
	authAccount.Expiry = acc.Expiry
	authAccount.RefreshToken = acc.RefreshToken
	authAccount.RefreshExpiry = acc.RefreshExpiry
	authAccount.TokenType = acc.TokenType
	authAccount.IDToken = acc.IDToken
	authAccount.Name = acc.Name
	authAccount.PreferredUsername = acc.PreferredUsername
	authAccount.Nickname = acc.Nickname
	authAccount.Picture = acc.Picture
	authAccount.Profile = acc.Profile

	if err = queryAccount.WithContext(ctx).Save(authAccount); err != nil {
		return account.User{}, err
	}

	if err = q.Commit(); err != nil {
		return account.User{}, err
	}
	return authAccount.User.MapToUser(), nil
}

func (s *GORMUserStore) UpdateUser(ctx context.Context, id account.UserID, user account.User) error {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return err
	}
	authUser := &models.AuthUser{
		ID:            userID,
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Username:      user.Username,
		Picture:       user.Picture,
	}
	return s.q.AuthUser.WithContext(ctx).Save(authUser)
}

func (s *GORMUserStore) UpdateAccount(ctx context.Context, id account.AccountID, acc account.Account) error {
	authAccount := &models.AuthAccount{
		Email:             acc.Email,
		EmailVerified:     acc.EmailVerified,
		AccessToken:       acc.AccessToken,
		Expiry:            acc.Expiry,
		RefreshToken:      acc.RefreshToken,
		RefreshExpiry:     acc.RefreshExpiry,
		TokenType:         acc.TokenType,
		IDToken:           acc.IDToken,
		Name:              acc.Name,
		PreferredUsername: acc.PreferredUsername,
		Nickname:          acc.Nickname,
		Picture:           acc.Picture,
		Profile:           acc.Profile,
	}
	queryAccount := s.q.AuthAccount
	_, err := queryAccount.WithContext(ctx).Where(
		queryAccount.Provider.Eq(acc.Provider),
		queryAccount.ProviderID.Eq(acc.ProviderID),
	).UpdateColumns(authAccount)
	return err
}

func (s *GORMUserStore) LinkAccount(ctx context.Context, id account.UserID, acc account.Account) error {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return err
	}
	queryUser := s.q.AuthUser
	queryAccount := s.q.AuthAccount
	authUser, err := queryUser.WithContext(ctx).Where(queryUser.ID.Eq(userID)).First()
	if err != nil {
		return account.ErrUserNotFound
	}
	authAccount := &models.AuthAccount{
		Provider:          acc.Provider,
		ProviderID:        acc.ProviderID,
		Email:             acc.Email,
		EmailVerified:     acc.EmailVerified,
		AccessToken:       acc.AccessToken,
		Expiry:            acc.Expiry,
		RefreshToken:      acc.RefreshToken,
		RefreshExpiry:     acc.RefreshExpiry,
		TokenType:         acc.TokenType,
		IDToken:           acc.IDToken,
		Name:              acc.Name,
		PreferredUsername: acc.PreferredUsername,
		Nickname:          acc.Nickname,
		Picture:           acc.Picture,
		Profile:           acc.Profile,
		UserID:            authUser.ID,
	}
	return queryAccount.WithContext(ctx).Create(authAccount)
}

func (s *GORMUserStore) UnLinkAccount(ctx context.Context, id account.UserID, accID account.AccountID) error {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return err
	}
	q := s.q.Begin()
	defer q.Rollback()
	queryUser := q.AuthUser
	queryAccount := q.AuthAccount
	authUser, err := queryUser.WithContext(ctx).
		Preload(queryUser.Accounts).
		Where(queryUser.ID.Eq(userID)).
		First()
	if err != nil {
		return account.ErrUserNotFound
	}

	if len(authUser.Accounts) <= 1 {
		return account.ErrLastAccount
	}

	idx := slices.IndexFunc(authUser.Accounts, func(a models.AuthAccount) bool {
		return a.Provider == accID.Provider && a.ProviderID == accID.Provider
	})
	if idx == -1 {
		return account.ErrAccountNotFound
	}

	acc := authUser.Accounts[idx]

	if _, err = queryAccount.WithContext(ctx).Delete(&acc); err != nil {
		return err
	}

	idx = slices.IndexFunc(authUser.Accounts, func(a models.AuthAccount) bool {
		return a.Provider != accID.Provider || a.ProviderID != accID.Provider
	})
	if idx == -1 {
		panic("unreachable")
	}

	newAcc := authUser.Accounts[idx]

	var userNeedsUpdate bool
	if authUser.Email == acc.Email {
		authUser.Email = newAcc.Email
		userNeedsUpdate = true
	}
	if authUser.Name == acc.Name {
		authUser.Name = newAcc.Name
		userNeedsUpdate = true
	}
	if authUser.Picture == acc.Picture {
		authUser.Picture = newAcc.Picture
		userNeedsUpdate = true
	}
	if authUser.Username == acc.PreferredUsername {
		authUser.Username = newAcc.PreferredUsername
		userNeedsUpdate = true
	}

	if userNeedsUpdate {
		if err := queryUser.WithContext(ctx).Save(authUser); err != nil {
			return err
		}
	}

	if err := q.Commit(); err != nil {
		return err
	}

	return nil
}
