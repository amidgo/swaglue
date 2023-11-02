package head_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
	"gotest.tools/v3/assert"
)

var (
	educationPeriodSchema = []byte{116, 121, 112, 101, 58, 32, 111, 98, 106, 101, 99, 116, 10, 112, 114, 111, 112, 101, 114, 116, 105, 101, 115, 58, 10, 32, 32, 105, 100, 58, 10, 32, 32, 32, 32, 36, 114, 101, 102, 58, 32, 34, 46, 46, 47, 73, 68, 51, 50, 46, 121, 97, 109, 108, 34, 10, 32, 32, 115, 116, 97, 114, 116, 95, 100, 97, 116, 101, 58, 10, 32, 32, 32, 32, 36, 114, 101, 102, 58, 32, 34, 69, 100, 117, 99, 97, 116, 105, 111, 110, 80, 101, 114, 105, 111, 100, 83, 116, 97, 114, 116, 68, 97, 116, 101, 46, 121, 97, 109, 108, 34, 10, 32, 32, 101, 110, 100, 95, 100, 97, 116, 101, 58, 10, 32, 32, 32, 32, 36, 114, 101, 102, 58, 32, 34, 69, 100, 117, 99, 97, 116, 105, 111, 110, 80, 101, 114, 105, 111, 100, 69, 110, 100, 68, 97, 116, 101, 46, 121, 97, 109, 108, 34, 10}
	educationPeriodData   = []byte{116, 121, 112, 101, 58, 32, 111, 98, 106, 101, 99, 116, 10, 112, 114, 111, 112, 101, 114, 116, 105, 101, 115, 58, 10, 32, 32, 115, 116, 97, 114, 116, 95, 100, 97, 116, 101, 58, 10, 32, 32, 32, 32, 36, 114, 101, 102, 58, 32, 34, 69, 100, 117, 99, 97, 116, 105, 111, 110, 80, 101, 114, 105, 111, 100, 83, 116, 97, 114, 116, 68, 97, 116, 101, 46, 121, 97, 109, 108, 34, 10, 32, 32, 101, 110, 100, 95, 100, 97, 116, 101, 58, 10, 32, 32, 32, 32, 36, 114, 101, 102, 58, 32, 34, 69, 100, 117, 99, 97, 116, 105, 111, 110, 80, 101, 114, 105, 111, 100, 69, 110, 100, 68, 97, 116, 101, 46, 121, 97, 109, 108, 34, 10}
)

func TestComponents(t *testing.T) {
	head, err := head.ParseHeadFromFile("testdata/swagger.yaml")
	assert.NilError(t, err, "failed open swagger.yaml")
	head.AppendComponent("schemas", []*model.Component{
		{
			Name:    "EducationPeriod",
			Content: bytes.NewReader(educationPeriodSchema),
		},
		{
			Name:    "EducationPeriodData",
			Content: bytes.NewReader(educationPeriodData),
		},
	})
	s := &strings.Builder{}
	err = head.SaveTo(s)
	assert.NilError(t, err, "failed save head")
	log.Println(s.String())
}
