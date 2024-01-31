package tester

import "testing"

type NamedTester interface {
	Name() string
	Test(t *testing.T)
}

type NamedContainer struct {
	namedTesters []NamedTester
}

func (c *NamedContainer) AddNamedTester(tester NamedTester) {
	c.namedTesters = append(c.namedTesters, tester)
}

func (c *NamedContainer) AddNamedTesters(testers ...NamedTester) {
	c.namedTesters = append(c.namedTesters, testers...)
}

func (c *NamedContainer) Test(t *testing.T) {
	for _, tester := range c.namedTesters {
		if tester == nil {
			continue
		}

		t.Run(tester.Name(), tester.Test)
	}
}
