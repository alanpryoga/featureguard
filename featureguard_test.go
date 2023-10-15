package featureguard

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewFeatureGuard(t *testing.T) {
	db, _ := redismock.NewClientMock()

	want := &FeatureGuard{
		db: db,
	}
	got := NewFeatureGuard(db)

	assert.Equal(t, want, got)
}
