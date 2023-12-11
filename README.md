# FeatureGuard
FeatureGuard is a feature flags library that uses Redis as the storage backend for Go. It provides a simple interface to enable, disable, toggle, and check the status of feature flags.

## Getting Started
These instructions will help you integrate FeatureGuard into your Go project.

### Prerequisites
- Redis

### Installation
```
go get -u github.com/alanpryoga/featureguard
```

## Usage
```go
package main

import (
	"context"
	"fmt"
	"github.com/alanpryoga/featureguard"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password
		DB:       0,  // Default DB
	})

	// Initialize FeatureGuard
	featureGuard := featureguard.NewFeatureGuard(redisClient, "feature:%s")

	// Enable a feature
	err := featureGuard.EnableFeature(context.Background(), "my_feature")
	if err != nil {
		fmt.Println("Error enabling feature:", err)
		return
	}

	// Check if the feature is enabled
	enabled, err := featureGuard.IsFeatureEnabled(context.Background(), "my_feature")
	if err != nil {
		fmt.Println("Error checking feature status:", err)
		return
	}

	if enabled {
		fmt.Println("My feature is enabled!")
	} else {
		fmt.Println("My feature is disabled.")
	}

	// Disable the feature
	err = featureGuard.DisableFeature(context.Background(), "my_feature")
	if err != nil {
		fmt.Println("Error disabling feature:", err)
		return
	}

	// Toggle the feature state
	err = featureGuard.ToggleFeature(context.Background(), "my_feature")
	if err != nil {
		fmt.Println("Error toggling feature:", err)
		return
	}
}
```

## Contributing

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
