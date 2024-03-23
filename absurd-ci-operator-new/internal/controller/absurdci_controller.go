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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

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
	l := log.FromContext(ctx)

	l.Info("From controller", "req", req)

	absurdCIconfig := &batchv1.AbsurdCI{}

	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, absurdCIconfig)

	// todo handle errors properly using kubernetes errors package
	if err != nil {

		fmt.Println("Error occured", err)
		return ctrl.Result{}, nil
	}

	fmt.Println(absurdCIconfig.ObjectMeta.DeletionTimestamp)

	if absurdCIconfig.ObjectMeta.DeletionTimestamp != nil {

		fmt.Println("Object is deleted")

		return ctrl.Result{}, nil
	}

	if !absurdCIconfig.Status.IsPipelineStarted {

		absurdCIconfig, _ = InitSpecAndStatus(absurdCIconfig)

		err := InitPVC(r, ctx, req, absurdCIconfig)

		if err != nil {

			fmt.Println("Error", err)
			return ctrl.Result{}, nil

		}

		err = r.Status().Update(ctx, absurdCIconfig)

		if err != nil {

			fmt.Println(err)
			return ctrl.Result{}, nil

		}

	} else {

		fmt.Println("Pod inited else")

		if absurdCIconfig.Status.APodExecutionContextInfo.TotalNumberOfSteps == absurdCIconfig.Status.APodExecutionContextInfo.TotalNumberOfStepsCompleted {

			fmt.Println("We ran all the steps. Please proceed to clean up")
			return ctrl.Result{}, nil

		}

		fmt.Println("Going to init pod creation")

		needUpdate := CreateStepPodCreationInfo(absurdCIconfig.Status.APodExecutionContextInfo.CurrentStep, absurdCIconfig)

		fmt.Println("statuses", absurdCIconfig.Status.APodExecutionContextInfo.Steps)

		// currently everything will be running on a single pod as init container
		if needUpdate {

			CreateWorkerPod(r, ctx, req, absurdCIconfig)
			err := r.Status().Update(ctx, absurdCIconfig)

			if err != nil {

				fmt.Println(err)
			}
			fmt.Println("Test Pod Created")
		}

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AbsurdCIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.AbsurdCI{}).
		Complete(r)
}

// func (r *AbsurdCIReconciler) CreateWorkerPod(ctx context.Context, req ctrl.Request, ciConfig *batchv1.AbsurdCI, stepCommands []batchv1.ACommand) error {

// 	/*

// 	* todo:
// 	* run client-gen to generate client
// 	* poll the status of current task runner
// 	* launch log collector on the same namespace (Instead of polling I can run a long loop to check the statuses)
// 	* The log collector update the Status field (CommandsRan) using kubectl command.
// 	* retrieve logs
// 	* clone the filesystem
// 	* launch new taskrunner pod with previous file system as working dir

// 	 */
// 	initWorkingDir := corev1.Container{
// 		Name:            "init-working-dir",
// 		Image:           "sujeshthekkepatt/absurd-ci-node-executor:v1.0.0",
// 		Command:         []string{"git", "clone", "https://github.com/sujeshthekkepatt/absurd-ci.git", "workspace/app"},
// 		ImagePullPolicy: corev1.PullAlways,
// 		VolumeMounts: []corev1.VolumeMount{
// 			{
// 				MountPath: "/workspace/app",
// 				Name:      "working-dir",
// 			},
// 		},
// 		Env: []corev1.EnvVar{
// 			{
// 				Name:  "GIT_SSH_COMMAND",
// 				Value: "ssh -o StrictHostKeyChecking=no",
// 			},
// 		},
// 	}
// 	initContainers := []corev1.Container{}
// 	initContainers = append(initContainers, initWorkingDir)

// 	for _, sCommand := range stepCommands {

// 		container := corev1.Container{
// 			Name:            fmt.Sprintf("%s-%s", sCommand.Name, ciConfig.Name),
// 			Image:           "sujeshthekkepatt/absurd-ci-node-executor:v1.0.0",
// 			Command:         []string{sCommand.Command},
// 			Args:            sCommand.Args,
// 			ImagePullPolicy: corev1.PullAlways,
// 			WorkingDir:      "/workspace/app",
// 			VolumeMounts: []corev1.VolumeMount{
// 				{
// 					MountPath: "/workspace/app",
// 					Name:      "working-dir",
// 				},
// 			},
// 			Env: []corev1.EnvVar{
// 				{
// 					Name:  "GIT_SSH_COMMAND",
// 					Value: "ssh -o StrictHostKeyChecking=no",
// 				},
// 			},
// 		}
// 		initContainers = append(initContainers, container)
// 	}

// 	pod := &corev1.Pod{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      ciConfig.Spec.Name,
// 			Namespace: ciConfig.Namespace,
// 		},
// 		Spec: corev1.PodSpec{
// 			Volumes: []corev1.Volume{{
// 				Name:         "working-dir",
// 				VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
// 			}},
// 			InitContainers: initContainers,

// 			Containers: []corev1.Container{{
// 				Name:            fmt.Sprintf("consolidate-log-%s", ciConfig.Spec.Name),
// 				Image:           "sujeshthekkepatt/absurd-ci-node-executor:v1.0.0",
// 				Command:         []string{"echo", "tasks are all ran"},
// 				ImagePullPolicy: corev1.PullAlways,
// 				WorkingDir:      "/workspace/app",
// 				VolumeMounts: []corev1.VolumeMount{
// 					{
// 						MountPath: "/workspace/app",
// 						Name:      "working-dir",
// 					},
// 				},
// 			}},
// 			RestartPolicy: corev1.RestartPolicyOnFailure,
// 		},
// 	}
// 	err := r.Create(ctx, pod)

// 	if err != nil {

// 		fmt.Println(err)
// 	}
// 	return nil
// }

// func (r *AbsurdCIReconciler) RunCommandInPod(ctx context.Context, req ctrl.Request, ciConfig *batchv1.AbsurdCI) error {

// 	pod := &corev1.Pod{}

// 	r.Client.Get(ctx, types.NamespacedName{Name: ciConfig.Spec.Name, Namespace: ciConfig.Namespace}, pod)

// 	if len(ciConfig.Status.CommandsRan) == 0 {

// 		ciConfig.Status.CommandsRan[ciConfig.Spec.Steps[0].Name] = true
// 	} else {

// 		for name, command := range ciConfig.Spec.Steps {

// 			if _, ok := ciConfig.Status.CommandsRan[name]; !ok {

// 			}
// 		}

// 	}
// 	pod.Spec.Containers[0].Command =
// 		fmt.Println("pod is receieved", pod)

// 	return nil
// }

// func (r *AbsurdCIReconciler) CreateWorkerJob(ctx context.Context, req ctrl.Request, ciConfig *batchv1.AbsurdCI, stepCommands map[string]string) error {

// 	job := &kbatchv1.Job{

// 		ObjectMeta: v1.ObjectMeta{
// 			Name: "my-job",
// 		},
// 		Spec: kbatchv1.JobSpec{

// 			Template: corev1.PodTemplateSpec{
// 				Spec: corev1.PodSpec{
// 					RestartPolicy: core,
// 				},
// 			},
// 		},
// 	}

// }
