package featureguard

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewFeatureGuard(t *testing.T) {
	db, _ := redismock.NewClientMock()
	keyPattern := "featureguardtest:%s"
	type args struct {
		db         *redis.Client
		keyPattern string
	}
	tests := []struct {
		name string
		args args
		want *FeatureGuard
	}{
		{
			name: "Initializes with default key pattern",
			args: args{
				db:         db,
				keyPattern: "",
			},
			want: &FeatureGuard{
				db:         db,
				keyPattern: defaultKeyPattern,
			},
		},
		{
			name: "Initializes with custom key pattern",
			args: args{
				db:         db,
				keyPattern: keyPattern,
			},
			want: &FeatureGuard{
				db:         db,
				keyPattern: keyPattern,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFeatureGuard(tt.args.db, tt.args.keyPattern)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFeatureGuard_getKey(t *testing.T) {
	type args struct {
		feature string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get valid key",
			args: args{
				feature: "unit",
			},
			want: "featureguardtest:unit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := redismock.NewClientMock()
			keyPattern := "featureguardtest:%s"
			fg := &FeatureGuard{
				db:         db,
				keyPattern: keyPattern,
			}
			got := fg.getKey(tt.args.feature)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFeatureGuard_EnableFeature(t *testing.T) {
	type fields struct {
		db         redismock.ClientMock
		keyPattern string
	}
	type args struct {
		ctx     context.Context
		feature string
	}
	tests := []struct {
		name    string
		args    args
		mock    func(f fields)
		wantErr error
	}{
		{
			name: "Succees to enables a feature",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectSet(key, true, 0).SetVal("")
			},
			wantErr: nil,
		},
		{
			name: "Error when enables a feature",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectSet(key, true, 0).SetErr(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			f := fields{
				db:         mock,
				keyPattern: defaultKeyPattern,
			}
			tt.mock(f)
			fg := &FeatureGuard{
				db:         db,
				keyPattern: f.keyPattern,
			}
			err := fg.EnableFeature(tt.args.ctx, tt.args.feature)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestFeatureGuard_DisableFeature(t *testing.T) {
	type fields struct {
		db         redismock.ClientMock
		keyPattern string
	}
	type args struct {
		ctx     context.Context
		feature string
	}
	tests := []struct {
		name    string
		args    args
		mock    func(f fields)
		wantErr error
	}{
		{
			name: "Succees to disables a feature",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectDel(key).SetVal(0)
			},
			wantErr: nil,
		},
		{
			name: "Error when disables a feature",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectDel(key).SetErr(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			f := fields{
				db:         mock,
				keyPattern: defaultKeyPattern,
			}
			tt.mock(f)
			fg := &FeatureGuard{
				db:         db,
				keyPattern: f.keyPattern,
			}
			err := fg.DisableFeature(tt.args.ctx, tt.args.feature)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestFeatureGuard_IsFeatureEnabled(t *testing.T) {
	type fields struct {
		db         redismock.ClientMock
		keyPattern string
	}
	type args struct {
		ctx     context.Context
		feature string
	}
	tests := []struct {
		name    string
		args    args
		mock    func(f fields)
		want    bool
		wantErr error
	}{
		{
			name: "Get feature is already enabled",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectGet(key).SetVal("true")
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "Get feature is currently disabled",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectGet(key).SetVal("false")
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "Get feature is currently disabled due to no key exists in the db",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectGet(key).SetErr(redis.Nil)
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "Error when check is feature enabled",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectGet(key).SetErr(assert.AnError)
			},
			want:    false,
			wantErr: assert.AnError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			f := fields{
				db:         mock,
				keyPattern: defaultKeyPattern,
			}
			tt.mock(f)
			fg := &FeatureGuard{
				db:         db,
				keyPattern: f.keyPattern,
			}
			got, err := fg.IsFeatureEnabled(tt.args.ctx, tt.args.feature)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestFeatureGuard_ToggleFeature(t *testing.T) {
	type fields struct {
		db         redismock.ClientMock
		keyPattern string
	}
	type args struct {
		ctx     context.Context
		feature string
	}
	tests := []struct {
		name    string
		args    args
		mock    func(f fields)
		wantErr error
	}{
		{
			name: "Toggles feature from disable to enable",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectGet(key).SetVal("false")
				f.db.ExpectSet(key, true, 0).SetVal("")
			},
			wantErr: nil,
		},
		{
			name: "Toggles feature from enable to disable",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectGet(key).SetVal("true")
				f.db.ExpectSet(key, false, 0).SetVal("")
			},
			wantErr: nil,
		},
		{
			name: "Error when check is feature enabled",
			args: args{
				ctx:     context.Background(),
				feature: "unit",
			},
			mock: func(f fields) {
				key := fmt.Sprintf(defaultKeyPattern, "unit")
				f.db.ExpectGet(key).SetErr(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			f := fields{
				db:         mock,
				keyPattern: defaultKeyPattern,
			}
			tt.mock(f)
			fg := &FeatureGuard{
				db:         db,
				keyPattern: f.keyPattern,
			}
			err := fg.ToggleFeature(tt.args.ctx, tt.args.feature)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
