package controller

import (
	"context"
	"fmt"
	"log"

	batchv1 "github.com/sujeshthekkepatt/absurd-ci/api/v1"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func InitSpecAndStatus(cr *batchv1.AbsurdCI) (*batchv1.AbsurdCI, error) {

	cr.Status.CRName = cr.Name
	cr.Status.IsPipelineStarted = true
	cr.Status.APodExecutionContextInfo, _ = createStepsList(cr.Spec)
	cr.Status.AStepPodCreationInfo = make(map[string]batchv1.AStepPodInfo)

	return cr, nil
}

func createStepsList(crSpec batchv1.AbsurdCISpec) (batchv1.APodExecutionContext, error) {

	stepList := []batchv1.AStep{}

	var podContext batchv1.APodExecutionContext

	podContext.Steps = stepList

	totalSteps := 0
	totalStepCommands := 0
	for _, step := range crSpec.Steps {
		totalSteps += 1
		for range step.Commands {
			totalStepCommands += 1
		}

		stepList = append(stepList, step)
	}
	initStep := stepList[0]

	podContext.CurrentStepName = initStep.Name
	podContext.Steps = stepList
	podContext.TotalExecutionTime = "nil"
	podContext.TotalNumberOfSteps = totalSteps
	podContext.TotalNumberOfStepsCompleted = 0

	return podContext, nil
}

func getNextItem(currentStep string, crSteps []batchv1.AStep) batchv1.AStep {

	var currentStepPosition int
	for position, element := range crSteps {

		if element.Name == currentStep {
			currentStepPosition = position
			break
		}
	}

	if currentStepPosition == len(crSteps)-1 {

		return batchv1.AStep{}
	}
	nextStep := crSteps[currentStepPosition+1]

	return nextStep
}

func getCurrentItem(currentStep string, crSteps []batchv1.AStep) batchv1.AStep {

	var currentStepPosition int
	for position, element := range crSteps {

		if element.Name == currentStep {
			currentStepPosition = position
			break
		}
	}

	nextStep := crSteps[currentStepPosition]

	return nextStep
}

// currentStep is actually nextStep
func CreateStepPodCreationInfo(currentStep string, cr *batchv1.AbsurdCI) bool {

	podInfo, exists := cr.Status.AStepPodCreationInfo[currentStep]

	fmt.Println("from pod creation info", podInfo, exists, cr.Status.AStepPodCreationInfo, currentStep)

	if exists {

		if (podInfo.PodStatus == "Running") || (podInfo.PodStatus == "Pending") {

			fmt.Printf("The Step:%s/pod is still running. No need to run the next step. Wait for the update", currentStep)
			return false
		} else {

			fmt.Println("Create and schedule new Step/Pod")
			step := getNextItem(currentStep, cr.Status.APodExecutionContextInfo.Steps)

			if step.Name != "" {
				cr.Status.AStepPodCreationInfo[step.Name] = batchv1.AStepPodInfo{
					PodName:        fmt.Sprintf("task-pod-%s", step.Name),
					ConatinerNames: []batchv1.AContainerNames{},
					PodStatus:      "Pending",
				}
				cr.Status.APodExecutionContextInfo.CurrentStepName = step.Name
				return true
			}
			return false
		}
	} else {

		fmt.Println("Voila we need to start the  step")

		step := getCurrentItem(currentStep, cr.Status.APodExecutionContextInfo.Steps)
		cr.Status.APodExecutionContextInfo.CurrentStepName = step.Name

		cr.Status.AStepPodCreationInfo[step.Name] = batchv1.AStepPodInfo{
			PodName:        fmt.Sprintf("task-pod-%s", step.Name),
			ConatinerNames: []batchv1.AContainerNames{},
			PodStatus:      "Pending",
		}
		return true
	}

}

func CreateWorkerPod(r *AbsurdCIReconciler, ctx context.Context, req ctrl.Request, ciConfig *batchv1.AbsurdCI) error {

	/*

	* todo:
	* run client-gen to generate client
	* poll the status of current task runner
	* launch log collector on the same namespace (Instead of polling I can run a long loop to check the statuses)
	* The log collector update the Status field (CommandsRan) using kubectl command.
	* retrieve logs
	* clone the filesystem
	* launch new taskrunner pod with previous file system as working dir

	 */

	var stepCommands []batchv1.ACommand

	currentStep := ciConfig.Status.APodExecutionContextInfo.CurrentStepName

	for _, step := range ciConfig.Spec.Steps {

		if step.Name == currentStep {

			stepCommands = append(stepCommands, step.Commands...)
			break
		}

	}

	initWorkingDir := corev1.Container{
		Name:            "init-working-dir",
		Image:           "sujeshthekkepatt/absurd-ci-node-executor:v1.0.0",
		Command:         []string{"git", "clone", "https://github.com/sujeshthekkepatt/absurd-ci.git", "workspace/app"},
		ImagePullPolicy: corev1.PullAlways,
		VolumeMounts: []corev1.VolumeMount{
			{
				MountPath: "/workspace/app",
				Name:      "working-dir",
			},
		},
		Env: []corev1.EnvVar{
			{
				Name:  "GIT_SSH_COMMAND",
				Value: "ssh -o StrictHostKeyChecking=no",
			},
		},
	}
	initContainers := []corev1.Container{}
	initContainers = append(initContainers, initWorkingDir)

	for _, sCommand := range stepCommands {

		container := corev1.Container{
			Name:            fmt.Sprintf("%s-%s", sCommand.Name, ciConfig.Name),
			Image:           "sujeshthekkepatt/absurd-ci-node-executor:v1.0.0",
			Command:         []string{sCommand.Command},
			Args:            sCommand.Args,
			ImagePullPolicy: corev1.PullAlways,
			WorkingDir:      "/workspace/app",
			VolumeMounts: []corev1.VolumeMount{
				{
					MountPath: "/workspace/app",
					Name:      "working-dir",
				},
			},
			Env: []corev1.EnvVar{
				{
					Name:  "GIT_SSH_COMMAND",
					Value: "ssh -o StrictHostKeyChecking=no",
				},
			},
		}
		initContainers = append(initContainers, container)
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      ciConfig.Status.AStepPodCreationInfo[currentStep].PodName,
			Namespace: ciConfig.Namespace,
		},

		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{{
				Name: "working-dir",

				VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: ciConfig.Status.PVCName,
				}},
			}},
			InitContainers: initContainers,

			Containers: []corev1.Container{{
				Name:            fmt.Sprintf("consolidate-log-%s", ciConfig.Spec.Name),
				Image:           "sujeshthekkepatt/absurd-ci-node-executor:v1.0.0",
				Command:         []string{"echo", "tasks are all ran"},
				ImagePullPolicy: corev1.PullAlways,
				WorkingDir:      "/workspace/app",
				VolumeMounts: []corev1.VolumeMount{
					{
						MountPath: "/workspace/app",
						Name:      "working-dir",
					},
				},
			}},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
	err := r.Create(ctx, pod)

	if err != nil {

		fmt.Println(err)
		return err
	}
	return nil
}

func InitPVC(r *AbsurdCIReconciler, ctx context.Context, req ctrl.Request, cr *batchv1.AbsurdCI) error {

	pvcName := fmt.Sprintf("pvc-%s", cr.Name)

	pvc := &corev1.PersistentVolumeClaim{}

	err := r.Get(ctx, types.NamespacedName{Namespace: cr.Namespace, Name: pvcName}, pvc)
	if err != nil {

		if kubeerrors.IsNotFound(err) {

			log.Println("PVC not exists. Creating a PVC")
			pvc := &corev1.PersistentVolumeClaim{
				ObjectMeta: v1.ObjectMeta{
					Name:      pvcName,
					Namespace: cr.Namespace,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources:   corev1.VolumeResourceRequirements{Requests: corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("1Gi")}},
				},
			}

			err := r.Create(ctx, pvc)
			if err != nil {

				log.Println("error while creating PVC", err)
				return err
			}
		} else {

			fmt.Println("Error while getting PVC", err)
			return err

		}
	} else {
		log.Println("PVC exists. Skipping PVC creation")

	}

	cr.Status.PVCName = pvcName

	return nil
}
