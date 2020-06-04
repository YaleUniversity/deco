package control

import (
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
func NewSSM(opts ...SSMOption) (*SSM, error) {
	Logger.Println("[INFO] creating new ssm provider")

	s := SSM{}
	s.config = aws.NewConfig()

	for _, opt := range opts {
		opt(&s)
	}

	sess := session.Must(session.NewSession(s.config))

	s.Service = ssm.New(sess)
	return &s, nil
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
