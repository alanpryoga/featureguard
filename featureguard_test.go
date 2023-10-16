package featureguard

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewFeatureGuard(t *testing.T) {
	redisClientMock, _ := redismock.NewClientMock()

	type args struct {
		redisClient    *redis.Client
		redisKeyPrefix string
	}
	tests := []struct {
		name string
		args args
		want *FeatureGuard
	}{
		{
			name: "Given initialization a new FeatureGuard instance " +
				"with the Redis client connection and no key prefix provided " +
				"Then it should return FeatureGuard instance with default key prefix",
			args: args{
				redisClient:    redisClientMock,
				redisKeyPrefix: "",
			},
			want: &FeatureGuard{
				client:    redisClientMock,
				keyPrefix: defaultRedisKeyPrefix,
			},
		},
		{
			name: "Given initialization a new FeatureGuard instance " +
				"with the Redis client connection and key prefix provided " +
				"Then it should return FeatureGuard instance with provided key prefix",
			args: args{
				redisClient:    redisClientMock,
				redisKeyPrefix: "prefix",
			},
			want: &FeatureGuard{
				client:    redisClientMock,
				keyPrefix: "prefix",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFeatureGuard(tt.args.redisClient, tt.args.redisKeyPrefix)

			assert.Equal(t, tt.want, got)
		})
	}
}
