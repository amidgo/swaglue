package gluer

import "fmt"

//go:generate mockery --all --output mocks --outpkg gluermocks --with-expecter

type Gluer interface {
	Glue() error
}

type Container struct {
	gluers []Gluer
}

func NewContainer() Container {
	return Container{
		gluers: make([]Gluer, 0),
	}
}

func (c *Container) AddGluer(gluer Gluer) {
	c.gluers = append(c.gluers, gluer)
}

func (c *Container) Glue() error {
	for _, gluer := range c.gluers {
		err := gluer.Glue()
		if err != nil {
			return fmt.Errorf("failed glue, %w", err)
		}
	}

	return nil
}
