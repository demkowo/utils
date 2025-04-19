package serviceauth

type RedisClient interface {
	GetServiceKey(service string) (string, error)
	SetServiceKey(service string, key string) error
	HasServiceKey(service string) (bool, error)
}
