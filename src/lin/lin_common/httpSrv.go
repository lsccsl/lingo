package lin_common

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type FUNC_HTTP_CALLBACK func(http.ResponseWriter, *http.Request)

type httpCallBack struct {
	help string
	fn FUNC_HTTP_CALLBACK
}
type MAP_HTTP_CALLBACK map[string]*httpCallBack

type HttpSrvMgr struct {
	ip          string
	port        int

	mapCallBack MAP_HTTP_CALLBACK
	mapCallBackLock sync.Mutex

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

	srv.addCallback("/help",
		func(writer http.ResponseWriter, request *http.Request){
			srv.mapCallBackLock.Lock()
			defer srv.mapCallBackLock.Unlock()

			var str string
			for key, val := range srv.mapCallBack {
				str += " cmd:" + key + " " + val.help
			}
			writer.Write([]byte(str))
		}, " help")

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
	fn := pthis.getCallback(r.URL.Path)
	if fn == nil {
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

func (pthis *HttpSrvMgr) HttpSrvAddCallback(path string, fn FUNC_HTTP_CALLBACK, help string) {
	pthis.addCallback(path, fn, help)
}

func (pthis *HttpSrvMgr) addCallback(path string, fn FUNC_HTTP_CALLBACK, help string) {
	pthis.mapCallBackLock.Lock()
	defer pthis.mapCallBackLock.Unlock()

	pthis.mapCallBack[path] = &httpCallBack{help, fn}
}

func (pthis *HttpSrvMgr) getCallback(path string) (fn FUNC_HTTP_CALLBACK) {
	pthis.mapCallBackLock.Lock()
	defer pthis.mapCallBackLock.Unlock()

	cb, _ := pthis.mapCallBack[path]
	if cb == nil {
		return nil
	}
	fn = cb.fn
	return
}
