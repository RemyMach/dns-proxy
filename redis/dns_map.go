package redis

import "errors"

func AddInDnsMap(entry string, address string) {
	redisClient := GetRedisClient()

	redisClient.Set("dns:" + entry, address, 0)
}

func DeleteInDnsMap(entry string) {
	redisClient := GetRedisClient()

	redisClient.Del("dns:" + entry)
}

func GetInDnsMap(entry string) (string, error) {
	redisClient := GetRedisClient()

	element := redisClient.Get("dns:" + entry)
	if element.Err() {
		return "", errors.New("error its impossible")
	}

	return element.Val(), nil
}

func GetAllDnsMap() []byte {
	redisClient := GetRedisClient()

	elements := redisClient.Get("dns:*")

	if elements.Err() {
		return []byte{}
	}

	return elements.Val()
}