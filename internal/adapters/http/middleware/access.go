package middleware

import (
	"github.com/valyala/fasthttp"

	"github.com/boris-army/server/internal/adapters/http/render"
	"github.com/boris-army/server/internal/core/domain"
	"github.com/boris-army/server/internal/core/ports"
)

type Access struct {
	Sessions ports.SessionDriver
}

type HandlerWithAccess = func(*fasthttp.RequestCtx, *domain.SessionHttpToken)

type AccessEnforcerFn = func(*domain.SessionHttpToken) bool

func (m *Access) Apply(next HandlerWithAccess, enforcerFn AccessEnforcerFn) fasthttp.RequestHandler {
	return func(req *fasthttp.RequestCtx) {
		tokRaw := m.getTokenRaw(req)
		if len(tokRaw) == 0 {
			render.ErrAccessTokenRequired(req)
			return
		}

		tok := domain.AcquireSessionHttpToken()
		defer domain.ReleaseSessionHttpToken(tok)

		if err := m.Sessions.DecodeHttpTokenTo(tok, tokRaw); err != nil {
			switch err {
			case domain.ErrValue:
				render.ErrAccessTokenInvalid(req)
				return

			case domain.ErrExpired:
				render.ErrAccessTokenExpired(req)
				return

			case domain.ErrSessionTerminated:
				render.ErrAccessTokenRevoked(req)
				return

			default:
				render.ErrInternal(req, "")
				return
			}
		}

		if !enforcerFn(tok) {
			render.ErrAccessTokenInsufficient(req)
			return
		}

		next(req, tok)
	}
}

func (m *Access) getTokenRaw(req *fasthttp.RequestCtx) []byte {
	authHdr := req.Request.Header.Peek(fasthttp.HeaderAuthorization)
	const nBearerPrefix = len("Bearer ")
	if len(authHdr) > nBearerPrefix {
		return authHdr[nBearerPrefix:]
	}

	tokValue := req.QueryArgs().Peek("access_token")
	if len(tokValue) > 0 {
		return tokValue
	}

	return req.Request.Header.Cookie("access_token")
}
