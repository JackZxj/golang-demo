package main

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/*
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  # name 必须匹配后面 spec 中的字段，且使用格式 <plural>.<group>
  name: crontabs.example.com
spec:
  # 组名，用于 REST API: /apis/<group>/<version>
  group: example.com
  # 此 CustomResourceDefinition 所支持的版本列表
  versions:
  - name: v1beta1
    # 每个 version 可以通过 served 标志启用或禁止
    served: true
    # 有且只能有一个 version 必须被标记为存储版本
    storage: true
    # schema 是必需字段
    schema:
      openAPIV3Schema:
        type: object
        properties: # crd 定义可以没有apiVersion,kind,metadata，apiserver会自动生成
          host:
            type: string
          port:
            type: string
  # conversion 节是 Kubernetes 1.13+ 版本引入的，其默认值为无转换，即 strategy 子字段设置为 None。
  conversion:
    # None 转换假定所有版本采用相同的模式定义，仅仅将定制资源的 apiVersion 设置为合适的值.
    strategy: None
  # 可以是 Namespaced 或 Cluster
  scope: Namespaced
  names:
    # 名称的复数形式，用于 URL: /apis/<group>/<version>/<plural>
    plural: crontabs
    # 名称的单数形式，用于在命令行接口和显示时作为其别名
    singular: crontab
    # kind 通常是驼峰编码（CamelCased）的单数形式，用于资源清单中
    kind: CronTab
    listKind: AnyThings  # 名称可以随意，但请求的时候还是得用CronTabList
    # shortNames 允许你在命令行接口中使用更短的字符串来匹配你的资源
    shortNames:
    - ct
---
# 创建的时候还是需要带上 gvk + name（namespace）
apiVersion: "example.com/v1"
kind: CronTab
metadata:
  name: my-new-cron-object
host: asd
*/

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/members.config")
	if err != nil {
		panic(err)
	}
	cli, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}

	gvk := schema.GroupVersionKind{
		Group:   "example.com",
		Version: "v1beta1",
		Kind:    "CronTab" + "List",
		//Kind: "AnyThings", // 不支持
	}

	ul := &unstructured.UnstructuredList{}
	ul.SetGroupVersionKind(gvk)
	err = cli.List(context.TODO(), ul)
	if err != nil {
		panic(err)
	}
	fmt.Println(ul)
	// kind 还是 AnyThings
	// &{map[apiVersion:example.com/v1beta1 kind:AnyThings metadata:map[continue: resourceVersion:112889]] [{map[apiVersion:example.com/v1beta1 host:a1 kind:CronTab metadata:map[annotations:map[kubectl.kubernetes.io/last-applied-configuration:{"apiVersion":"example.com/v1beta1","host":"a1","kind":"CronTab","metadata":{"annotations":{},"name":"my-new-cron-object","namespace":"default"}}] creationTimestamp:2023-05-24T11:43:30Z generation:1 managedFields:[map[apiVersion:example.com/v1beta1 fieldsType:FieldsV1 fieldsV1:map[f:host:map[] f:metadata:map[f:annotations:map[.:map[] f:kubectl.kubernetes.io/last-applied-configuration:map[]]]] manager:kubectl-client-side-apply operation:Update time:2023-05-24T11:43:30Z]] name:my-new-cron-object namespace:default resourceVersion:110503 uid:36c4ecba-d448-479f-aa49-1636b0cfd9ee]]}]}
}
