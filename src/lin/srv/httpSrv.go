package main

import (
	"lin/lin_common"
	"net/http"
	"strconv"
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
	go func() {
		err := srv.httpSrv.ListenAndServe()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()
	return srv, nil
}

func (pthis *HttpSrvMgr) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	if r == nil {
		lin_common.LogDebug("request is nil")
		return
	}
	if r.URL == nil {
		lin_common.LogDebug("request url is nil")
		return
	}
	err := r.ParseForm()
	if err != nil {
		lin_common.LogDebug("parse form err")
		return
	}
	lin_common.LogDebug("request ", r.URL.Path)

	//find method from self.mapCallBack
	fn, ok := pthis.mapCallBack[r.URL.Path]
	if !ok {
		lin_common.LogDebug("no func for", r.URL.Path)
		return
	}
	func(){
		defer func() {
			err := recover()
			if err != nil {
				lin_common.LogErr(r.URL.Path, err)
			}
		}()
		fn(w, r)
	}()
}

func (pthis *HttpSrvMgr) HttpSrvAddCallback(path string, fn FUNC_HTTP_CALLBACK) {
	pthis.mapCallBack[path] = fn
}
