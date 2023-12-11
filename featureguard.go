package featureguard

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const defaultKeyPattern = "featureguard:%s"

// FeatureGuard represents a feature flags management system using Redis as the storage backend.
type FeatureGuard struct {
	db         *redis.Client
	keyPattern string
}

// NewFeatureGuard initializes a new FeatureGuard instance with the provided Redis client connection.
// If keyPattern is empty, it uses the default key pattern.
func NewFeatureGuard(db *redis.Client, keyPattern string) *FeatureGuard {
	if keyPattern == "" {
		keyPattern = defaultKeyPattern
	}

	return &FeatureGuard{
		db:         db,
		keyPattern: keyPattern,
	}
}

// getKey formats the feature key using the specified pattern.
func (fg *FeatureGuard) getKey(feature string) string {
	return fmt.Sprintf(fg.keyPattern, feature)
}

// EnableFeature enables the specified feature flag.
func (fg *FeatureGuard) EnableFeature(ctx context.Context, feature string) error {
	key := fg.getKey(feature)
	return fg.db.Set(ctx, key, true, 0).Err()
}

// DisableFeature disables the specified feature flag.
func (fg *FeatureGuard) DisableFeature(ctx context.Context, feature string) error {
	key := fg.getKey(feature)
	return fg.db.Del(ctx, key).Err()
}

// IsFeatureEnabled checks if the specified feature flag is enabled.
func (fg *FeatureGuard) IsFeatureEnabled(ctx context.Context, feature string) (bool, error) {
	key := fg.getKey(feature)

	result, err := fg.db.Get(ctx, key).Bool()
	if err != nil {
		if err == redis.Nil {
			// Key does not exist, feature is considered disabled
			return false, nil
		}

		// Other error occurred
		return false, err
	}

	return result, nil
}

// ToggleFeature toggles the state of the specified feature flag.
func (fg *FeatureGuard) ToggleFeature(ctx context.Context, feature string) error {
	key := fg.getKey(feature)

	// Check the current state
	currentState, err := fg.IsFeatureEnabled(ctx, feature)
	if err != nil {
		return err
	}

	// Toggle the state
	return fg.db.Set(ctx, key, !currentState, 0).Err()
}
