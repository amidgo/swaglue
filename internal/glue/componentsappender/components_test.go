package componentsappender_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/amidgo/node/yaml"
	"github.com/amidgo/swaglue/internal/glue/componentsappender"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
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
	hd, err := head.ParseHeadFromFile("testdata/schemas/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open swagger.yaml")

	appender := componentsappender.New(hd, new(yaml.Decoder))

	err = appender.AppendComponent("schemas", []model.Item{
		{
			Name:    "EducationPeriod",
			Content: bytes.NewReader(educationPeriodSchema),
		},
		{
			Name:    "EducationPeriodData",
			Content: bytes.NewReader(educationPeriodData),
		},
	})
	require.NoError(t, err, "append components")

	buf := &bytes.Buffer{}
	err = hd.SaveTo(buf, &yaml.Encoder{Indent: 2})
	require.NoError(t, err, "save file")

	assert.Equal(t, componentExpectedSwagger, buf.Bytes())
}

func Test_Head_AppendComponent_ExistsComponent(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/schemas/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open swagger.yaml")

	appender := componentsappender.New(hd, new(yaml.Decoder))

	err = appender.AppendComponent("random_component", []model.Item{
		{
			Name:    "EducationPeriod",
			Content: bytes.NewReader(educationPeriodSchema),
		},
		{
			Name:    "EducationPeriodData",
			Content: bytes.NewReader(educationPeriodData),
		},
	})
	require.NoError(t, err, "append components")

	buf := &bytes.Buffer{}
	err = hd.SaveTo(buf, &yaml.Encoder{Indent: 2})
	require.NoError(t, err, "save file")

	assert.Equal(t, string(existsComponentExpectedSwagger), buf.String())
}

func Test_Head_AppendComponent_ExistsComponentItemName(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/schemas/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open swagger.yaml")

	appender := componentsappender.New(hd, new(yaml.Decoder))

	err = appender.AppendComponent("random_component", []model.Item{
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

	require.ErrorIs(t, err, componentsappender.ErrComponentItemNameExists)
}
