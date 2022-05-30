// Package worker contains worker instances
package worker

// Work worker pool work instance
type Work interface {
	Do()
	Stop()
}
