package manage

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ExecuteItem interface {
	PreRun(client.Client) error
	PostRun(client.Client) error
	PreCheck(client.Client) (bool, error)
	PostCheck(client.Client) (bool, error)
	Run(*OperatorManage) error
}

type Object interface {
	runtime.Object
	metav1.Object
}
