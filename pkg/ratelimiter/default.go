/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package ratelimiter contains suggested default ratelimiters for Crossplane.
package ratelimiter

import (
	"time"

	"golang.org/x/time/rate"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/ratelimiter"
)

const (
	// DefaultGlobalRPS is the recommended default average requeues per
	// second tolerated by Crossplane controller managers.
	DefaultGlobalRPS = 1

	// DefaultProviderRPS is the recommended default average requeues per
	// second tolerated by a Crossplane provider.
	//
	// Deprecated: Use DefaultGlobalRPS
	DefaultProviderRPS = DefaultGlobalRPS
)

// NewGlobal returns a token bucket rate limiter meant for limiting the number
// of average total requeues per second for all controllers registered with a
// controller manager. The bucket size is a linear function of the requeues per
// second.
func NewGlobal(rps int) *workqueue.BucketRateLimiter {
	return &workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(rps), rps*10)}
}

// NewController returns a rate limiter that takes the maximum delay between the
// passed rate limiter and a per-item exponential backoff limiter. The
// exponential backoff limiter has a base delay of 1s and a maximum of 60s.
func NewController(global ratelimiter.RateLimiter) ratelimiter.RateLimiter {
	return workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(1*time.Second, 60*time.Second),
		global,
	)
}

// NewDefaultProviderRateLimiter returns a token bucket rate limiter meant for
// limiting the number of average total requeues per second for all controllers
// registered with a controller manager. The bucket size is a linear function of
// the requeues per second.
//
// Deprecated: Use NewGlobal.
func NewDefaultProviderRateLimiter(rps int) *workqueue.BucketRateLimiter {
	return NewGlobal(rps)
}

// NewDefaultManagedRateLimiter returns a rate limiter that takes the maximum
// delay between the passed provider and a per-item exponential backoff limiter.
// The exponential backoff limiter has a base delay of 1s and a maximum of 60s.
//
// Deprecated: Use NewController.
func NewDefaultManagedRateLimiter(provider ratelimiter.RateLimiter) ratelimiter.RateLimiter {
	return NewController(provider)
}
