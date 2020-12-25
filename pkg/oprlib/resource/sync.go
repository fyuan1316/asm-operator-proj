package resource

import (
	"fmt"
	"github.com/fyuan1316/asm-operator/pkg/oprlib/manage"
	"io/ioutil"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"
)

type SyncManager struct {
	K8sResource map[string]SyncResource
}

type SyncResource struct {
	manage.Object
	Sync func(client.Client, manage.Object) error
}

func (m *SyncManager) LoadFile(filePath string, res SyncResource) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return m.Load(string(bytes), res)
}
func (m *SyncManager) Load(objectStr string, res SyncResource) error {
	var err error
	object := res.Object
	err = yaml.Unmarshal([]byte(objectStr), object)
	if err != nil {
		return err
	}
	objKey := fmt.Sprintf("%s-%s-%s", object.GetObjectKind().GroupVersionKind().Kind,
		object.GetNamespace(),
		object.GetName(),
	)
	if m.K8sResource == nil {
		m.K8sResource = make(map[string]SyncResource)
	}
	m.K8sResource[objKey] = res
	return err
}

func (m *SyncManager) Sync(om *manage.OperatorManage) error {
	for _, res := range m.K8sResource {
		err := controllerutil.SetControllerReference(om.Object, res.Object, om.Scheme)
		if err != nil {
			return err
		}
		_ = res.Sync(om.K8sClient, res.Object)
	}
	return nil
}
