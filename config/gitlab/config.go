package gitlab

type Config interface {
	GetURL() string
}
