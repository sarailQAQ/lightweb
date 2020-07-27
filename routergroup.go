package lightweb

import (
	"net/http"
	"regexp"
)

type RouterGroup struct {
	basePath string
	handlers HandlerChain
	engine *Engine
	root bool
}

type IRoutes interface {
	Use(...Handler) IRoutes

	Handle(string, string, ...Handler) IRoutes

	GET(string, ...Handler) IRoutes
	POST(string, ...Handler) IRoutes
}

func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...Handler) IRoutes {
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		panic("http method " + httpMethod + " is not valid")
	}
	return group.handle(httpMethod, relativePath, handlers)
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlerChain) IRoutes {
	absolutePath := joinPath(group.basePath,relativePath)
	handlers = group.combineHandlers(handlers)
	ok := group.engine.router.addRouter(httpMethod,absolutePath,handlers)
	if !ok {
		panic("same route")
	}
	return group.returnObj()
}

func (group *RouterGroup) Use(middleware ...Handler) IRoutes {
	group.handlers = append(group.handlers, middleware...)
	return group.returnObj()
}

func (group *RouterGroup) GET(relativePath string,handlers ...Handler) IRoutes {
	return group.handle(http.MethodGet, relativePath, handlers)
}

func (group *RouterGroup) POST(relativePath string,handlers ...Handler) IRoutes {
	return group.handle(http.MethodPost, relativePath, handlers)
}

func (group *RouterGroup) Group(relativePath string,handlers ...Handler) *RouterGroup {
	return &RouterGroup{
		basePath: joinPath(group.basePath,relativePath),
		handlers: group.combineHandlers(handlers),
		engine:   group.engine,
		root:	  false,
	}
}

func (group *RouterGroup) combineHandlers(handlers HandlerChain) HandlerChain {
	finalSize := len(group.handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlerChain, finalSize)
	copy(mergedHandlers, group.handlers)
	copy(mergedHandlers[len(group.handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}
