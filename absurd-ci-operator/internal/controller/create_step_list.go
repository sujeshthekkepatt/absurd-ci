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

	err := OrderSteps(cr)

	//fmt.Println("DAG", cr.Status.Dag)

	if err != nil {
		//fmt.Println("error occured while performing ordering")

		return nil, err
	}

	cr.Status.APodExecutionContextInfo, _ = createStepsList(cr)
	cr.Status.AStepPodCreationInfo = make(map[string]batchv1.AStepPodInfo)

	//fmt.Println("astep pod creation info", cr.Status.APodExecutionContextInfo)

	return cr, nil
}

func createStepsList(crSpec *batchv1.AbsurdCI) (batchv1.APodExecutionContext, error) {

	stepList := []batchv1.AStep{}

	var podContext batchv1.APodExecutionContext

	podContext.Steps = stepList

	totalSteps := 0
	totalStepCommands := 0
	var initStep batchv1.AStep
	for _, step := range crSpec.Status.Dag {
		totalSteps += 1
		for range step.Commands {
			totalStepCommands += 1
		}

		if len(step.Environments.Envs) <= 0 && step.Environments.SecretName == "" {

			//fmt.Println("no env")

			step.Environments = batchv1.AStepEnv{

				SecretName:    "",
				ConfigMapName: "",
				Envs:          []batchv1.AEnv{},
				MountOptions: batchv1.AMountOptions{
					VolumeName:    "",
					MountToEnv:    false,
					MountToVolume: false,
					MappingConfig: []batchv1.AMappingConfig{},
				},
			}

		}

		if step.Order == 0 {
			initStep = step
		}

		// fmt.Println("envvars", envVars)
		// fmt.Println("envFromSource", envFromSource)
		// fmt.Println("envFromVolume", envFromVolume)
		// fmt.Println("volumeMounts", volumeMounts)

		// volumeMounts = append(volumeMounts, corev1.VolumeMount{MountPath: "/workspace/app",
		// Name: "working-dir",
		// })

		//fmt.Println("steps are", step.Environments)

		stepList = append(stepList, step)
	}

	//fmt.Println("current step", initStep)

	podContext.CurrentStep = initStep
	podContext.Steps = stepList
	podContext.TotalExecutionTime = "nil"
	podContext.TotalNumberOfSteps = totalSteps
	podContext.TotalNumberOfStepsCompleted = 0

	return podContext, nil
}

func getNextItem(currentStep batchv1.AStep, crSteps []batchv1.AStep) batchv1.AStep {

	var currentStepPosition int
	var nextStep batchv1.AStep
	for position, element := range crSteps {

		if element.Order == currentStep.Order+1 {
			currentStepPosition = position
			nextStep = element
			break
		}
	}

	if currentStepPosition == len(crSteps) {

		return batchv1.AStep{}
	}
	return nextStep

}

