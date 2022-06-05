package middleware

import (
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/boris-army/server/internal/adapters/http/render"
	"github.com/boris-army/server/internal/core/domain"
	"github.com/boris-army/server/internal/impl/session"
)

func TestAccess_Apply(t *testing.T) {
	type tc struct {
		name          string
		hasSomeTok    bool
		expDriverErr  error
		expEnforcerFn bool
		expResCode    int
		expResBody    []byte
	}

	tcs := []tc{
		{
			name:          "ok",
			hasSomeTok:    true,
			expDriverErr:  nil,
			expResCode:    200,
			expEnforcerFn: true,
			expResBody:    []byte(`ok`),
		},
		{
			name:         "no token",
			hasSomeTok:   false,
			expDriverErr: nil,
			expResCode:   fasthttp.StatusUnauthorized,
			expResBody:   []byte(`{"err":{"code":"` + render.CodeTokenRequired + `"}}`),
		},
		{
			name:         "invalid token",
			hasSomeTok:   true,
			expDriverErr: domain.ErrValue,
			expResCode:   fasthttp.StatusUnauthorized,
			expResBody:   []byte(`{"err":{"code":"` + render.CodeTokenInvalid + `"}}`),
		},
		{
			name:         "expired token",
			hasSomeTok:   true,
			expDriverErr: domain.ErrExpired,
			expResCode:   fasthttp.StatusUnauthorized,
			expResBody:   []byte(`{"err":{"code":"` + render.CodeTokenExpired + `"}}`),
		},
		{
			name:         "revoked token",
			hasSomeTok:   true,
			expDriverErr: domain.ErrSessionTerminated,
			expResCode:   fasthttp.StatusUnauthorized,
			expResBody:   []byte(`{"err":{"code":"` + render.CodeTokenRevoked + `"}}`),
		},
		{
			name:         "insufficient token",
			hasSomeTok:   true,
			expDriverErr: nil,
			expResCode:   fasthttp.StatusForbidden,
			expResBody:   []byte(`{"err":{"code":"` + render.CodeTokenInsufficient + `"}}`),
		},
		{
			name:         "internal",
			hasSomeTok:   true,
			expDriverErr: io.ErrShortWrite,
			expResCode:   fasthttp.StatusInternalServerError,
			expResBody:   []byte(`{"err":{"code":"` + render.CodeInternal + `"}}`),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sess := session.NewMockSessionDriver(ctrl)
			if tc.hasSomeTok {
				sess.EXPECT().
					DecodeHttpTokenTo(&domain.SessionHttpToken{}, []byte("tok")).
					Return(tc.expDriverErr)
			}

			mw := Access{Sessions: sess}
			handler := mw.Apply(func(req *fasthttp.RequestCtx, token *domain.SessionHttpToken) {
				_, _ = req.WriteString("ok")
			}, func(_ *domain.SessionHttpToken) bool {
				return tc.expEnforcerFn
			})

			req := fasthttp.RequestCtx{}
			if tc.hasSomeTok {
				req.Request.Header.Set(fasthttp.HeaderAuthorization, "Bearer tok")
			}

			handler(&req)
			assert.Equal(t, tc.expResCode, req.Response.StatusCode())
			assert.Equal(t, tc.expResBody, req.Response.Body())
		})
	}
}

func TestAccess_getTokenRaw(t *testing.T) {
	type tc struct {
		name string
		req  fasthttp.Request
	}
	headerReq := fasthttp.Request{}
	headerReq.Header.Set(fasthttp.HeaderAuthorization, "Bearer tok")

	queryReq := fasthttp.Request{}
	queryReq.SetRequestURI("https://example.com/private?access_token=tok")

	cookieReq := fasthttp.Request{}
	cookieReq.Header.SetCookie("access_token", "tok")

	tcs := []tc{
		{
			name: "header",
			req:  headerReq,
		},
		{
			name: "query",
			req:  queryReq,
		},
		{
			name: "cookie",
			req:  cookieReq,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := fasthttp.RequestCtx{Request: tc.req}
			tok := (&Access{}).getTokenRaw(&ctx)
			assert.Equal(t, tok, []byte("tok"))
		})
	}
}
