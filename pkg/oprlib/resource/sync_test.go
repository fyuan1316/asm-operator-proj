package resource

import (
	"context"
	"fmt"
	"github.com/fyuan1316/asm-operator/pkg/oprlib/manage"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func TestSyncManager_LoadFile(t *testing.T) {
	type fields struct {
		K8sResource SyncResource
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test-load-file",
			fields:  fields{},
			args:    args{filePath: "./pod.yaml"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SyncManager{
				K8sResource: map[string]SyncResource{
					"test1": tt.fields.K8sResource,
				},
			}
			res := SyncResource{
				Object: &appsv1.Deployment{},
				Sync: func(client client.Client, object manage.Object) error {
					deploy := appsv1.Deployment{}
					err := client.Get(context.Background(),
						types.NamespacedName{Namespace: object.GetNamespace(), Name: object.GetName()},
						&deploy,
					)
					if err != nil {
						if errors.IsNotFound(err) {
							errCreate := client.Create(context.Background(), &deploy)
							if errCreate != nil {
								return errCreate
							}
						}
						return err
					} else {
						//update
						wanted := object.(*appsv1.Deployment)
						if !equality.Semantic.DeepDerivative(deploy.Spec, wanted.Spec) {
							deploy.Spec = wanted.Spec
							if errUpd := client.Update(context.Background(), &deploy); errUpd != nil {
								return errUpd
							}
						}
					}
					return nil
				},
			}
			if err := m.LoadFile(tt.args.filePath, res); (err != nil) != tt.wantErr {
				t.Errorf("LoadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println()
		})
	}
}
