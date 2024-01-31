package mail

type Config interface {
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetTemplateFilePath() string
}
