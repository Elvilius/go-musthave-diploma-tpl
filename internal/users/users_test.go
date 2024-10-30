package users

import (
	"context"
	"testing"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	mocks_users "github.com/Elvilius/go-musthave-diploma-tpl.git/internal/users/mocks"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CreateNewUser(t *testing.T) {
	type args struct {
		registerUser models.UserLogin
	}
	tests := []struct {
		name string
		want models.JWTToken
		args args
	}{
		{name: "Create new user", want: models.JWTToken("secret"), args: args{registerUser: models.UserLogin{Login: "test", Password: "test"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks_users.NewMockStorer(ctrl)

			token := jwt.NewMockToken()

			s := &Service{
				store: m,
				token: token,
				cfg:   nil,
			}

			m.EXPECT().CreateUser(gomock.Any(), tt.args.registerUser.Login, gomock.Any()).Return(1, nil)

			got, err := s.CreateNewUser(context.TODO(), tt.args.registerUser)

			assert.Equal(t, got, tt.want)
			require.NoError(t, err)
		})
	}
}
