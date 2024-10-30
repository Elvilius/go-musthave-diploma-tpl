package handler

import (
	"context"
	"net/http"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/middleware"
)

func SetAuth(request *http.Request) *http.Request {
	request = request.WithContext(context.WithValue(request.Context(), middleware.UserIDKey, 1))
	request.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjJ9.6Oz5eGuwTSWswdvgsxbhvDIBkd9YKzxJSyd9mg4auBM")
	return request
}
