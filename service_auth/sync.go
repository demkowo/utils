package serviceauth

import "errors"

func HandleSyncAPIKey(payload SyncAPIKeyPayload, redis RedisClient) error {
	if payload.Service == "" || payload.Key == "" {
		return errors.New("invalid sync payload")
	}
	return redis.SetServiceKey(payload.Service, payload.Key)
}
