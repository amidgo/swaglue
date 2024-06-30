package componentsappender_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/amidgo/node"
	"github.com/amidgo/node/yaml"
	"github.com/amidgo/swaglue/internal/glue/componentsappender"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/item"
	"github.com/amidgo/swaglue/internal/model"
	"github.com/amidgo/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/schemas/EducationPeriod.yaml
	educationPeriodSchema []byte
	//go:embed testdata/schemas/EducationPeriodData.yaml
	educationPeriodData []byte

	//go:embed testdata/schemas/component_expected_swagger.yaml
	componentExpectedSwagger string

	//go:embed testdata/schemas/exists_component_expected_swagger.yaml
	existsComponentExpectedSwagger string

	//go:embed testdata/schemas/two_components_expected.swagger.yaml
	twoComponentsExpectedSwagger string
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

	assert.Equal(t, componentExpectedSwagger, buf.String())
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

	assert.Equal(t, existsComponentExpectedSwagger, buf.String())
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

type IterationStepTest struct {
	CaseName        string
	IterationsSteps []componentsappender.ComponentIterationStep
	ExpectedError   error
	ExpectedData    string
}

func (i *IterationStepTest) Name() string {
	return i.CaseName
}

func (i *IterationStepTest) Test(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/schemas/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open swagger.yaml")

	iterStep := componentsappender.NewIterationStep(i.IterationsSteps...)

	genNode := node.NewIterationMapSource(hd.Node(), iterStep)

	mapNode, err := genNode.MapNode()
	require.ErrorIs(t, err, i.ExpectedError)

	buf := &bytes.Buffer{}

	enc := yaml.Encoder{Indent: 2}
	err = enc.EncodeTo(buf, mapNode)
	require.NoError(t, err)

	assert.Equal(t, i.ExpectedData, buf.String())
}

func Test_IterationStep(t *testing.T) {
	dec := new(yaml.Decoder)

	tester.RunNamedTesters(t,
		&IterationStepTest{
			CaseName: "new component",
			IterationsSteps: []componentsappender.ComponentIterationStep{
				componentsappender.NewComponentIterationStep(
					"schemas",
					item.SliceSource{
						{
							Name:    "EducationPeriod",
							Content: node.MustDecode(dec, educationPeriodSchema),
						},
						{
							Name:    "EducationPeriodData",
							Content: node.MustDecode(dec, educationPeriodData),
						},
					},
				),
			},
			ExpectedData: componentExpectedSwagger,
		},
		&IterationStepTest{
			CaseName: "exists component",
			IterationsSteps: []componentsappender.ComponentIterationStep{
				componentsappender.NewComponentIterationStep(
					"random_component",
					item.SliceSource{
						{
							Name:    "EducationPeriod",
							Content: node.MustDecode(dec, educationPeriodSchema),
						},
						{
							Name:    "EducationPeriodData",
							Content: node.MustDecode(dec, educationPeriodData),
						},
					},
				),
			},
			ExpectedData: existsComponentExpectedSwagger,
		},
		&IterationStepTest{
			CaseName: "two components",
			IterationsSteps: []componentsappender.ComponentIterationStep{
				componentsappender.NewComponentIterationStep(
					"schemas",
					item.SliceSource{
						{
							Name:    "EducationPeriod",
							Content: node.MustDecode(dec, educationPeriodSchema),
						},
						{
							Name:    "EducationPeriodData",
							Content: node.MustDecode(dec, educationPeriodData),
						},
					},
				),
				componentsappender.NewComponentIterationStep(
					"random_component",
					item.SliceSource{
						{
							Name:    "EducationPeriod",
							Content: node.MustDecode(dec, educationPeriodSchema),
						},
						{
							Name:    "EducationPeriodData",
							Content: node.MustDecode(dec, educationPeriodData),
						},
					},
				),
			},
			ExpectedData: twoComponentsExpectedSwagger,
		},
	)
}
