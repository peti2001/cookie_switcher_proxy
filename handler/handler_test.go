package handler_test

import (
	"net/http"

	"github.com/elazarl/goproxy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/peti2001/csrf_changer/handler"
)

var _ = Describe("Handler", func() {
	Describe("CookieSwitcherHandler", func() {
		Context("Request Handler", func() {
			It("Should replace the session and xsrf-token cookies", func() {
				//Arrange
				cs := handler.NewCookieSwitcherHandler([]string{"laravel_session", "XSRF-TOKEN"})
				cookies := map[string]string{"laravel_session": "new_value", "XSRF-TOKEN": "new_value2"}
				cs.SetReplacedCookies(cookies)
				r := &http.Request{
					Header: map[string][]string{"Cookie": {"LS_CSRF_TOKEN=a; XSRF-TOKEN=b; laravel_session=c"}},
				}
				ctx := &goproxy.ProxyCtx{}

				//Act
				req, _ := cs.RequestHandler(r, ctx)

				//Assert
				sessionCookie, err := req.Cookie("laravel_session")
				Expect(err).To(BeNil())
				xsrfTokenCookie, err := req.Cookie("XSRF-TOKEN")
				Expect(err).To(BeNil())

				Expect(sessionCookie.Value).To(Equal("new_value"))
				Expect(xsrfTokenCookie.Value).To(Equal("new_value2"))
			})
			It("Should not remove a cookie even if no need to change", func() {
				//Arrange
				cs := handler.NewCookieSwitcherHandler([]string{})
				cs.SetReplacedCookies(map[string]string{})
				r := &http.Request{
					Header: map[string][]string{"Cookie": {"LS_CSRF_TOKEN=a; XSRF-TOKEN=b; laravel_session=c"}},
				}
				ctx := &goproxy.ProxyCtx{}

				//Act
				req, _ := cs.RequestHandler(r, ctx)

				//Assert
				csrfTokenCookie, err := req.Cookie("LS_CSRF_TOKEN")
				Expect(err).To(BeNil())
				sessionCookie, err := req.Cookie("laravel_session")
				Expect(err).To(BeNil())
				xsrfTokenCookie, err := req.Cookie("XSRF-TOKEN")
				Expect(err).To(BeNil())

				Expect(csrfTokenCookie.Value).To(Equal("a"))
				Expect(xsrfTokenCookie.Value).To(Equal("b"))
				Expect(sessionCookie.Value).To(Equal("c"))
			})
		})
	})
})
