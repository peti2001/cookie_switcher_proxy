package handler

import (
	"net/http"
	"strings"
	"log"

	"github.com/elazarl/goproxy"
)

type cookieSwitcher struct {
	replacedCookies   map[string]string
	needsToBeReplaced []string
}

func NewCookieSwitcherHandler(needsToBeReplaced []string) *cookieSwitcher {
	cs := &cookieSwitcher{
		make(map[string]string),
		needsToBeReplaced,
	}
	return cs
}

func (cs *cookieSwitcher) SetReplacedCookies(cookies map[string]string) {
	cs.replacedCookies = cookies
}

func (cs *cookieSwitcher) ReplacedCookies() map[string]string {
	return cs.replacedCookies
}

func (cs *cookieSwitcher) isInTheList(cookieName string) bool {
	for _, c := range cs.needsToBeReplaced {
		if c == cookieName {
			return true
		}
	}

	return false
}

func (cs *cookieSwitcher) ResponseHandler(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	newCookies := resp.Cookies()
	for _, cookie := range newCookies {
		cs.replacedCookies[cookie.Name] = cookie.Value
		log.Println("Set Cookie:", cookie.Name, ":", cookie.Value)
	}

	return resp
}

func (cs *cookieSwitcher) RequestHandler(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	newCookies := make([]string, len(r.Cookies()))
	for i, cookie := range r.Cookies() {
		if cs.isInTheList(cookie.Name) {
			if newValue, ok := cs.replacedCookies[cookie.Name]; ok {
				newCookies[i] = cookie.Name + "=" + newValue
			} else {
				newCookies[i] = cookie.Name + "=" + cookie.Value
			}
		} else {
			newCookies[i] = cookie.Name + "=" + cookie.Value
		}
	}

	r.Header.Set("Cookie", strings.Join(newCookies, "; "))

	return r, nil
}
