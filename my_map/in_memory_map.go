package mymap

import (
	"encoding/json"
	"errors"
	"sync"
)

var dnsMap = make(map[string]string, 10000)
var dnsMapMutex = &sync.Mutex{}

func AddInDnsMap(entry string, address string) {

	dnsMapMutex.Lock()
	dnsMap[entry] = address
	dnsMapMutex.Unlock()
}

func DeleteInDnsMap(entry string) {
	dnsMapMutex.Lock()
	delete(dnsMap, entry)
	dnsMapMutex.Unlock()
}

func GetInDnsMap(entry string) (string, error) {
	dnsMapMutex.Lock()
	ip, ok := dnsMap[entry]
	dnsMapMutex.Unlock()

	if ok {
		return ip, nil
	} else {
		return "", errors.New("error to get entry")
	}
}

func GetAllDnsMap() []byte {
	dnsMapMutex.Lock()
	defer dnsMapMutex.Unlock()

	b, err := json.Marshal(dnsMap)
	if err != nil {
		return nil
	}

	return b
}

func InitMap() {
	dnsMap["pomme.worker.stuga-cloud.tech."] = "65.109.94.8"
}
