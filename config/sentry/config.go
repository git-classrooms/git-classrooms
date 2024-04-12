package sentry

type Config interface {
	GetDSN() string
	GetEnvironment() string
	IsEnabled() bool
}
