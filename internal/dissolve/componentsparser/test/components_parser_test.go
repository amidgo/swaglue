package test_test

import (
	"bytes"
	_ "embed"
	"io"
	"slices"
	"testing"

	"github.com/amidgo/node/yaml"
	"github.com/amidgo/swaglue/internal/dissolve/componentsparser"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
	"github.com/amidgo/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/schemas/AdmissionOrderDate.yaml
	admissionOrderDate []byte
	//go:embed testdata/schemas/AdmissionOrderNumber.yaml
	admissionOrderNumber []byte
	//go:embed testdata/schemas/AdmissionOrderPageData.yaml
	admissionOrderPageData []byte
	//go:embed testdata/schemas/AttendanceReport.yaml
	attendanceReport []byte
	//go:embed testdata/schemas/AttendanceReportMark.yaml
	attendanceReportMark []byte
	//go:embed testdata/schemas/CreateStudentAttendanceRow.yaml
	createStudentAttendanceRow []byte
	//go:embed testdata/schemas/Curriculum.yaml
	curriculum []byte
	//go:embed testdata/schemas/GroupAttendance.yaml
	groupAttendance []byte
	//go:embed testdata/schemas/ID32.yaml
	id32 []byte
	//go:embed testdata/schemas/PeriodAttendance.yaml
	periodAttendance []byte
)

var (
	//go:embed testdata/securitySchemes/adminAuth.yaml
	adminAuth []byte
	//go:embed testdata/securitySchemes/bearerTokenAuth.yaml
	bearerTokenAuth []byte
	//go:embed testdata/securitySchemes/cookieAuth.yaml
	cookieAuth []byte
)

type ComponentsParserTester struct {
	ComponentName string
	Head          *head.Head

	ExpectedComponentItems []*model.Item
	ExpectedErr            error
}

func (c *ComponentsParserTester) Name() string {
	return "component parser tester, component " + c.ComponentName
}

func (c *ComponentsParserTester) Test(t *testing.T) {
	componentParser := componentsparser.NewComponentsParser(
		c.Head,
		c.ComponentName,
		&yaml.Encoder{Indent: 2},
	)

	err := componentParser.Parse()
	require.ErrorIs(t, err, c.ExpectedErr)

	components := componentParser.Components()
	c.assertComponentsEqual(t, components)
}

func (c *ComponentsParserTester) assertComponentsEqual(t *testing.T, components []*model.Item) {
	t.Helper()

	equal := slices.EqualFunc(
		c.ExpectedComponentItems,
		components,
		func(expected, actual *model.Item) bool {
			expectedData, err := io.ReadAll(expected.Content)
			require.NoError(t, err)
			actualData, err := io.ReadAll(actual.Content)
			require.NoError(t, err)

			assert.Equal(t, expected.Name, actual.Name)
			assert.Equal(t, expectedData, actualData)

			return true
		},
	)

	assert.True(t, equal, "slices not equal")
}

func Test_ComponentsParser(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err)

	headWithEmptyComponent, err := head.ParseHeadFromFile("testdata/invalid.swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err)

	tester.RunNamedTesters(t,
		&ComponentsParserTester{
			ComponentName: "schemas",
			Head:          hd,

			ExpectedComponentItems: []*model.Item{
				{
					Name:    "AdmissionOrderDate",
					Content: bytes.NewReader(admissionOrderDate),
				},
				{
					Name:    "AdmissionOrderNumber",
					Content: bytes.NewReader(admissionOrderNumber),
				},
				{
					Name:    "AdmissionOrderPageData",
					Content: bytes.NewReader(admissionOrderPageData),
				},
				{
					Name:    "AttendanceReport",
					Content: bytes.NewReader(attendanceReport),
				},
				{
					Name:    "AttendanceReportMark",
					Content: bytes.NewReader(attendanceReportMark),
				},
				{
					Name:    "CreateStudentAttendanceRow",
					Content: bytes.NewReader(createStudentAttendanceRow),
				},
				{
					Name:    "GroupAttendance",
					Content: bytes.NewReader(groupAttendance),
				},
				{
					Name:    "PeriodAttendance",
					Content: bytes.NewReader(periodAttendance),
				},
				{
					Name:    "ID32",
					Content: bytes.NewReader(id32),
				},
				{
					Name:    "Curriculum",
					Content: bytes.NewReader(curriculum),
				},
			},
		},
		&ComponentsParserTester{
			ComponentName: "securitySchemes",
			Head:          hd,

			ExpectedComponentItems: []*model.Item{
				{
					Name:    "adminAuth",
					Content: bytes.NewReader(adminAuth),
				},
				{
					Name:    "bearerTokenAuth",
					Content: bytes.NewReader(bearerTokenAuth),
				},
				{
					Name:    "cookieAuth",
					Content: bytes.NewReader(cookieAuth),
				},
			},
		},
		&ComponentsParserTester{
			ComponentName: "notExistsComponent",
			Head:          hd,

			ExpectedErr: componentsparser.ErrComponentNotFound,
		},
		&ComponentsParserTester{
			ComponentName: "emptyComponent",
			Head:          headWithEmptyComponent,

			ExpectedErr: componentsparser.ErrEmptyComponents,
		},
	)
}
