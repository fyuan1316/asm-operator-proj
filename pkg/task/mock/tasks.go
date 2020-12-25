package mock

import (
	"context"
	"fmt"
	"github.com/fyuan1316/asm-operator/pkg/migration"
	"github.com/fyuan1316/asm-operator/pkg/oprlib/manage"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetDeployStages() [][]manage.ExecuteItem {
	tasks := [][]manage.ExecuteItem{
		{
			migration.ChangeCrdTask,
			PatchTask{},
		},
		{
			deployTask,
		},
	}

	return tasks
}

type PatchTask struct {
}

func (m PatchTask) Run(manage *manage.OperatorManage) error {
	fmt.Println("PatchTask Run")
	client := manage.K8sClient
	ns := corev1.Namespace{}
	err := client.Get(context.Background(), types.NamespacedName{Name: "default"}, &ns)
	if err != nil {
		return err
	}
	if len(ns.Labels) == 0 {
		ns.Labels = make(map[string]string)
	}
	ns.Labels["asm-opr-patch"] = "test-fy"
	if err := client.Update(context.Background(), &ns); err != nil {
		return err
	}
	return nil
}

var _ manage.ExecuteItem = PatchTask{}

func (m PatchTask) PreRun(client client.Client) error {
	fmt.Println("PatchTask prerun")
	return nil
}

func (m PatchTask) PostRun(client client.Client) error {
	fmt.Println("PatchTask PostRun")
	return nil
}

func (m PatchTask) PreCheck(client client.Client) (bool, error) {
	fmt.Println("PatchTask PreCheck")
	return true, nil
}

func (m PatchTask) PostCheck(client client.Client) (bool, error) {
	fmt.Println("PatchTask PostCheck")
	return true, nil
}
