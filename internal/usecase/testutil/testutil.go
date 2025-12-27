package testutil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MockLogger struct {
	Infos  []string
	Errors []string
}

func (l *MockLogger) Infof(format string, v ...any) {
	l.Infos = append(l.Infos, sprintf(format, v...))
}

func (l *MockLogger) Errorf(format string, v ...any) {
	l.Errors = append(l.Errors, sprintf(format, v...))
}

func (l *MockLogger) Fatal(format string, v ...any) {
	l.Errorf(format, v...)
}

func (l *MockLogger) Stop() {}

func DefaultSettings() *domain.Settings {
	return &domain.Settings{
		ID:                        0,
		RequireTwoFactor:          false,
		TokenTTLMinutes:           60,
		MinPasswordLength:         8,
		RequirePasswordComplexity: true,
		UpdatedAt:                 time.Now(),
	}
}

// MemoryUserRepo is a lightweight in-memory implementation for tests.
type MemoryUserRepo struct {
	mu     sync.Mutex
	users  map[int64]*domain.User
	nextID int64
}

func NewMemoryUserRepo() *MemoryUserRepo {
	return &MemoryUserRepo{
		users:  make(map[int64]*domain.User),
		nextID: 1,
	}
}

func (r *MemoryUserRepo) GetByID(_ context.Context, id int64) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	copy := *user
	return &copy, nil
}

func (r *MemoryUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == email {
			copy := *u
			return &copy, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (r *MemoryUserRepo) GetList(
	_ context.Context,
	offset,
	limit int32,
	minPermissionLevel *int32,
	includeDisabled bool,
) ([]*domain.User, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var filtered []*domain.User
	for _, u := range r.users {
		if minPermissionLevel != nil && u.PermissionLevel != *minPermissionLevel {
			continue
		}
		if !includeDisabled && u.Disabled {
			continue
		}
		copy := *u
		filtered = append(filtered, &copy)
	}

	total := int64(len(filtered))
	end := int(offset + limit)
	if end > len(filtered) {
		end = len(filtered)
	}
	start := int(offset)
	if start > len(filtered) {
		start = len(filtered)
	}

	return filtered[start:end], total, nil
}

func (r *MemoryUserRepo) Create(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.nextID++
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	copy := *user
	r.users[user.ID] = &copy
	return nil
}

func (r *MemoryUserRepo) Update(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	copy := *user
	r.users[user.ID] = &copy
	return nil
}

func (r *MemoryUserRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.users, id)
	return nil
}

// MemoryAccessTokenRepo stores tokens in memory for tests.
type MemoryAccessTokenRepo struct {
	mu     sync.Mutex
	tokens map[string]*domain.AccessToken
}

func NewMemoryAccessTokenRepo() *MemoryAccessTokenRepo {
	return &MemoryAccessTokenRepo{
		tokens: make(map[string]*domain.AccessToken),
	}
}

func (r *MemoryAccessTokenRepo) Create(_ context.Context, token *domain.AccessToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *token
	r.tokens[token.Token] = &copy
	return nil
}

func (r *MemoryAccessTokenRepo) GetByToken(_ context.Context, token string) (*domain.AccessToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tokens[token]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	copy := *t
	return &copy, nil
}

func (r *MemoryAccessTokenRepo) Revoke(_ context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t, ok := r.tokens[token]; ok {
		t.Revoked = true
	}
	return nil
}

func (r *MemoryAccessTokenRepo) RevokeAllByUser(_ context.Context, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, t := range r.tokens {
		if t.UserID == userID {
			t.Revoked = true
		}
	}
	return nil
}

func (r *MemoryAccessTokenRepo) DeleteExpired(_ context.Context, now time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for token, t := range r.tokens {
		if t.ExpiresAt.Before(now) {
			delete(r.tokens, token)
		}
	}
	return nil
}

// MemorySettingsRepo stores settings in memory for tests.
type MemorySettingsRepo struct {
	mu       sync.Mutex
	settings *domain.Settings
}

func NewMemorySettingsRepo(settings *domain.Settings) *MemorySettingsRepo {
	return &MemorySettingsRepo{settings: settings}
}

func (r *MemorySettingsRepo) Get(_ context.Context) (*domain.Settings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.settings == nil {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *r.settings
	return &copy, nil
}

func (r *MemorySettingsRepo) Update(_ context.Context, settings *domain.Settings) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *settings
	r.settings = &copy
	return nil
}

func (r *MemorySettingsRepo) Reset(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.settings = &domain.Settings{
		ID:                        0,
		RequireTwoFactor:          false,
		TokenTTLMinutes:           60,
		MinPasswordLength:         8,
		RequirePasswordComplexity: true,
		UpdatedAt:                 time.Now(),
	}
	return nil
}

// MemoryEmailUsecase captures sent messages during tests.
type MemoryEmailUsecase struct {
	Messages []EmailMessage
}

type EmailMessage struct {
	To      string
	Subject string
	Body    string
}

func (m *MemoryEmailUsecase) Send(_ context.Context, to, subject, body string) error {
	m.Messages = append(m.Messages, EmailMessage{To: to, Subject: subject, Body: body})
	return nil
}

// PasswordHash is a helper to prepare hashed passwords in tests.
func PasswordHash(raw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	return string(hash)
}

func sprintf(format string, v ...any) string {
	return fmt.Sprintf(format, v...)
}
