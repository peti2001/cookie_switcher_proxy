package handler

import (
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

type cookieSwitcher struct {
	cookies map[string]string
	replace []string
}

func NewCookieSwitcherHandler(needsToBeReplaced []string) cookieSwitcher {
	cs := cookieSwitcher{
		make(map[string]string),
		needsToBeReplaced,
	}
	return cs
}

func (cs cookieSwitcher) isInTheList(cookieName string) bool {
	for _, c := range cs.replace {
		if c == cookieName {
			return true
		}
	}

	return false
}

func (cs cookieSwitcher) ResponseHandler(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	return resp
}

func (cs cookieSwitcher) RequestHandler(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	newCookies := make([]string, len(r.Cookies()))
	for i, cookie := range r.Cookies() {
		if cs.isInTheList(cookie.Name) {
			newCookies[i] = cookie.Name + "=new_value"
		}
	}

	r.Header.Set("Cookie", strings.Join(newCookies, "; "))

	return r, nil
}
