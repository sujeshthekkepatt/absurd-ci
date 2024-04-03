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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ACommand struct {
	Name string `json:"name"`
	// +kubebuilder:validation:Optional
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type AEnv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AMappingConfig struct {
	Key  string `json:"key"`
	Path string `json:"path"`
}

type AMountOptions struct {
	MountToEnv    bool             `json:"mountToEnv"`
	MountToVolume bool             `json:"mountToVolume"`
	MappingConfig []AMappingConfig `json:"mappingConfig"`
}

type AStepEnv struct {
	SecretName    string        `json:"secretName"`
	ConfigMapName string        `json:"configMapName"`
	Envs          []AEnv        `json:"envs"`
	MountOptions  AMountOptions `json:"mountOptions"`
}

type AStep struct {
	Name     string     `json:"name"`
	Executor string     `json:"executor"`
	Commands []ACommand `json:"commands"`
	RunAfter string     `json:"runAfter"`
	// +kubebuilder:validation:Optional
	Order int `json:"order"`
	// +kubebuilder:validation:Optional
	SecretName string `json:"secretName"`
	// +kubebuilder:validation:Optional
	Environments AStepEnv `json:"stepEnvs"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AbsurdCISpec defines the desired state of AbsurdCI
type AbsurdCISpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AbsurdCI. Edit absurdci_types.go to remove/update
	Name    string  `json:"name,omitempty"`
	Version string  `json:"version,omitempty"`
	Steps   []AStep `json:"steps,omitempty"`
}

type ACommandStatus struct {
	CommandName   string `json:"commandName"`
	CommandStatus string `json:"commandStatus"`
}

type ACommandRan struct {
	StepName      string           `json:"taskName"`
	CommandStatus []ACommandStatus `json:"commandStatus"`
	IsFail        bool             `json:"isFail"`
}

type APodExecutionContext struct {
	CurrentStep                 AStep   `json:"currentStep"`
	Steps                       []AStep `json:"steps"`
	TotalNumberOfTasks          int     `json:"totalNumberOfTasks"`
	TotalNumberOfSteps          int     `json:"totalNumberOfSteps"`
	TotalNumberOfTasksCompleted int     `json:"totalNUmberOfTasksCompleted"`
	TotalNumberOfStepsCompleted int     `json:"totalNumberOfStepsCompleted"`
	TotalExecutionTime          string  `json:"totalExecutionTime"`
}

type AContainerNames struct {
	ContainerName   string `json:"containerName"`
	CommandStatus   string `json:"commandStatus"`
	ContainerStatus string `json:"containerStatus"`
	CommandLogs     string `json:"commandLog"`
}

type AStepPodInfo struct {
	PodName        string            `json:"podname"`
	ConatinerNames []AContainerNames `json:"containerNames"`
	PodStatus      string            `json:"podStatus"`
}

// AbsurdCIStatus defines the observed state of AbsurdCI
type AbsurdCIStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	IsPipelineStarted        bool                    `json:"isPipelineStarted"`
	APodExecutionContextInfo APodExecutionContext    `json:"apodExecutionContextInfo"`
	AStepPodCreationInfo     map[string]AStepPodInfo `json:"astepPodCreationInfo"`
	CRName                   string                  `json:"crName"`
	Namespace                string                  `json:"namespace"`
	PVCName                  string                  `json:"pvcName"`
	Dag                      []AStep                 `json:"dag"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AbsurdCI is the Schema for the absurdcis API
type AbsurdCI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AbsurdCISpec   `json:"spec,omitempty"`
	Status AbsurdCIStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AbsurdCIList contains a list of AbsurdCI
type AbsurdCIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AbsurdCI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AbsurdCI{}, &AbsurdCIList{})
}
