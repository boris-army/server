package render

import "github.com/valyala/fasthttp"

const (
	CodeValue      = "VALUE"
	CodeUserExists = "USER_EXISTS"
	CodeInternal   = "INTERNAL"
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
