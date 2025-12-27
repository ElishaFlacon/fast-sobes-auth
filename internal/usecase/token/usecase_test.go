package tokens

import (
	"context"
	"testing"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/jwtmanager"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/testutil"
)

func TestVerifyToken(t *testing.T) {
	ctx := context.Background()
	log := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()

	_ = users.Create(ctx, &domain.User{Email: "user@example.com", PasswordHash: testutil.PasswordHash("pass1234")})
	user, _ := users.GetByEmail(ctx, "user@example.com")

	jwt := jwtmanager.New("secret")
	now := time.Now()
	exp := now.Add(time.Hour)
	tokenString, err := jwt.Sign(jwtmanager.Claims{
		UserID:          user.ID,
		Email:           user.Email,
		PermissionLevel: user.PermissionLevel,
		ExpiresAt:       exp.Unix(),
		IssuedAt:        now.Unix(),
	})
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	_ = tokens.Create(ctx, &domain.AccessToken{
		Token:     tokenString,
		UserID:    user.ID,
		ExpiresAt: exp,
		CreatedAt: now,
	})

	uc := NewUsecase(log, tokens, users, "secret")
	uc.now = func() time.Time { return now }

	info, err := uc.VerifyToken(ctx, tokenString)
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if !info.Valid || info.UserID != user.ID {
		t.Fatalf("unexpected token info: %+v", info)
	}
}

func TestVerifyTokenFailures(t *testing.T) {
	ctx := context.Background()
	log := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()

	_ = users.Create(ctx, &domain.User{Email: "user@example.com", PasswordHash: testutil.PasswordHash("pass1234")})
	user, _ := users.GetByEmail(ctx, "user@example.com")

	jwt := jwtmanager.New("secret")
	now := time.Now()
	exp := now.Add(-time.Hour) // expired
	tokenString, _ := jwt.Sign(jwtmanager.Claims{
		UserID:          user.ID,
		Email:           user.Email,
		PermissionLevel: user.PermissionLevel,
		ExpiresAt:       exp.Unix(),
		IssuedAt:        now.Unix(),
	})

	_ = tokens.Create(ctx, &domain.AccessToken{
		Token:     tokenString,
		UserID:    user.ID,
		ExpiresAt: exp,
		CreatedAt: now,
	})

	uc := NewUsecase(log, tokens, users, "secret")
	uc.now = func() time.Time { return now }

	if _, err := uc.VerifyToken(ctx, tokenString); err == nil {
		t.Fatal("expected expired token error")
	}

	_ = tokens.Create(ctx, &domain.AccessToken{
		Token:     "revoked",
		UserID:    user.ID,
		ExpiresAt: now.Add(time.Hour),
		Revoked:   true,
		CreatedAt: now,
	})
	if _, err := uc.VerifyToken(ctx, "revoked"); err == nil {
		t.Fatal("expected revoked token error")
	}
}
