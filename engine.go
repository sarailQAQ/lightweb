package lightweb

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler func(*Context)

type HandlerChain []Handler

const defaultMultipartMemory = 32 << 20

func (hc HandlerChain) cat(nxt HandlerChain) {
	hc = append(hc, nxt...)
}


type Engine struct {
	RouterGroup

	MaxMultipartMemory int64

	router MethodHash
}

func New() *Engine{
	return &Engine{
		RouterGroup:        RouterGroup{
			basePath: "",
			root: false,
		},
		MaxMultipartMemory: defaultMultipartMemory,
		router:             make(MethodHash),
	}
}

func Default() *Engine {
	engine := New()
	engine.engine = engine
	engine.RouterGroup.root = true
	return engine
}

func (engine *Engine) Run(port int) {
	portS := strconv.FormatInt(int64(port), 10)
	http.Handle("/", engine)
	fmt.Println("Listen and server http on :",portS)
	if err := http.ListenAndServe(":"+portS, engine); err != nil {
		log.Println(err)
	}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	httpMethod := req.Method
	uri := req.RequestURI
	uris := strings.Split(uri, "?")
	if len(uris) < 1 {
		return
	}

	hs,params,ok := engine.router.getRouter(httpMethod,uris[0])
	if !ok {
		Handler404(w, req)
		return
	}

	c := NewContext(w, req,hs,engine)
	c.Params = params
	if len(hs) < 1 {
		return
	}
	c.handlers[0](&c)
}

func Handler404(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("404 not found"))
}

