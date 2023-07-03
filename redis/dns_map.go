package redis

import (
	"errors"
	"log"
)

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
	if element.Err() != nil {
		return "", errors.New("error its impossible")
	}

	return element.Val(), nil
}

func GetAllDnsMap() map[string]string {
	redisClient := GetRedisClient()

	var cursor uint64
	keyValueMap := make(map[string]string)
	for {
		var keys []string
		var err error
		keys, cursor, err = redisClient.Scan(cursor, "dns:*", 10).Result()
		if err != nil {
			log.Println("Error when try to get all keys")
		}

		for _, key := range keys {
			val, err := redisClient.Get(key).Result()
			if err != nil {
				log.Printf("Error when getting value for key %s: %v", key, err)
			} else {
				keyValueMap[key] = val
			}
		}

		if cursor == 0 {
			break
		}
	}

	return keyValueMap
}