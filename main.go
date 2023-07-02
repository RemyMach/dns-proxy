package main

import (
	"dns-proxy/config"
	"dns-proxy/redis"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	mymap "dns-proxy/my_map"

	"github.com/miekg/dns"
)

type entry struct {
	Domain string `json:"domain"`
	IP     string `json:"ip"`
}

func main() {
	config.Init()
	dns.HandleFunc(".", handleRequest)
	go func() {
		server := &dns.Server{Addr: ":53", Net: "udp"}
		log.Fatal(server.ListenAndServe())
	}()

	redis.InitRedis()
	http.HandleFunc("/health", handleHealthCheck)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/list", handleList)     // Ajout de la route pour lister les entrées
	http.HandleFunc("/delete", handleDelete) // Ajout de la route pour supprimer les entrées
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkToken(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token == os.Getenv("API_KEY")
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	if !checkToken(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Nous supposons que le corps de la requête contient un JSON
	// avec les champs "domain" et "ip".
	var data struct {
		Domain string `json:"domain"`
		IP     string `json:"ip"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if data.Domain == "" || data.IP == "" {
		http.Error(w, "Domain or IP not provided", http.StatusBadRequest)
		return
	}
	if os.Getenv("CONFIG_STORAGE") != "REDIS" {
		mymap.AddInDnsMap(data.Domain+".", data.IP)
	} else {
		redis.AddInDnsMap(data.Domain+".", data.IP)
	}
}

func handleList(w http.ResponseWriter, r *http.Request) {
	if !checkToken(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// dnsMapMutex.Lock()
	// defer dnsMapMutex.Unlock()

	// b, err := json.Marshal(dnsMap)
	// if err != nil {
	// 	http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
	// 	return
	// }
	var list []byte
	if os.Getenv("CONFIG_STORAGE") != "REDIS" {
		list = mymap.GetAllDnsMap()
	} else {
		list = redis.GetAllDnsMap()
	}
	w.Write(list)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if !checkToken(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Nous supposons que le corps de la requête contient un JSON avec le champ "domain".
	var data struct {
		Domain string `json:"domain"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if data.Domain == "" {
		http.Error(w, "Domain not provided", http.StatusBadRequest)
		return
	}

	// dnsMapMutex.Lock()
	// delete(dnsMap, data.Domain+".")
	// dnsMapMutex.Unlock()
	if os.Getenv("CONFIG_STORAGE") != "REDIS" {
		mymap.DeleteInDnsMap(data.Domain + ".")
	} else {
		redis.DeleteInDnsMap(data.Domain + ".")
	}
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	for _, question := range r.Question {
		switch question.Qtype {
		case dns.TypeA:
			// dnsMapMutex.Lock()
			// ip, ok := dnsMap[question.Name]
			// dnsMapMutex.Unlock()
			var ip string
			var err error
			if os.Getenv("CONFIG_STORAGE") != "REDIS" {
				ip, err = mymap.GetInDnsMap(question.Name)
			} else {
				ip, err = redis.GetInDnsMap(question.Name)
			}

			if err == nil {
				rr, _ := dns.NewRR(question.Name + " IN A " + ip)
				m.Answer = append(m.Answer, rr)
			} else {
				c := new(dns.Client)
				msg, _, _ := c.Exchange(r, net.JoinHostPort("8.8.8.8", "53"))
				m.Answer = append(m.Answer, msg.Answer...)
			}
		}
	}

	w.WriteMsg(m)
}
