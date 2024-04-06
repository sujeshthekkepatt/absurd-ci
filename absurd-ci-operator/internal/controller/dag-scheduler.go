package controller

import (
	"fmt"

	dag "github.com/dominikbraun/graph"
	batchv1 "github.com/sujeshthekkepatt/absurd-ci/api/v1"
)

func StepHash(step string) string {

	return step
}
func OrderSteps(ciConfig *batchv1.AbsurdCI) error {

	absurdDag := dag.New(StepHash, dag.Directed(), dag.Acyclic(), dag.PreventCycles())

	for _, val := range ciConfig.Spec.Steps {
		fmt.Println("from add vertex", val.Name)
		absurdDag.AddVertex(val.Name)
	}

	for _, val := range ciConfig.Spec.Steps {
		if val.RunAfter != "" {

			absurdDag.AddEdge(StepHash(val.RunAfter), StepHash(val.Name))
		}
	}

	tdag, err := dag.TransitiveReduction(absurdDag)
	if err != nil {

		fmt.Println("Error while running transitive reduction sort")
		return err
	}
	orderedSteps, err := dag.TopologicalSort(tdag)
	if err != nil {

		fmt.Println("Error while running topo sort")
		return err
	}
	fmt.Println("dag edges: topo order", orderedSteps)

	orderedAsteps := []batchv1.AStep{}

	for order, val := range orderedSteps {

		step := setOrderInfoOnAstep(order, val, ciConfig)
		orderedAsteps = append(orderedAsteps, step)

	}

	// ciConfig.Spec.Steps = orderedAsteps
	// fmt.Println("spec steps", ciConfig.Spec.Steps)

	ciConfig.Status.Dag = orderedAsteps
	return nil

}

func setOrderInfoOnAstep(order int, stepName string, ciConfig *batchv1.AbsurdCI) batchv1.AStep {

	var stepVal batchv1.AStep
	for _, val := range ciConfig.Spec.Steps {

		if val.Name == stepName {
			if len(val.Environments.Envs) <= 0 && val.Environments.SecretName == "" {

				fmt.Println("no env")

				val.Environments = batchv1.AStepEnv{

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
			val.Order = order
			stepVal = val

			break

		}

	}

	return stepVal
}