func getFirstStep(crSteps []batchv1.AStep) batchv1.AStep {

	var nextStep batchv1.AStep
	for _, element := range crSteps {

		if element.Order == 0 {
			nextStep = element
			break
		}
	}

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

func getPodStatus(r *AbsurdCIReconciler, ctx context.Context, cr *batchv1.AbsurdCI, podName string) (corev1.PodStatus, error) {

	pod := &corev1.Pod{}
	err := r.Get(ctx, types.NamespacedName{Namespace: cr.Namespace, Name: podName}, pod)

	if err != nil {
		fmt.Println("error getting pod sttaus", err)
		return pod.Status, err
	}

	return pod.Status, nil
}

// currentStep is actually nextStep
func CreateStepPodCreationInfo(r *AbsurdCIReconciler, ctx context.Context, currentStep batchv1.AStep, cr *batchv1.AbsurdCI) (bool, bool) {

	podInfo, exists := cr.Status.AStepPodCreationInfo[currentStep.Name]

	//fmt.Println("getting pod info", podInfo, exists)

	if exists {

		// fmt.Println("from pod creation info", podInfo, exists, cr.Status.AStepPodCreationInfo, currentStep)

		//fmt.Println("getting pod status of", podInfo.PodName)
		status, err := getPodStatus(r, ctx, cr, podInfo.PodName)

		//fmt.Println("status of the pod", status, podInfo.PodName)

		if err != nil {
			if kubeerrors.IsNotFound(err) {

				return true, true
			}

			//fmt.Println("error getting pod status")
			return false, false
		}

		//fmt.Println("pod status ", status.Phase, podInfo.PodName)

		if (status.Phase == "Running") || (status.Phase == "Pending") {

			fmt.Printf("The Step:%s/pod is still running. No need to run the next step. Wait for the update", currentStep.Name)
			return false, false
		} else if status.Phase == "Succeeded" {

			//fmt.Println("Create and schedule new Step/Pod")
			step := getNextItem(currentStep, cr.Status.Dag)
			//fmt.Println("nextstep", step)

			cr.Status.APodExecutionContextInfo.TotalNumberOfStepsCompleted = cr.Status.APodExecutionContextInfo.TotalNumberOfStepsCompleted + 1
			cr.Status.APodExecutionContextInfo.TotalNumberOfTasksCompleted = cr.Status.APodExecutionContextInfo.TotalNumberOfTasksCompleted + len(currentStep.Commands)

			if step.Name != "" {
				cr.Status.AStepPodCreationInfo[step.Name] = batchv1.AStepPodInfo{
					PodName:        fmt.Sprintf("task-pod-%s", step.Name),
					ConatinerNames: []batchv1.AContainerNames{},
					PodStatus:      "Pending",
				}
				cr.Status.APodExecutionContextInfo.CurrentStep = step

				return true, true
			}
			return false, true
		} else {
			//fmt.Println("status is down", status.Phase)
			return false, false
		}
	} else {

		//fmt.Println("Voila we need to start the  step")

		step := getCurrentItem(currentStep.Name, cr.Status.APodExecutionContextInfo.Steps)
		cr.Status.APodExecutionContextInfo.CurrentStep = step

		//fmt.Println("current step  and order is", step.Name, step.Order)
		cr.Status.AStepPodCreationInfo[step.Name] = batchv1.AStepPodInfo{
			PodName:        fmt.Sprintf("task-pod-%s", step.Name),
			ConatinerNames: []batchv1.AContainerNames{},
			PodStatus:      "Pending",
		}
		return true, true
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

	//fmt.Println("from create worker pod")

	var stepCommands []batchv1.ACommand

	currentStep := ciConfig.Status.APodExecutionContextInfo.CurrentStep

	// We don't need a loop here.

	for _, step := range ciConfig.Spec.Steps {

		if step.Name == currentStep.Name {

			stepCommands = append(stepCommands, step.Commands...)
			break
		}

	}

	// this container will throw error if the mount path already exists and if not empty
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

	stepContainers := []corev1.Container{}

	// var envFrom corev1.EnvFromSource
	// var volumeFromEnv corev1.Volume

	var envVars []corev1.EnvVar

	var envFromSource []corev1.EnvFromSource

	var envFromVolume []corev1.Volume

	var volumeMounts []corev1.VolumeMount
	envVars, envFromSource, envFromVolume, volumeMounts = ProcessEnvVarsForThePod(currentStep)

	//fmt.Println("envvars", envVars)
	//fmt.Println("envFromSource", envFromSource)
	//fmt.Println("envFromVolume", envFromVolume)
	// fmt.Println("volumeMounts", volumeMounts)

	volumeMounts = append(volumeMounts, corev1.VolumeMount{MountPath: "/workspace/app",
		Name: "working-dir",
	})

	//each step command ran as a container in the pod. Ideally each step should contain a single command
	for _, sCommand := range stepCommands {
		//fmt.Println("scommand", sCommand.Command)
		var container corev1.Container

		if sCommand.Command == "" {

			container = corev1.Container{
				Name:            fmt.Sprintf("%s-%s", sCommand.Name, ciConfig.Name),
				Image:           currentStep.Executor,
				Args:            sCommand.Args,
				ImagePullPolicy: corev1.PullAlways,
				WorkingDir:      "/workspace/app",
				VolumeMounts:    volumeMounts,
				EnvFrom:         envFromSource,
				Env:             envVars,
			}
		} else {

			container = corev1.Container{
				Name:            fmt.Sprintf("%s-%s", sCommand.Name, ciConfig.Name),
				Image:           currentStep.Executor,
				Command:         []string{sCommand.Command},
				Args:            sCommand.Args,
				ImagePullPolicy: corev1.PullAlways,
				WorkingDir:      "/workspace/app",
				VolumeMounts:    volumeMounts,
				EnvFrom:         envFromSource,
				Env:             envVars,
			}

			//fmt.Println("container", container)
		}

		//fmt.Println("container is", container)
		stepContainers = append(stepContainers, container)
	}

	// this executor container should poll statuses of the step containers
	nodeExecutorContiner := corev1.Container{
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
	}

	stepContainers = append(stepContainers, nodeExecutorContiner)

	envFromVolume = append(envFromVolume, corev1.Volume{
		Name: "working-dir",

		VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
			ClaimName: ciConfig.Status.PVCName,
		}},
	})

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      ciConfig.Status.AStepPodCreationInfo[currentStep.Name].PodName,
			Namespace: ciConfig.Namespace,
		},

		Spec: corev1.PodSpec{
			Volumes:       envFromVolume,
			Containers:    stepContainers,
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	err := r.Create(ctx, pod)

	if err != nil {

		//fmt.Println(err)
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

			//fmt.Println("Error while getting PVC", err)
			return err

		}
	} else {
		log.Println("PVC exists. Skipping PVC creation")

	}

	cr.Status.PVCName = pvcName

	return nil
}

