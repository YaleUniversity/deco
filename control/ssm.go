package control

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// SSMOption is a function to set ssm options
type SSMOption func(*SSM)

type SSM struct {
	Service ssmiface.SSMAPI
	config  *aws.Config
}

// New creates an SSM from a list of SSMOption functions
func NewSSM(opts ...SSMOption) *SSM {
	Logger.Println("[INFO] creating new ssm provider")

	region := "us-east-1"
	if r, ok := os.LookupEnv("AWS_REGION"); ok {
		region = r
	} else if r, ok := os.LookupEnv("AWS_DEFAULT_REGION"); ok {
		region = r
	}

	s := SSM{}
	s.config = aws.NewConfig().WithRegion(region)

	for _, opt := range opts {
		opt(&s)
	}

	sess := session.Must(session.NewSession(s.config))
	s.Service = ssm.New(sess)
	return &s
}

func (s *SSM) GetParameter(path string) (string, error) {
	Logger.Println("[INFO] Fetching control from SSM path", path)

	out, err := s.Service.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		Logger.Println("[Error] Failed to fetch SSM paramter", err)
		return "", err
	}

	return aws.StringValue(out.Parameter.Value), nil
}
