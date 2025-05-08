package main

import (
	"context"
	"encoding/base64"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
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
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-patch",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"bb": {},
		},
	}
	_, err = c.CoreV1().Secrets("default").Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}
	patcher := []map[string]interface{}{{"op": "replace", "path": "/data/aa", "value": base64.StdEncoding.EncodeToString([]byte("aaa"))}}
	patchData, err := json.Marshal(patcher)
	if err != nil {
		panic(err)
	}
	patchDatas := fmt.Sprintf(`[{"op": "replace", "path": "/data/%s", "value": "%s"}]`,
		"secret", base64.StdEncoding.EncodeToString([]byte("abc")))
	got, err := c.CoreV1().Secrets("default").Patch(ctx, "test-patch", types.JSONPatchType, []byte(patchDatas), metav1.PatchOptions{})
	if err != nil {
		fmt.Printf("failed to patch kubeconfig secret into control cluster with data: %s\n", patchData)
		panic(err)
	}
	fmt.Printf("%+v\n", *got)
}
