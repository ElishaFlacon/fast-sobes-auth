package user

import (
	"context"
	"testing"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/testutil"
)

func TestChangePassword(t *testing.T) {
	ctx := context.Background()
	log := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())
	email := &testutil.MemoryEmailUsecase{}

	oldHash := testutil.PasswordHash("oldpass1")
	_ = users.Create(ctx, &domain.User{Email: "user@example.com", PasswordHash: oldHash})
	_ = tokens.Create(ctx, &domain.AccessToken{
		Token:     "t1",
		UserID:    1,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	})

	uc := NewUsecase(log, users, settings, tokens, email)
	if err := uc.ChangePassword(ctx, "1", "oldpass1", "newpass2"); err != nil {
		t.Fatalf("change password: %v", err)
	}

	user, _ := users.GetByID(ctx, 1)
	if user.PasswordHash == oldHash {
		t.Fatal("password hash was not updated")
	}
	token, _ := tokens.GetByToken(ctx, "t1")
	if !token.Revoked {
		t.Fatal("token should be revoked after password change")
	}
}

func TestChangeEmail(t *testing.T) {
	ctx := context.Background()
	log := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())
	email := &testutil.MemoryEmailUsecase{}

	hash := testutil.PasswordHash("pass1234")
	_ = users.Create(ctx, &domain.User{Email: "old@example.com", PasswordHash: hash})
	_ = tokens.Create(ctx, &domain.AccessToken{
		Token:     "t1",
		UserID:    1,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	})

	uc := NewUsecase(log, users, settings, tokens, email)
	if err := uc.ChangeEmail(ctx, "1", "new@example.com", "pass1234"); err != nil {
		t.Fatalf("change email: %v", err)
	}

	user, _ := users.GetByID(ctx, 1)
	if user.Email != "new@example.com" {
		t.Fatalf("email not updated: %s", user.Email)
	}
	token, _ := tokens.GetByToken(ctx, "t1")
	if !token.Revoked {
		t.Fatal("token should be revoked after email change")
	}
	if len(email.Messages) == 0 {
		t.Fatal("expected email stub call")
	}
}

func TestUpdatePermissionsAndStatus(t *testing.T) {
	ctx := context.Background()
	log := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())

	_ = users.Create(ctx, &domain.User{Email: "user@example.com", PermissionLevel: 1, PasswordHash: testutil.PasswordHash("pass1234")})

	uc := NewUsecase(log, users, settings, tokens, &testutil.MemoryEmailUsecase{})

	user, err := uc.UpdatePermissions(ctx, "1", 5)
	if err != nil {
		t.Fatalf("update permissions: %v", err)
	}
	if user.PermissionLevel != 5 {
		t.Fatal("permission level not updated")
	}

	if err := uc.DisableUser(ctx, "1"); err != nil {
		t.Fatalf("disable user: %v", err)
	}
	disabledUser, _ := users.GetByID(ctx, 1)
	if !disabledUser.Disabled {
		t.Fatal("user not disabled")
	}

	if err := uc.EnableUser(ctx, "1"); err != nil {
		t.Fatalf("enable user: %v", err)
	}
	enabledUser, _ := users.GetByID(ctx, 1)
	if enabledUser.Disabled {
		t.Fatal("user not enabled")
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	log := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())

	_ = users.Create(ctx, &domain.User{Email: "user@example.com", PasswordHash: testutil.PasswordHash("pass1234")})
	_ = tokens.Create(ctx, &domain.AccessToken{
		Token:     "t1",
		UserID:    1,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	})

	uc := NewUsecase(log, users, settings, tokens, &testutil.MemoryEmailUsecase{})

	if err := uc.DeleteUser(ctx, "1"); err != nil {
		t.Fatalf("delete user: %v", err)
	}

	if _, err := users.GetByID(ctx, 1); err == nil {
		t.Fatal("user should be deleted")
	}
	token, _ := tokens.GetByToken(ctx, "t1")
	if !token.Revoked {
		t.Fatal("token should be revoked on delete")
	}
}

func TestGetUserList(t *testing.T) {
	ctx := context.Background()
	log := &testutil.MockLogger{}
	users := testutil.NewMemoryUserRepo()
	tokens := testutil.NewMemoryAccessTokenRepo()
	settings := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())

	_ = users.Create(ctx, &domain.User{Email: "user1@example.com", PermissionLevel: 1, PasswordHash: testutil.PasswordHash("pass123")})
	_ = users.Create(ctx, &domain.User{Email: "user2@example.com", PermissionLevel: 2, PasswordHash: testutil.PasswordHash("pass123")})

	uc := NewUsecase(log, users, settings, tokens, &testutil.MemoryEmailUsecase{})

	list, err := uc.UsersList(ctx, 1, 10, nil, true)
	if err != nil {
		t.Fatalf("users list: %v", err)
	}
	if list.Total != 2 || len(list.Users) != 2 {
		t.Fatalf("unexpected list result: %+v", list)
	}
	for _, u := range list.Users {
		if u.PasswordHash != "" {
			t.Fatal("password hash should be sanitized")
		}
	}
}
