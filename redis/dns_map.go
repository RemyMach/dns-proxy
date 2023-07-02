package redis

func AddInDnsMap(entry string, address string) {
	redisClient := GetRedisClient()

	redisClient.Set("dns:" + entry, address, 0)
}

func DeleteInDnsMap(entry string) {
	redisClient := GetRedisClient()

	redisClient.Del("dns:" + entry)
}

func GetInDnsMap(entry string) {
	redisClient := GetRedisClient()

	redisClient.Get("dns:" + entry)
}