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

func main() {
	flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	scache.NewGroup("scores", 2<<10, f)
	addr := "localhost:8080"
	peers := scache.NewHTTPPool(addr)
	log.Println("scache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
