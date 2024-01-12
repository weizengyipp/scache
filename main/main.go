package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/weizengyipp/scache/scache"
)

var FDB = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

var LoadCounts = make(map[string]int, len(FDB))

var f scache.Getter = scache.GetterFunc(func(key string) ([]byte, error) {
	log.Println("db searching:", key)
	if v, ok := FDB[key]; ok {
		LoadCounts[key] += 1
		return []byte(v), nil
	}
	return nil, fmt.Errorf("%s not exist", key)
})

func createGroup() *scache.Group {
	return scache.NewGroup("scores", 2<<10, f)
}

func startCacheServer(addr string, addrs []string, group *scache.Group) {
	peers := scache.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Println("scache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, group *scache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "scache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	sca := createGroup()
	if api {
		go startAPIServer(apiAddr, sca)
	}
	startCacheServer(addrMap[port], addrs, sca)
}
