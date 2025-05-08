package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	ctx := context.Background()
	config, err := ctrl.GetConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(config.Host)
	c := kubernetes.NewForConfigOrDie(config)

	name := fmt.Sprintf("jack-%d", rand.Intn(100000000))
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "jack",
		},
		Data: map[string]string{},
	}
	cm, err = c.CoreV1().ConfigMaps("jack").Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	watchInterface, err := c.CoreV1().ConfigMaps("jack").Watch(ctx, metav1.ListOptions{FieldSelector: fields.OneTermEqualSelector("metadata.name", name).String()})
	if err != nil {
		panic(err)
	}
	defer watchInterface.Stop()
	go func() {
		cm.Data = map[string]string{"foo": "bar"}
		_, err = c.CoreV1().ConfigMaps("jack").Update(ctx, cm, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println("update cm err:", err)
		}
		fmt.Println("update done")
	}()

	timeout := time.After(10 * time.Second)
LOOP:
	for {
		select {
		case event := <-watchInterface.ResultChan():
			if event.Type == watch.Modified {
				configCR := event.Object.(*corev1.ConfigMap)
				fmt.Println("get cm:", *configCR)
				break LOOP

			} else {
				fmt.Println("get cm event:", event)
				continue
			}
		case <-timeout:
			fmt.Println("watch cluster config timeout")
			break LOOP
		}
	}

	time.Sleep(time.Second)
	// cli := Client{"foo", c}
	// cli.run1()
	// cli.run2()
}

type Client struct {
	name   string
	client kubernetes.Interface
}

func (c Client) run1() {
	n, e := c.client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	fmt.Printf("run1:\n %s err: %v nodes: %v \n\n\n", c.name, e, n)
}
func (c *Client) run2() {
	n, e := c.client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	fmt.Printf("run2:\n %s err: %v nodes: %v \n\n\n", c.name, e, n)
}
