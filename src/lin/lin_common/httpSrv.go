package lin_common

import (
	"net/http"
	"strconv"
	"time"
)
type FUNC_HTTP_CALLBACK func(http.ResponseWriter, *http.Request)
type MAP_HTTP_CALLBACK map[string]FUNC_HTTP_CALLBACK
type HttpSrvMgr struct {
	ip          string
	port        int
	mapCallBack MAP_HTTP_CALLBACK
	httpSrv     *http.Server
}

func StartHttpSrvMgr(ip string, port int) (*HttpSrvMgr, error) {
	srv := &HttpSrvMgr{
		ip:ip,
		port:port,
		mapCallBack:make(MAP_HTTP_CALLBACK),
	}
	srv.httpSrv = &http.Server{Addr: ip + ":" + strconv.Itoa(port), Handler: srv}
	srv.httpSrv.IdleTimeout = time.Second * 10
	go func() {
		err := srv.httpSrv.ListenAndServe()
		if err != nil {
			LogErr(err)
		}
	}()
	return srv, nil
}

func (pthis *HttpSrvMgr) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	if r == nil {
		LogDebug("request is nil")
		return
	}
	if r.URL == nil {
		LogDebug("request url is nil")
		return
	}
	err := r.ParseForm()
	if err != nil {
		LogDebug("parse form err")
		return
	}
	//LogDebug("request ", r.URL.Path, r.PostForm)

	//find method from self.mapCallBack
	fn, ok := pthis.mapCallBack[r.URL.Path]
	if !ok {
		//LogDebug("no func for", r.URL.Path)
		return
	}
	func(){
		defer func() {
			err := recover()
			if err != nil {
				LogErr(r.URL.Path, err)
			}
		}()
		fn(w, r)
	}()
}

func (pthis *HttpSrvMgr) HttpSrvAddCallback(path string, fn FUNC_HTTP_CALLBACK) {
	pthis.mapCallBack[path] = fn
}
