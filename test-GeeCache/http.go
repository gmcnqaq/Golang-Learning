package test_GeeCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

// 首先创建一个结构体 HTTPPool，作为承载节点间 HTTP 通信的核心数据结构，此处实现服务端

// HTTPPool 第一个参数 self，用来记录自己的地址，包括主机名/IP 和端口
// 第二个参数 basePath，作为节点间通信地址的前缀，默认是 /geecache/，所以 http://example.com/_geecache/ 开头的请求就用于节点间的访问。
type HTTPPool struct {
	self     string
	basePath string
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basePath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupNmae := parts[0]
	key := parts[1]

	group := GetGroup(groupNmae)
	if group == nil {
		http.Error(w, "no such group"+groupNmae, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
