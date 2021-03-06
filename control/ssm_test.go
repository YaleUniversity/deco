package control

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

var testDecoString string = `
{
    "filters": {
        "test/file1": {
            "string1": "value1",
            "string2": "value2",
            "string3": "value3"
        },
        "test/file2": {
            "string1": "othervalue1"
        }
    }
}
`

var testParam = ssm.Parameter{
	ARN:   aws.String("arn:aws:ssm:us-east-1:846761448161:parameter/spinup/testapi/dev/deco.json"),
	Value: aws.String(testDecoString),
}

func (m *mockSSMClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &ssm.GetParameterOutput{
		Parameter: &testParam,
	}, nil
}

// mockSSMClient is a fake ssm client
type mockSSMClient struct {
	ssmiface.SSMAPI
	t   *testing.T
	err error
}

func newmockSSMClient(t *testing.T, err error) ssmiface.SSMAPI {
	return &mockSSMClient{
		t:   t,
		err: err,
	}
}

func TestNewSSM(t *testing.T) {
	s := NewSSM()
	to := reflect.TypeOf(s).String()
	if to != "*control.SSM" {
		t.Errorf("expected type to be '*control.SSM', got %s", to)
	}
}

func TestGetParameter(t *testing.T) {
	p := SSM{Service: newmockSSMClient(t, nil)}
	expected := aws.StringValue(testParam.Value)

	out, err := p.GetParameter("/spinup/testapi/dev/deco.json")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	if expected != out {
		t.Errorf("expected %s, got %s", expected, out)
	}
}