func ProcessEnvVarsForThePod(currentStep batchv1.AStep) ([]corev1.EnvVar, []corev1.EnvFromSource, []corev1.Volume, []corev1.VolumeMount) {

	stepEnv := currentStep.Environments

	var envVars []corev1.EnvVar
	var envFrom []corev1.EnvFromSource
	var volumesFromEnv []corev1.Volume

	var volumeMounts []corev1.VolumeMount

	//fmt.Println("envvars length", len(stepEnv.Envs))
	if len(stepEnv.Envs) > 0 {

		for _, val := range stepEnv.Envs {
			//fmt.Println("envvars ", val)

			envVars = append(envVars, corev1.EnvVar{Name: val.Key, Value: val.Value})
		}

	}

	//fmt.Println("envvars", envVars)

	if stepEnv.SecretName != "" {

		if stepEnv.MountOptions.MountToEnv {

			envFromSecret := corev1.EnvFromSource{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: stepEnv.SecretName,
					},
				},
			}
			envFrom = append(envFrom, envFromSecret)

		}

		if stepEnv.MountOptions.MountToVolume {

			if len(stepEnv.MountOptions.MappingConfig) > 0 {

				for _, val := range stepEnv.MountOptions.MappingConfig {
					volumeFrom := corev1.Volume{
						Name: val.VolumeName, // Name of the volume. This should be accept from AbsurdCI Spec
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: currentStep.SecretName, // Name of the Secret to mount
								Items: []corev1.KeyToPath{
									{
										Key:  val.Key,
										Path: val.Path,
									},
								},
							},
						},
					}

					volumeMount := corev1.VolumeMount{
						MountPath: val.MountPath,
						Name:      val.VolumeName,
					}

					volumeMounts = append(volumeMounts, volumeMount)

					volumesFromEnv = append(volumesFromEnv, volumeFrom)
				}
			} else {

				//map the secret as a whole to a volume

				volumeFrom := corev1.Volume{
					Name: stepEnv.MountOptions.VolumeName, // Name of the volume. This should be accept from AbsurdCI Spec
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: currentStep.SecretName, // Name of the Secret to mount
						},
					},
				}

				volumesFromEnv = append(volumesFromEnv, volumeFrom)
			}
		}

	}
	return envVars, envFrom, volumesFromEnv, volumeMounts
}
