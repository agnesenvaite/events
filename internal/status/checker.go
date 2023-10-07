package status

import (
	"time"
)

const (
	OK   = "OK"
	DOWN = "DOWN"
)
const checkTimeout = 500 * time.Millisecond

type Checker interface {
	Name() string
	Status(timeout time.Duration) string
}

type Pool interface {
	Status() string
	Details() map[string]string
}

type checkResult struct {
	name   string
	status string
}

func NewCheckerPool(checkers ...Checker) *CheckerPool {
	return &CheckerPool{checkers}
}

type CheckerPool struct {
	checkers []Checker
}

func (cp *CheckerPool) Status() string {
	resultChannel := cp.check()

	for i := 0; i < len(cp.checkers); i++ {
		result := <-resultChannel
		if result.status != OK {
			return result.status
		}
	}

	return OK
}

func (cp *CheckerPool) Details() map[string]string {
	details := make(map[string]string)
	resultChannel := cp.check()

	for i := 0; i < len(cp.checkers); i++ {
		result := <-resultChannel

		details[result.name] = result.status
	}

	return details
}

func (cp *CheckerPool) check() <-chan checkResult {
	resultChannel := make(chan checkResult)

	for _, checker := range cp.checkers {
		go func(ch Checker) {
			resultChannel <- checkResult{name: ch.Name(), status: ch.Status(checkTimeout)}
		}(checker)
	}

	return resultChannel
}
