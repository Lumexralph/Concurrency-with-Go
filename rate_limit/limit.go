// Package main is an implementation of the token bucket algorithm to
// rate limit requests.
//
// It assumes there's an access to an API and a Go client has been
// created to utilize it. The API has two endpoints, one for reading
// a file, and another to resolve domain name to an IP address
package main

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

// RateLimiter to allow MultiLimiter define other MultiLimiter instances.
type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

// Multilimiter combines separate limiters into one rate limiter.
func Multilimiter(limiters ...RateLimiter) *multiLimiter {
	// function to sort the limiters by Limit in ascending order
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}

	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

// Per - rate limits in terms of number operations per time measurement.
func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

// APIConnection -
type APIConnection struct {
	apiLimit,
	diskLimit,
	networkLimit RateLimiter
}

// Open - initiates the API connection with multiple rate limits
func Open() *APIConnection {

	return &APIConnection{
		apiLimit: Multilimiter(
			// limit per sec
			rate.NewLimiter(Per(2, time.Second), 1),
			// limit per min
			rate.NewLimiter(Per(10, time.Minute), 10),
		),
		diskLimit: Multilimiter(
			// one read pe sec
			rate.NewLimiter(rate.Limit(1), 1),
		),
		networkLimit: Multilimiter(
			// 3 requests per secs
			rate.NewLimiter(Per(3, time.Second), 3),
		),
	}
}

// ReadFile takes context as first parameter in case we need
// to cancel the request or pass values over to the server.
func (a *APIConnection) ReadFile(ctx context.Context) error {
	// apply the rate limiter for every request
	err := Multilimiter(a.apiLimit, a.diskLimit).Wait(ctx)
	if err != nil {
		return err
	}
	// perform some work
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	// apply the rate limiter for every request
	// wait for it to have enough access token to complete the request
	err := Multilimiter(a.apiLimit, a.networkLimit).Wait(ctx)
	if err != nil {
		return err
	}
	// perform some work
	return nil
}

type multiLimiter struct {
	limiters []RateLimiter
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	// return the most restrictive limit i.e smallest limit
	return l.limiters[0].Limit()
}

func main() {
	defer log.Print("Done")
	// location to direct the logs
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot read file: %v", err)
			}
			log.Print("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot resolve url: %v", err)
			}
			log.Print("ResolveAddress")
		}()
	}

	wg.Wait()
}
