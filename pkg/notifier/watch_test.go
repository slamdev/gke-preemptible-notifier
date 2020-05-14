package notifier

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
)

func Test_Throw(t *testing.T) {
	ctx := context.Background()
	out := logrus.StandardLogger().Out

	if err := Watch(ctx, out); err != nil {
		t.Fatalf("%+v", err)
	}
}
