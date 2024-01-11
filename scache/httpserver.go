package scache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/scache"

type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		http.Error(w, "no such cache server: "+r.URL.Path, http.StatusNotFound)
		return
	}
	log.Printf("URL: %s", r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.Split(r.URL.Path, "/")
	log.Println("url split:", parts)
	if len(parts) != 4 {
		log.Println("len(parts) != 4", len(parts))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[2]
	key := parts[3]
	log.Printf("groupname %s key %s", groupName, key)
	group := GetGroup(groupName)
	if group == nil {
		log.Println("no such group: " + groupName)
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
