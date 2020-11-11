package notifier

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func Test_Throw(t *testing.T) {
	ctx := context.Background()
	out := logrus.StandardLogger().Out

	// 		WithField("created", node.CreationTimestamp.Time).
	//		WithField("lifetime", time.Now().Sub(node.CreationTimestamp.Time)).

	s := time.Now().Add(-time.Hour*1)
	e := time.Now()
	d := e.Sub(s)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.WithField("d", d.String()).Info("aaa")

	if err := Watch(ctx, out); err != nil {
		t.Fatalf("%+v", err)
	}
}
