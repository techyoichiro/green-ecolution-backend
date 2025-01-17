package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRestrospectToken(t *testing.T) {
	t.Run("should return result when success", func(t *testing.T) {
		// given
		token := "token"
		identityConfig := &config.IdentityAuthConfig{}
		expected := &entities.IntroSpectTokenResult{
			Active:   utils.P(true),
			Exp:      utils.P(123),
			AuthTime: utils.P(123),
			Type:     utils.P("token"),
		}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		authRepo.EXPECT().RetrospectToken(context.Background(), token).Return(expected, nil)
		resp, err := svc.RetrospectToken(context.Background(), token)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("should return error when failed to retrospect token", func(t *testing.T) {
		// given
		token := "token"
		identityConfig := &config.IdentityAuthConfig{}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		authRepo.EXPECT().RetrospectToken(context.Background(), token).Return(nil, errors.New("failed to retrospect token"))
		_, err := svc.RetrospectToken(context.Background(), token)

		// then
		assert.Error(t, err)
	})
}
