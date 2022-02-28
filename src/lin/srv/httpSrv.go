package main

import (
	"lin/log"
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

func NewHttpSrvMgr(ip string, port int) (*HttpSrvMgr, error) {
	srv := &HttpSrvMgr{
		ip:ip,
		port:port,
		mapCallBack:make(MAP_HTTP_CALLBACK),
	}
	srv.httpSrv = &http.Server{Addr: ip + ":" + strconv.Itoa(port), Handler: srv}
	go func() {
		err := srv.httpSrv.ListenAndServe()
		if err != nil {
			log.LogErr(err)
		}
	}()
	return srv, nil
}

func (self *HttpSrvMgr) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	if r == nil {
		log.LogDebug("request is nil")
		return
	}
	if r.URL == nil {
		log.LogDebug("request url is nil")
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.LogDebug("parse form err")
		return
	}
	log.LogDebug("request ", r.URL.Path)

	//find method from self.mapCallBack
	fn, ok := self.mapCallBack[r.URL.Path]
	if !ok {
		log.LogDebug("no func for", r.URL.Path)
		return
	}
	fn(w, r)
}

func (self *HttpSrvMgr) HttpSrvAddCallback(path string, fn FUNC_HTTP_CALLBACK) {
	self.mapCallBack[path] = fn
}
