package notifier

import (
	"cloud.google.com/go/compute/metadata"
	"context"
	"github.com/sirupsen/logrus"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/rest"
	"time"
)

func Watch(ctx context.Context, out io.Writer) error {
	client, err := client()
	if err != nil {
		return err
	}
	wait.Forever(func() {
		err := metadata.Subscribe("instance/preempted", func(state string, exists bool) error {
			if !exists {
				logrus.Error("metadata API deleted unexpectedly")
				return nil
			}
			if state == "TRUE" {
				nodeName, err := metadata.InstanceName()
				if err != nil {
					logrus.Error(err)
					return nil
				}
				if err := notify(ctx, client, nodeName); err != nil {
					logrus.Error(err)
					return nil
				}
			}
			return nil
		})
		if err != nil {
			logrus.Error("failed to get metadata")
		}
	}, 10*time.Second)
	return nil
}

func client() (corev1.CoreV1Interface, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return clientSet.CoreV1(), nil
}

func notify(ctx context.Context, client corev1.CoreV1Interface, nodeName string) error {
	node, err := client.Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	options := metav1.ListOptions{FieldSelector: fields.OneTermEqualSelector("spec.nodeName", nodeName).String()}
	podsList, err := client.Pods(metav1.NamespaceAll).List(ctx, options)
	if err != nil {
		return err
	}

	podNames := make([]string, len(podsList.Items))
	for i := range podsList.Items {
		podNames[i] = podsList.Items[i].Name
	}

	logrus.
		WithField("node", node.Name).
		WithField("created", node.CreationTimestamp).
		WithField("labels", node.Labels).
		WithField("pods", podNames).
		Info("node is going to be preempted soon")
	return nil
}
