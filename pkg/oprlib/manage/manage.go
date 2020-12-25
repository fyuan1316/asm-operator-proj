package manage

import (
	"context"
	"github.com/fyuan1316/asm-operator/pkg/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

var (
	logger = logging.RegisterScope("controller.oprlib")
)

type OperatorManage struct {
	K8sClient client.Client
	Object    Object
	Scheme    *runtime.Scheme
}

func NewOperatorManage(client client.Client, object Object, scheme *runtime.Scheme) *OperatorManage {
	return &OperatorManage{
		K8sClient: client,
		Object:    object,
		Scheme:    scheme,
	}
}

func (m *OperatorManage) Reconcile(stages [][]ExecuteItem) error {
	//delete
	//if !m.Object.GetDeletionTimestamp().IsZero(){
	//
	//}
	//sync
	return m.ProcessStages(stages)
}

func (m *OperatorManage) ProcessStages(stages [][]ExecuteItem) error {
	for _, items := range stages {
		for _, item := range items {
			//if item.PreCheck != nil {
			logger.Debugf("run precheck")
			if err := loopUntil(context.Background(), 5*time.Second, 10, item.PreCheck, m.K8sClient); err != nil {
				//item.err = err
				return err
			}
			//}
		}
		for _, item := range items {
			//if item.PreRun != nil {
			logger.Debugf("run prerun")
			if err := item.PreRun(m.K8sClient); err != nil {
				//item.err = err
				return err
			}
			//}
		}
		//if err := m.(items); err != nil {
		//	return err
		//}
		for _, item := range items {
			//if item.Run != nil {
			logger.Debugf("execute run")
			if err := item.Run(m); err != nil {
				//item.err = err
				return err
			}
			//}
		}

		for _, item := range items {
			//if item.PostRun != nil {
			logger.Debugf("run postrun")
			if err := item.PostRun(m.K8sClient); err != nil {
				//item.err = err
				return err
			}
			//}
		}
		for _, item := range items {
			//if item.PostCheck != nil {
			logger.Debugf("run postcheck")
			if err := loopUntil(context.Background(), 5*time.Second, 10, item.PostCheck, m.K8sClient); err != nil {
				//item.err = err
				return err
			}
			//}
		}
	}
	return nil
}
