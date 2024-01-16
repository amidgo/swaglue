package head_test

import (
	"bytes"
	_ "embed"
	"testing"

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
)

//go:embed testdata/test_component_expected_swagger.yaml
var componentExpectedSwagger []byte

func Test_Head_AppendComponent(t *testing.T) {
	head, err := head.ParseHeadFromFile("testdata/swagger.yaml")
	require.NoError(t, err, "failed open swagger.yaml")

	err = head.AppendComponent("schemas", []*model.SwaggerComponentItem{
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
