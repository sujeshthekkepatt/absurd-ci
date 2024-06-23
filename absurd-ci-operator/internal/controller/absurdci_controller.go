/*
Copyright 2023.

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

package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1 "github.com/sujeshthekkepatt/absurd-ci/api/v1"
)

// AbsurdCIReconciler reconciles a AbsurdCI object
type AbsurdCIReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch.api.absurd-ci.xyz,resources=absurdcis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.api.absurd-ci.xyz,resources=absurdcis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.api.absurd-ci.xyz,resources=absurdcis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AbsurdCI object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *AbsurdCIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// l := log.FromContext(ctx)

	// l.Info("From controller")

	absurdCIconfig := &batchv1.AbsurdCI{}

	fmt.Println("TS:", time.Now(), req.Name, req.Namespace)
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, absurdCIconfig)

	// todo handle errors properly using kubernetes errors package
	if err != nil {

		fmt.Println("Error occured", err)
		return ctrl.Result{}, nil
	}

	//fmt.Println(absurdCIconfig.ObjectMeta.DeletionTimestamp)

	if absurdCIconfig.ObjectMeta.DeletionTimestamp != nil {

		fmt.Println("Object is deleted")

		return ctrl.Result{}, nil
	}

	if !absurdCIconfig.Status.IsPipelineStarted {

		absurdCIconfig, _ = InitSpecAndStatus(absurdCIconfig)

		err := InitPVC(r, ctx, req, absurdCIconfig)

		if err != nil {

			fmt.Println("Error", err)
			return ctrl.Result{}, err

		}

		err = r.Status().Update(ctx, absurdCIconfig)

		if err != nil {

			fmt.Println("error while updating status", err)
			return ctrl.Result{}, err

		}

	} else {

		fmt.Println("Pod inited else")

		if absurdCIconfig.Status.APodExecutionContextInfo.TotalNumberOfSteps == absurdCIconfig.Status.APodExecutionContextInfo.TotalNumberOfStepsCompleted {

			fmt.Println("We ran all the steps. Please proceed to clean up")
			return ctrl.Result{}, nil

		}

		fmt.Println("Going to init pod creation")

		// for {

		needUpdate, needCrUpdate := CreateStepPodCreationInfo(r, ctx, absurdCIconfig.Status.APodExecutionContextInfo.CurrentStep, absurdCIconfig)

		//fmt.Println("need update is", needUpdate)
		// currently everything will be running on a single pod as init container
		if needUpdate {

			CreateWorkerPod(r, ctx, req, absurdCIconfig)
			err := r.Status().Update(ctx, absurdCIconfig)

			if err != nil {

				fmt.Println(err)
			}
			//fmt.Println("Test Pod Created")
		} else {
			if needCrUpdate {
				err := r.Status().Update(ctx, absurdCIconfig)

				if err != nil {

					fmt.Println(err)
				}
			}
		}
		// }

	}

	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AbsurdCIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.AbsurdCI{}).
		Complete(r)
}
