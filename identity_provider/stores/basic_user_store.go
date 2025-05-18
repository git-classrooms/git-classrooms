package stores

import (
	"context"
	"slices"

	"github.com/akatranlp/go-pkg/its"
	"github.com/akatranlp/sentinel/account"
	userbasestore "github.com/akatranlp/sentinel/account/base_store"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/identity_provider/stores/models"
	"gitlab.hs-flensburg.de/gitlab-classroom/identity_provider/stores/query"
	"gorm.io/gorm"
)

type BasicGORMUserStore struct {
	db *gorm.DB
	q  *query.Query
}

var _ userbasestore.Repository = (*BasicGORMUserStore)(nil)
var _ userbasestore.UserIDAndPRoviderGetter = (*BasicGORMUserStore)(nil)

func NewBasicGORMUserStore(db *gorm.DB) account.UserStore {
	repo := &BasicGORMUserStore{
		db: db,
		q:  query.Use(db),
	}
	return userbasestore.NewBaseUserStore(repo)
}

func (s *BasicGORMUserStore) CreateUser(ctx context.Context, user account.User) (account.User, error) {
	authUser := models.AuthUser{
		Name:          user.Name,
		Username:      user.Username,
		Picture:       user.Picture,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	}
	if err := s.q.AuthUser.WithContext(ctx).Save(&authUser); err != nil {
		return account.User{}, err
	}

	return authUser.MapToUser(), nil
}

func (s *BasicGORMUserStore) GetUserByID(ctx context.Context, id account.UserID) (account.User, error) {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return account.User{}, err
	}
	queryUser := s.q.AuthUser
	user, err := queryUser.WithContext(ctx).
		Where(queryUser.ID.Eq(userID)).
		First()
	if err != nil {
		// TODO: check if it is not found err
		return account.User{}, account.ErrUserNotFound
	}
	return user.MapToUser(), nil
}

func (s *BasicGORMUserStore) GetUserByAccountID(ctx context.Context, id account.AccountID) (account.User, error) {
	queryAccount := s.q.AuthAccount
	authAcc, err := queryAccount.WithContext(ctx).
		Preload(queryAccount.User).
		Where(
			queryAccount.Provider.Eq(id.Provider),
			queryAccount.ProviderID.Eq(id.ProviderID),
		).
		First()
	if err != nil {
		// TODO: check if it is not found err
		return account.User{}, account.ErrUserNotFound
	}
	return authAcc.User.MapToUser(), nil
}

func (s *BasicGORMUserStore) UpdateUser(ctx context.Context, id account.UserID, user account.User) error {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return err
	}
	queryUser := s.q.AuthUser
	_, err = queryUser.WithContext(ctx).
		Where(queryUser.ID.Eq(userID)).
		Updates(models.AuthUser{
			Name:          user.Name,
			Username:      user.Username,
			Picture:       user.Picture,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
		})
	return err
}

func (s *BasicGORMUserStore) CreateAccount(ctx context.Context, acc account.Account) (account.Account, error) {
	userID, err := uuid.Parse(string(acc.UserID))
	if err != nil {
		return account.Account{}, err
	}
	authAcc := models.AuthAccount{
		Provider:          acc.Provider,
		ProviderID:        acc.ProviderID,
		AccessToken:       acc.AccessToken,
		Expiry:            acc.Expiry,
		RefreshToken:      acc.RefreshToken,
		RefreshExpiry:     acc.RefreshExpiry,
		IDToken:           acc.IDToken,
		TokenType:         acc.TokenType,
		Email:             acc.Email,
		EmailVerified:     acc.EmailVerified,
		Name:              acc.Name,
		PreferredUsername: acc.PreferredUsername,
		Picture:           acc.Picture,
		Nickname:          acc.Nickname,
		Profile:           acc.Profile,
		UserID:            userID,
	}
	if err := s.q.AuthAccount.WithContext(ctx).Save(&authAcc); err != nil {
		return account.Account{}, err
	}
	return authAcc.MapToAccount(), nil
}

func (s *BasicGORMUserStore) GetAccountByID(ctx context.Context, id account.AccountID) (account.Account, error) {
	queryAccount := s.q.AuthAccount
	authAcc, err := queryAccount.WithContext(ctx).
		Where(
			queryAccount.Provider.Eq(id.Provider),
			queryAccount.ProviderID.Eq(id.ProviderID),
		).
		First()
	if err != nil {
		// TODO: check if it is not found err
		return account.Account{}, account.ErrAccountNotFound
	}
	return authAcc.MapToAccount(), nil
}

func (s *BasicGORMUserStore) GetAccountsByUserID(ctx context.Context, id account.UserID) ([]account.Account, error) {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}
	queryUser := s.q.AuthUser
	user, err := queryUser.WithContext(ctx).
		Preload(queryUser.Accounts).
		Where(queryUser.ID.Eq(userID)).
		First()
	if err != nil {
		// TODO: check if it is not found err
		return nil, account.ErrUserNotFound
	}
	return slices.Collect(
		its.Map(slices.Values(user.Accounts),
			func(acc models.AuthAccount) account.Account { return acc.MapToAccount() },
		),
	), nil
}

func (s *BasicGORMUserStore) UpdateAccount(ctx context.Context, id account.AccountID, acc account.Account) error {
	queryAccount := s.q.AuthAccount
	_, err := queryAccount.WithContext(ctx).
		Where(
			queryAccount.Provider.Eq(id.Provider),
			queryAccount.ProviderID.Eq(id.ProviderID),
		).
		Updates(models.AuthAccount{
			AccessToken:       acc.AccessToken,
			Expiry:            acc.Expiry,
			RefreshToken:      acc.RefreshToken,
			RefreshExpiry:     acc.RefreshExpiry,
			IDToken:           acc.IDToken,
			TokenType:         acc.TokenType,
			Email:             acc.Email,
			EmailVerified:     acc.EmailVerified,
			Name:              acc.Name,
			PreferredUsername: acc.PreferredUsername,
			Picture:           acc.Picture,
			Nickname:          acc.Nickname,
			Profile:           acc.Profile,
		})
	return err
}

func (s *BasicGORMUserStore) DeleteAccountByID(ctx context.Context, id account.AccountID) error {
	queryAccount := s.q.AuthAccount
	_, err := queryAccount.WithContext(ctx).
		Where(
			queryAccount.Provider.Eq(id.Provider),
			queryAccount.ProviderID.Eq(id.ProviderID),
		).
		Delete()
	return err
}

func (s *BasicGORMUserStore) GetAccountByUserIDAndProvider(ctx context.Context, id account.UserID, provider string) (account.Account, error) {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return account.Account{}, err
	}
	queryAccount := s.q.AuthAccount
	authAcc, err := queryAccount.WithContext(ctx).
		Where(
			queryAccount.UserID.Eq(userID),
			queryAccount.Provider.Eq(provider),
		).
		First()
	if err != nil {
		// TODO: check if it is not found err
		return account.Account{}, account.ErrAccountNotFound
	}
	return authAcc.MapToAccount(), nil
}
