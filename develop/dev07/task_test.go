package main

import (
	"testing"
	"time"
)

func TestOrSingleChannel(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(sig(1 * time.Second))
	elapsed := time.Since(start)

	if elapsed < 1*time.Second {
		t.Errorf("Expected at least 1 second, but got %v", elapsed)
	}
}

func TestOrMultipleChannels(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	elapsed := time.Since(start)

	if elapsed >= 1*time.Second && elapsed < 2*time.Second {
		t.Logf("Test passed: done after %v", elapsed)
	} else {
		t.Errorf("Expected close in approximately 1 second, but got %v", elapsed)
	}
}

func TestOrImmediateChannel(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(sig(0), sig(2*time.Second))
	elapsed := time.Since(start)

	if elapsed >= 0 && elapsed < 1*time.Second {
		t.Logf("Test passed: done immediately after %v", elapsed)
	} else {
		t.Errorf("Expected immediate close, but got %v", elapsed)
	}
}

func TestOrNoChannels(t *testing.T) {
	start := time.Now()
	done := make(chan interface{})

	go func() {
		time.Sleep(500 * time.Millisecond)
		close(done)
	}()

	<-or(done)
	elapsed := time.Since(start)

	if elapsed >= 500*time.Millisecond && elapsed < 1*time.Second {
		t.Logf("Test passed: done after %v", elapsed)
	} else {
		t.Errorf("Expected close in approximately 500 milliseconds, but got %v", elapsed)
	}
}
