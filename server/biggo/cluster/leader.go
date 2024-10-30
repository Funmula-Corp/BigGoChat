package cluster

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/util/homedir"
)

func (p *BiggoCluster) getKubeConfig() (config *rest.Config, err error) {
	if config, err = rest.InClusterConfig(); err == rest.ErrNotInCluster {
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	}
	return
}

func (p *BiggoCluster) JoinVote(id, lockName, lockNamespace string) (err error) {
	if p.KubeClient, err = kubernetes.NewForConfig(p.KubeConfig); err != nil {
		return
	}

	lock := &resourcelock.LeaseLock{
		Client: p.KubeClient.CoordinationV1(),
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockName,
			Namespace: lockNamespace,
		},
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	var ctx context.Context
	ctx, p.CancelVote = context.WithCancel(context.Background())
	defer p.CancelVote()

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   30 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				p.Leader.Swap(true)
				p.UpdateService()
			},
			OnStoppedLeading: func() {
				p.Leader.Swap(false)
			},
			OnNewLeader: func(identity string) {
			},
		},
	})
	return
}

func (p *BiggoCluster) UpdateService() {
	ctx := context.TODO()

	hostname, _ := os.Hostname()
	namespace := os.Getenv("NAMESPACE")
	serviceLabel := os.Getenv("LEADER_SERVICE_LABEL")
	if namespace == "" || serviceLabel == "" {
		return
	}

	services, err := p.KubeClient.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: serviceLabel,
	})
	if err != nil {
		mlog.Error("cluster.leader.service.list", mlog.Err(err))
		return
	}

	if services != nil {
		for idx := range services.Items {
			for key := range services.Items[idx].Spec.Selector {
				if strings.Contains(services.Items[idx].Spec.Selector[key], "pod-name") {
					services.Items[idx].Spec.Selector[key] = hostname
					if _, err = p.KubeClient.CoreV1().Services(namespace).Update(ctx, &services.Items[idx], metav1.UpdateOptions{}); err != nil {
						mlog.Error("cluster.leader.service.update", mlog.Err(err))
						return
					}
				}
			}
		}
	} else {
		mlog.Error("cluster.leader.service.update", mlog.Err(fmt.Errorf("no service found")))
	}
}
