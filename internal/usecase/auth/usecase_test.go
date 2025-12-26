package auth

import (
	"context"
	"testing"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/jwtmanager"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterAndLogin(t *testing.T) {
	ctx := context.Background()
	logger := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())
	email := &testutil.MemoryEmailUsecase{}

	uc := NewUsecase(logger, users, tokens, settings, email, "secret")
	uc.now = func() time.Time { return time.Unix(1000, 0) }

	res, err := uc.Register(ctx, "user@example.com", "pass1234")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if res.User.Email != "user@example.com" || res.User.PasswordHash != "" {
		t.Fatalf("unexpected user: %+v", res.User)
	}

	loginRes, err := uc.Login(ctx, "user@example.com", "pass1234")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if loginRes.Auth.User.Email != "user@example.com" {
		t.Fatalf("unexpected login user: %+v", loginRes.Auth.User)
	}
	if loginRes.Auth.AccessToken == "" {
		t.Fatal("expected token")
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	logger := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())

	uc := NewUsecase(logger, users, tokens, settings, nil, "secret")

	if _, err := uc.Register(ctx, "dup@example.com", "pass1234"); err != nil {
		t.Fatalf("first register: %v", err)
	}
	if _, err := uc.Register(ctx, "dup@example.com", "pass1234"); err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestLoginFailures(t *testing.T) {
	ctx := context.Background()
	logger := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	_ = users.Create(ctx, &domain.User{Email: "user@example.com", PasswordHash: string(hash)})

	uc := NewUsecase(logger, users, tokens, settings, nil, "secret")

	if _, err := uc.Login(ctx, "user@example.com", "wrong"); err == nil {
		t.Fatal("expected invalid credentials")
	}

	user, _ := users.GetByEmail(ctx, "user@example.com")
	user.Disabled = true
	_ = users.Update(ctx, user)

	if _, err := uc.Login(ctx, "user@example.com", "secret"); err == nil {
		t.Fatal("expected disabled error")
	}
}

func TestLogoutRevokesToken(t *testing.T) {
	ctx := context.Background()
	logger := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())

	uc := NewUsecase(logger, users, tokens, settings, nil, "secret")
	token := "token123"
	_ = tokens.Create(ctx, &domain.AccessToken{
		Token:     token,
		UserId:    1,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	})

	if err := uc.Logout(ctx, token); err != nil {
		t.Fatalf("logout: %v", err)
	}

	stored, _ := tokens.GetByToken(ctx, token)
	if !stored.Revoked {
		t.Fatal("token should be revoked")
	}
}

func TestTokenIssuedWithConfiguredSecret(t *testing.T) {
	ctx := context.Background()
	logger := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())

	uc := NewUsecase(logger, users, tokens, settings, nil, "secret-123")
	uc.now = func() time.Time { return time.Unix(2000, 0) }

	if err := users.Create(ctx, &domain.User{Email: "user@example.com", PasswordHash: testutil.PasswordHash("pass1234")}); err != nil {
		t.Fatalf("create user: %v", err)
	}
	user, _ := users.GetByEmail(ctx, "user@example.com")

	authRes, err := uc.issueAuth(ctx, user)
	if err != nil {
		t.Fatalf("issueAuth: %v", err)
	}

	jwt := jwtmanager.New("secret-123")
	claims, err := jwt.Verify(authRes.AccessToken)
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if claims.UserID != user.Id {
		t.Fatalf("unexpected user id in token: %d", claims.UserID)
	}
}
