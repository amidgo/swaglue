package head_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/amidgo/swaglue/internal/glue/head"
	"github.com/amidgo/swaglue/internal/glue/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/schemas/EducationPeriod.yaml
	educationPeriodSchema []byte
	//go:embed testdata/schemas/EducationPeriodData.yaml
	educationPeriodData []byte

	//go:embed testdata/schemas/component_expected_swagger.yaml
	componentExpectedSwagger []byte

	//go:embed testdata/schemas/exists_component_expected_swagger.yaml
	existsComponentExpectedSwagger []byte
)

func Test_Head_AppendComponent(t *testing.T) {
	head, err := head.ParseHeadFromFile("testdata/schemas/swagger.yaml")
	require.NoError(t, err, "failed open swagger.yaml")

	err = head.AppendComponent("schemas", []*model.Item{
		{
			Name:    "EducationPeriod",
			Content: bytes.NewReader(educationPeriodSchema),
		},
		{
			Name:    "EducationPeriodData",
			Content: bytes.NewReader(educationPeriodData),
		},
	})
	require.NoError(t, err, "failed append components")

	buf := &bytes.Buffer{}
	err = head.SaveTo(buf)
	require.NoError(t, err, "failed save file")

	assert.Equal(t, componentExpectedSwagger, buf.Bytes())
}

func Test_Head_AppendComponent_ExistsComponent(t *testing.T) {
	head, err := head.ParseHeadFromFile("testdata/schemas/swagger.yaml")
	require.NoError(t, err, "failed open swagger.yaml")

	err = head.AppendComponent("random_component", []*model.Item{
		{
			Name:    "EducationPeriod",
			Content: bytes.NewReader(educationPeriodSchema),
		},
		{
			Name:    "EducationPeriodData",
			Content: bytes.NewReader(educationPeriodData),
		},
	})
	require.NoError(t, err, "failed append components")

	buf := &bytes.Buffer{}
	err = head.SaveTo(buf)
	require.NoError(t, err, "failed save file")

	assert.Equal(t, existsComponentExpectedSwagger, buf.Bytes())
}

func Test_Head_AppendComponent_ExistsComponentItemName(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/schemas/swagger.yaml")
	require.NoError(t, err, "failed open swagger.yaml")

	err = hd.AppendComponent("random_component", []*model.Item{
		{
			Name:    "EducationPeriod",
			Content: bytes.NewReader(educationPeriodSchema),
		},
		{
			Name:    "Hello",
			Content: bytes.NewReader(educationPeriodData),
		},
		{
			Name:    "EducationPeriodData",
			Content: bytes.NewReader(educationPeriodData),
		},
	})

	require.ErrorIs(t, err, head.ErrComponentItemNameExists)
}
