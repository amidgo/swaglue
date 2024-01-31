package test_test

import (
	"io"
	"testing"

	"github.com/amidgo/swaglue/internal/glue/gluer"
	gluermocks "github.com/amidgo/swaglue/internal/glue/gluer/mocks"
	"github.com/stretchr/testify/require"
)

func Test_Container_AllSuccess(t *testing.T) {
	gluer1 := gluermocks.NewGluer(t)
	gluer2 := gluermocks.NewGluer(t)
	gluer3 := gluermocks.NewGluer(t)

	gluer1.EXPECT().Glue().Return(nil).Once()
	gluer2.EXPECT().Glue().Return(nil).Once()
	gluer3.EXPECT().Glue().Return(nil).Once()

	container := gluer.NewContainer()

	container.AddGluer(gluer1)
	container.AddGluer(gluer2)
	container.AddGluer(gluer3)

	err := container.Glue()
	require.NoError(t, err)
}

func Test_Container_GluerError(t *testing.T) {
	gluer2Error := io.ErrNoProgress

	gluer1 := gluermocks.NewGluer(t)
	gluer2 := gluermocks.NewGluer(t)
	gluer3 := gluermocks.NewGluer(t)

	gluer1.EXPECT().Glue().Return(nil).Once()
	gluer2.EXPECT().Glue().Return(gluer2Error).Once()

	container := gluer.NewContainer()

	container.AddGluer(gluer1)
	container.AddGluer(gluer2)
	container.AddGluer(gluer3)

	err := container.Glue()
	require.ErrorIs(t, err, gluer2Error)
}
