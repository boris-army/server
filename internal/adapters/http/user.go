package http

import (
	"sync"

	"github.com/valyala/fasthttp"

	"github.com/boris-army/server/internal/adapters/http/render"
	"github.com/boris-army/server/internal/core/domain"
	"github.com/boris-army/server/internal/core/ports"
)

//go:generate easyjson $GOFILE

//easyjson:json
type UserPostCtx struct {
	Email      string                  `json:"email,nocopy"`
	Surname    string                  `json:"surname,nocopy"`
	GivenNames string                  `json:"given_names,nocopy"`
	Password   string                  `json:"password,nocopy"`
	CreateUser ports.CommandUserCreate `json:"-"`
}

func (r *UserPostCtx) Reset() {
	r.Email = ""
	r.Surname = ""
	r.GivenNames = ""
	r.Password = ""
	r.CreateUser.Reset()
}

var userPostCtxPool = sync.Pool{
	New: func() any {
		return &UserPostCtx{}
	},
}

func (a *Adapter) UserPost(req *fasthttp.RequestCtx) {
	ctx := userPostCtxPool.Get().(*UserPostCtx)
	defer func() {
		ctx.Reset()
		userPostCtxPool.Put(ctx)
	}()

	if err := ctx.UnmarshalJSON(req.PostBody()); err != nil {
		return
	}

	cmd := &ctx.CreateUser
	cmd.Email = ctx.Email
	cmd.Surname = ctx.Surname
	cmd.GivenNames = ctx.GivenNames
	cmd.Password = ctx.Password
	if err := a.Users.Create(cmd); err != nil {
		switch err {
		case domain.ErrValue:
			render.ErrBadReq(req, render.CodeValue, "")
			return

		case domain.ErrExists:
			render.ErrConflict(req, render.CodeUserExists, "")
			return

		default:
			render.ErrInternal(req, "")
			return
		}
	}

	req.SetContentType("application/json")
	_, _ = req.WriteString(`{"res":"24h email confirmation"}`)
}
