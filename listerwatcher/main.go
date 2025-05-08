package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var lw cache.ListerWatcher

func main() {
	cfg, err := buildConfig("/root/.kube/dev-mechine-dcp-test.config")
	if err != nil {
		panic(err)
	}
	client := clientset.NewForConfigOrDie(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	lw = &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Pods(corev1.NamespaceAll).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Pods(corev1.NamespaceAll).Watch(ctx, options)
		},
	}

	rv, err := list(ctx, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println("rv:", rv)

	watcher, err := lw.Watch(metav1.ListOptions{ResourceVersion: rv,
		LabelSelector: "!pod-template-hash"})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 15; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("ctx done")
			return
		case event := <-watcher.ResultChan():
			po := event.Object.(*corev1.Pod)
			switch event.Type {
			case watch.Added:
				fmt.Printf("%d %s   add: %s/%s %s %d\n", i, time.Now(), po.Namespace, po.Name, po.ResourceVersion, po.Generation)
			case watch.Modified:
				fmt.Printf("%d %smodify: %s/%s %s %d\n", i, time.Now(), po.Namespace, po.Name, po.ResourceVersion, po.Generation)
			case watch.Deleted:
				fmt.Printf("%d %sdelete: %s/%s %s %d\n", i, time.Now(), po.Namespace, po.Name, po.ResourceVersion, po.Generation)
			case watch.Error:
				fmt.Printf("watch channel closed unexpectedly: %v\n", event.Object)
				return
			default:
				continue
			}
		}
	}
	cancel()
	time.Sleep(time.Second)
}

func list(ctx context.Context, pageSize int64) (rv string, err error) {
	var resourceVersion string
	listOpts := metav1.ListOptions{
		Limit:         pageSize,
		Continue:      "",
		LabelSelector: "!pod-template-hash",
	}

	for {
		select {
		case <-ctx.Done():
			return "", nil
		default:
			resp, err := lw.List(listOpts)
			if err != nil {
				return "", fmt.Errorf("failed to list resource: %w", err)
			}

			list, err := meta.ListAccessor(resp)
			if err != nil {
				return "", fmt.Errorf("failed to list resource: %w", err)
			}

			listOpts.Continue = list.GetContinue()
			resourceVersion = list.GetResourceVersion()

			// err = meta.EachListItem(resp, func(o runtime.Object) error {
			// 	po := o.(*corev1.Pod)
			// 	fmt.Printf("%s, %s, %s, %d\n", po.Namespace, po.Name, po.ResourceVersion, po.Generation)
			// 	return nil
			// })

			// if err != nil {
			// 	return "", fmt.Errorf("failed to process list item: %w", err)
			// }
		}

		if len(listOpts.Continue) == 0 {
			break
		}
	}

	return resourceVersion, nil
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func handleSignals(cancel context.CancelFunc) {
	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	for sig := range sigc {
		switch sig {
		case syscall.SIGINT:
			cancel()
			fmt.Println("good bye!")
			os.Exit(1)
		case syscall.SIGTERM:
			cancel()
			return
		}
	}
}
