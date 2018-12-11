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
				cs := handler.NewCookieSwitcherHandler([]string{"XSRF-TOKEN"})
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

		Context("Response handler", func() {
			It("Should update the replaced cookies if there is a Set-Cookie header", func() {
				//Arrange
				cs := handler.NewCookieSwitcherHandler([]string{"laravel_session", "XSRF-TOKEN"})
				r := &http.Response{
					Header: map[string][]string{"Set-Cookie": {
						"laravel_session=new_value2; expires=Mon, 10-Dec-2018 22:01:30 GMT; Max-Age=7200; path=/; httponly",
						"XSRF-TOKEN=new_value; expires=Mon, 10-Dec-2018 22:01:30 GMT; Max-Age=7200; path=/",
					}},
				}
				ctx := &goproxy.ProxyCtx{}

				//Act
				cs.ResponseHandler(r, ctx)

				//Assert
				replacedCookies := cs.ReplacedCookies()
				Expect(len(replacedCookies)).To(Equal(2))
				sessionCookie, ok := replacedCookies["laravel_session"]
				Expect(ok).To(Equal(true))
				Expect(sessionCookie).To(Equal("new_value2"))
				xsrfTokenCookie, ok := replacedCookies["XSRF-TOKEN"]
				Expect(ok).To(Equal(true))
				Expect(xsrfTokenCookie).To(Equal("new_value"))

			})
		})
	})
})
