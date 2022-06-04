package render

import "github.com/valyala/fasthttp"

const (
	CodeValue             = "VALUE"
	CodeUserExists        = "USER_EXISTS"
	CodeInternal          = "INTERNAL"
	CodeTokenRequired     = "ACCESS_TOKEN_REQUIRED"
	CodeTokenInvalid      = "ACCESS_TOKEN_INVALID"
	CodeTokenExpired      = "ACCESS_TOKEN_EXPIRED"
	CodeTokenRevoked      = "ACCESS_TOKEN_REVOKED"
	CodeTokenInsufficient = "FORBIDDEN"
)

const (
	authRealm   = "boris.army"
	authPrelude = `Bearer authRealm="` + authRealm + `"`
)

func ErrBadReq(w *fasthttp.RequestCtx, code, mes string) {
	w.SetStatusCode(fasthttp.StatusBadRequest)
	Err(w, code, mes)
}

func ErrConflict(w *fasthttp.RequestCtx, code, mes string) {
	w.SetStatusCode(fasthttp.StatusConflict)
	Err(w, code, mes)
}

func ErrInternal(w *fasthttp.RequestCtx, mes string) {
	w.SetStatusCode(fasthttp.StatusInternalServerError)
	Err(w, CodeInternal, mes)
}

func Err(w *fasthttp.RequestCtx, code, mes string) {
	w.SetContentType("application/json")
	_, _ = w.WriteString(`{"err":{"code":"`)
	_, _ = w.WriteString(code)
	if len(mes) > 0 {
		_, _ = w.WriteString(`","mes":"`)
		_, _ = w.WriteString(mes)
	}
	_, _ = w.WriteString(`"}}`)
}

func ErrAccessTokenRequired(w *fasthttp.RequestCtx) {
	setNoTokenHeaders(w)
	Err(w, CodeTokenRequired, "")
}

func ErrAccessTokenInvalid(w *fasthttp.RequestCtx) {
	setInvalidTokenHeaders(w)
	Err(w, CodeTokenInvalid, "")
}

func ErrAccessTokenExpired(w *fasthttp.RequestCtx) {
	setInvalidTokenHeaders(w)
	const desc = `error_description="The access token expired"`
	w.Response.Header.Add(fasthttp.HeaderWWWAuthenticate, desc)
	Err(w, CodeTokenExpired, "")
}

func ErrAccessTokenRevoked(w *fasthttp.RequestCtx) {
	setInvalidTokenHeaders(w)
	const desc = `error_description="The access token had been revoked"`
	w.Response.Header.Add(fasthttp.HeaderWWWAuthenticate, desc)
	Err(w, CodeTokenRevoked, "")
}

func ErrAccessTokenInsufficient(w *fasthttp.RequestCtx) {
	w.SetStatusCode(fasthttp.StatusForbidden)
	w.Response.Header.Add(fasthttp.HeaderWWWAuthenticate, authPrelude)
	w.Response.Header.Add(fasthttp.HeaderWWWAuthenticate, `error="insufficient_scope"`)
	Err(w, CodeTokenInsufficient, "")
}

func setNoTokenHeaders(w *fasthttp.RequestCtx) {
	w.SetStatusCode(fasthttp.StatusUnauthorized)
	w.Response.Header.Set(fasthttp.HeaderWWWAuthenticate, authPrelude)
}

func setInvalidTokenHeaders(w *fasthttp.RequestCtx) {
	setNoTokenHeaders(w)
	w.Response.Header.Add(fasthttp.HeaderWWWAuthenticate, `error="invalid_token"`)
}
