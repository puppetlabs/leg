package sns

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/puppetlabs/horsehead/v2/instrumentation/alerts/trackers"
	"github.com/puppetlabs/horsehead/v2/instrumentation/errors"
)

type SNS struct {
	arn   string
	sopts session.Options
}

func (s SNS) NewCapturer() trackers.Capturer {
	return &Capturer{
		arn:   s.arn,
		sopts: s.sopts,
	}
}

type Builder struct {
	arn   string
	sopts session.Options
}

func (b *Builder) WithEnvironment(environment string) *Builder {
	return b
}

func (b *Builder) WithRelease(release string) *Builder {
	return b
}

func (b *Builder) Build() *SNS {
	return &SNS{
		arn:   b.arn,
		sopts: b.sopts,
	}
}

func NewBuilder(arn string, sopts session.Options) (*Builder, errors.Error) {
	b := &Builder{
		arn:   arn,
		sopts: sopts,
	}
	return b, nil
}
