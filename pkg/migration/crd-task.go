package migration

import (
	"context"
	"fmt"
	"github.com/fyuan1316/asm-operator/pkg/oprlib/manage"
	"io/ioutil"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

type MigCrdTask struct {
}

var ChangeCrdTask MigCrdTask

func init() {
	ChangeCrdTask = MigCrdTask{}

}

func (m MigCrdTask) Run(manage *manage.OperatorManage) error {

	fmt.Println("ChangeCrdTask Run")
	crdList := &apiextensionsv1.CustomResourceDefinitionList{}
	err := manage.K8sClient.List(context.Background(), crdList)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile("pkg/migration/crds/1-vs.yaml")
	if err != nil {
		return err
	}
	crd := apiextensionsv1.CustomResourceDefinition{}
	err = yaml.Unmarshal(bytes, &crd)
	if err != nil {
		return err
	}
	err = manage.K8sClient.Create(context.Background(), &crd)
	if errors.IsAlreadyExists(err) {
		err = nil
	}
	return err
}

var _ manage.ExecuteItem = MigCrdTask{}

func (m MigCrdTask) PreRun(client client.Client) error {
	fmt.Println("ChangeCrdTask prerun")
	crdList := &apiextensionsv1.CustomResourceDefinitionList{}
	err := client.List(context.Background(), crdList)
	fmt.Println(len(crdList.Items), err)
	return nil
}

func (m MigCrdTask) PostRun(client client.Client) error {
	fmt.Println("ChangeCrdTask PostRun")
	return nil
}

func (m MigCrdTask) PreCheck(client client.Client) (bool, error) {
	fmt.Println("ChangeCrdTask PreCheck")
	return true, nil
}

func (m MigCrdTask) PostCheck(client client.Client) (bool, error) {
	fmt.Println("ChangeCrdTask PostCheck")
	return true, nil
}
