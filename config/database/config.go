package database

type Config interface {
	Dsn() string
}
