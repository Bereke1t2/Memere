package auth

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// fakeUserRepo is an in-memory UserRepository for usecase tests. It is keyed by
// id and indexed by lowercased email, mimicking the unique-email constraint.
type fakeUserRepo struct {
	mu        sync.Mutex
	byID      map[uuid.UUID]*entity.User
	createErr error
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{byID: map[uuid.UUID]*entity.User{}}
}

func (f *fakeUserRepo) Create(_ context.Context, u *entity.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.createErr != nil {
		return f.createErr
	}
	for _, ex := range f.byID {
		if ex.Email == u.Email {
			return apperror.Conflict("EMAIL_TAKEN", nil)
		}
	}
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	cp := *u
	f.byID[u.ID] = &cp
	return nil
}

func (f *fakeUserRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[id]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, apperror.NotFound("user not found", nil)
}

func (f *fakeUserRepo) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.byID {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, apperror.NotFound("user not found", nil)
}

func (f *fakeUserRepo) Update(_ context.Context, u *entity.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.byID[u.ID]; !ok {
		return apperror.NotFound("user not found", nil)
	}
	cp := *u
	f.byID[u.ID] = &cp
	return nil
}

func (f *fakeUserRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.byID, id)
	return nil
}

func (f *fakeUserRepo) SetLastLogin(_ context.Context, id uuid.UUID, t time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[id]; ok {
		u.LastLoginAt = &t
	}
	return nil
}

// fakeTokenRepo is an in-memory RefreshTokenRepository keyed by hash.
type fakeTokenRepo struct {
	mu     sync.Mutex
	byHash map[string]*entity.RefreshToken
	byID   map[uuid.UUID]*entity.RefreshToken
}

func newFakeTokenRepo() *fakeTokenRepo {
	return &fakeTokenRepo{
		byHash: map[string]*entity.RefreshToken{},
		byID:   map[uuid.UUID]*entity.RefreshToken{},
	}
}

func (f *fakeTokenRepo) Create(_ context.Context, t *entity.RefreshToken) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	t.CreatedAt = time.Now()
	cp := *t
	f.byHash[t.TokenHash] = &cp
	f.byID[t.ID] = &cp
	return nil
}

func (f *fakeTokenRepo) FindByHash(_ context.Context, hash string) (*entity.RefreshToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if t, ok := f.byHash[hash]; ok {
		cp := *t
		return &cp, nil
	}
	return nil, apperror.NotFound("refresh token not found", nil)
}

func (f *fakeTokenRepo) Revoke(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if t, ok := f.byID[id]; ok {
		now := time.Now()
		t.RevokedAt = &now
	}
	return nil
}

func (f *fakeTokenRepo) RevokeAllForUser(_ context.Context, userID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	for _, t := range f.byID {
		if t.UserID == userID {
			t.RevokedAt = &now
		}
	}
	return nil
}

func (f *fakeTokenRepo) DeleteExpired(_ context.Context) error { return nil }

// fakeSessionRepo is an in-memory SessionRepository keyed by user id.
type fakeSessionRepo struct {
	mu     sync.Mutex
	byUser map[uuid.UUID]string
}

func newFakeSessionRepo() *fakeSessionRepo {
	return &fakeSessionRepo{byUser: map[uuid.UUID]string{}}
}

func (f *fakeSessionRepo) SetSession(_ context.Context, userID uuid.UUID, tokenHash string, _ time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.byUser[userID] = tokenHash
	return nil
}

func (f *fakeSessionRepo) GetSession(_ context.Context, userID uuid.UUID) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.byUser[userID], nil
}

func (f *fakeSessionRepo) DeleteSession(_ context.Context, userID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.byUser, userID)
	return nil
}
