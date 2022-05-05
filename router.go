package tgin

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	paths := strings.Split(pattern, "/")
	s := make([]string, 0)
	for _, path := range paths {
		if path != "" {
			s = append(s, path)
			if path[0] == '*' {
				break
			}
		}
	}
	return s
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	paths := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, paths, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	searchPaths := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchPaths, 0)
	if n != nil {
		paths := parsePattern(n.pattern)
		for index, path := range paths {
			if path[0] == ':' {
				params[path[1:]] = searchPaths[index]
			}
			if path[0] == '*' && len(path) > 1 {
				params[path[1:]] = strings.Join(searchPaths[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(ctx *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
