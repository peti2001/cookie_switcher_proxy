package handler_test

import (
	"net/http"

	"github.com/elazarl/goproxy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/peti2001/csrf_changer/handler"
)

var _ = Describe("Handler", func() {
	Describe("CookieSwitcher", func() {
		Context("Request Handler", func() {
			It("Should replace the session and xsrf-token cookies", func() {
				//Arrange
				cs := handler.NewCookieSwitcherHandler([]string{"laravel_session", "XSRF-TOKEN"})
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
				Expect(xsrfTokenCookie.Value).To(Equal("new_value"))
			})
		})
	})
})
