package settings

import (
	"context"
	"testing"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/testutil"
)

func TestGetAndResetSettings(t *testing.T) {
	ctx := context.Background()
	log := &testutil.MockLogger{}
	initial := &domain.Settings{
		Id:                        0,
		RequireTwoFactor:          true,
		TokenTTLMinutes:           30,
		MinPasswordLength:         10,
		RequirePasswordComplexity: true,
		UpdatedAt:                 time.Now(),
	}
	repo := testutil.NewMemorySettingsRepo(initial)

	uc := NewUsecase(log, repo)

	got, err := uc.GetSettings(ctx)
	if err != nil {
		t.Fatalf("get settings: %v", err)
	}
	if got.TokenTTLMinutes != 30 || !got.RequireTwoFactor {
		t.Fatalf("unexpected settings: %+v", got)
	}

	if _, err := uc.ResetSettings(ctx); err != nil {
		t.Fatalf("reset settings: %v", err)
	}
	reset, _ := uc.GetSettings(ctx)
	if reset.RequireTwoFactor {
		t.Fatal("expected default requireTwoFactor=false after reset")
	}
}

func TestUpdateSettings(t *testing.T) {
	ctx := context.Background()
	log := &testutil.MockLogger{}
	repo := testutil.NewMemorySettingsRepo(testutil.DefaultSettings())

	uc := NewUsecase(log, repo)

	requireTwoFactor := true
	tokenTTL := int32(120)
	minLen := int32(12)
	requireComplexity := false

	updated, err := uc.UpdateSettings(ctx, &requireTwoFactor, &tokenTTL, &minLen, &requireComplexity)
	if err != nil {
		t.Fatalf("update settings: %v", err)
	}
	if updated.RequireTwoFactor != requireTwoFactor || updated.TokenTTLMinutes != tokenTTL || updated.MinPasswordLength != minLen || updated.RequirePasswordComplexity != requireComplexity {
		t.Fatalf("settings not updated correctly: %+v", updated)
	}
}
