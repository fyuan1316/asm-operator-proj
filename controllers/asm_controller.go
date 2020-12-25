/*
Copyright 2020 fyuan.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/fyuan1316/asm-operator/pkg/oprlib/manage"
	"github.com/fyuan1316/asm-operator/pkg/task"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sync"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	operatorv1alpha1 "github.com/fyuan1316/asm-operator/api/v1alpha1"
)

// AsmReconciler reconciles a Asm object
type AsmReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var once = sync.Once{}
var mgr *manage.OperatorManage
var tasks [][]manage.ExecuteItem

// +kubebuilder:rbac:groups=operator.alauda.io,resources=asms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.alauda.io,resources=asms/status,verbs=get;update;patch

func (r *AsmReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("asm", req.NamespacedName)

	// your logic here
	instance := &operatorv1alpha1.Asm{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	mgr = manage.NewOperatorManage(r.Client, instance, r.Scheme)
	once.Do(func() {
		tasks = task.GetDeployStages()
	})
	err = mgr.Reconcile(tasks)
	return ctrl.Result{}, nil
}

func (r *AsmReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.Asm{}).
		Complete(r)
}
