package tester

import "testing"

type Tester interface {
	Test(t *testing.T)
}

type Container struct {
	testers []Tester
}

func (tt *Container) AddTester(t Tester) {
	tt.testers = append(tt.testers, t)
}

func (tt *Container) Test(t *testing.T) {
	for _, tester := range tt.testers {
		if tester == nil {
			continue
		}
		tester.Test(t)
	}
}
