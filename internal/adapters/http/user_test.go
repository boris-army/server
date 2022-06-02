package http

import (
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/boris-army/server/internal/core/domain"
	"github.com/boris-army/server/internal/core/ports"
	"github.com/boris-army/server/internal/impl/user"
)

func TestUserPost_Response(t *testing.T) {
	type tc struct {
		name         string
		driverErr    error
		expRes       string
		expResStatus int
	}
	tcs := []tc{
		{"user exists", domain.ErrExists, `{"err":{"code":"USER_EXISTS"}}`, fasthttp.StatusConflict},
		{"value error", domain.ErrValue, `{"err":{"code":"VALUE"}}`, fasthttp.StatusBadRequest},
		{"internal error", io.ErrShortWrite /* other */, `{"err":{"code":"INTERNAL"}}`, fasthttp.StatusInternalServerError},
		{"ok", nil, `{"res":"24h email confirmation"}`, fasthttp.StatusOK},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userDriver := user.NewMockDriverUser(ctrl)
			a := Adapter{Users: userDriver}

			const reqBody = `
				{	
					"email": "pgarin@old.me",
					"surname": "Garin",
					"given_names": "Pyotr",
					"password": "qwerty123"
				}
			`
			userDriver.EXPECT().Create(&ports.CommandUserCreate{
				Email:      "pgarin@old.me",
				Surname:    "Garin",
				GivenNames: "Pyotr",
				Password:   "qwerty123",
			}).Return(tc.driverErr)

			req := &fasthttp.RequestCtx{}
			req.Request.SetBody([]byte(reqBody))
			a.UserPost(req)

			resBody := req.Response.Body()
			assert.Equal(t, req.Response.StatusCode(), tc.expResStatus)
			assert.Equal(t, string(resBody), tc.expRes)
			assert.Equal(t, string(req.Response.Header.ContentType()), "application/json")
		})
	}
}
