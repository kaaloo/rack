package helpers

import (
	"os"
	"regexp"

	"github.com/convox/rack/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/convox/rack/Godeps/_workspace/src/github.com/ddollar/logger"
	"github.com/convox/rack/Godeps/_workspace/src/github.com/segmentio/analytics-go"
	"github.com/convox/rack/Godeps/_workspace/src/github.com/stvp/rollbar"
)

var regexpEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
var segment *analytics.Client

func init() {
	rollbar.Token = os.Getenv("ROLLBAR_TOKEN")
	rollbar.Environment = os.Getenv("CLIENT_ID")

	segment = analytics.New(os.Getenv("SEGMENT_WRITE_KEY"))

	if os.Getenv("DEVELOPMENT") == "true" {
		segment.Size = 1
	}

	if regexpEmail.MatchString(os.Getenv("CLIENT_ID")) {
		segment.Identify(&analytics.Identify{
			UserId: os.Getenv("CLIENT_ID"),
			Traits: map[string]interface{}{
				"email": os.Getenv("CLIENT_ID"),
			},
		})
	}
}

func Error(log *logger.Logger, err error) {
	if log != nil {
		log.Error(err)
	}

	if rollbar.Token != "" {
		extraData := map[string]string{
			"AWS_REGION": os.Getenv("AWS_REGION"),
			"RACK":       os.Getenv("RACK"),
			"RELEASE":    os.Getenv("RELEASE"),
			"VPC":        os.Getenv("VPC"),
		}
		extraField := &rollbar.Field{"env", extraData}
		rollbar.Error(rollbar.ERR, err, extraField)
	}
}

func TrackEvent(event string, params map[string]interface{}) {
	log := logrus.WithFields(logrus.Fields{"ns": "api.helpers", "at": "TrackEvent"})

	params["client_id"] = os.Getenv("CLIENT_ID")
	params["stack_id"] = os.Getenv("STACK_ID")

	userId := os.Getenv("CLIENT_ID")

	if stackId := os.Getenv("STACK_ID"); stackId != "" {
		userId = stackId
	}

	log.WithFields(logrus.Fields(params)).Info()

	segment.Track(&analytics.Track{
		Event:      event,
		UserId:     userId,
		Properties: params,
	})
}
