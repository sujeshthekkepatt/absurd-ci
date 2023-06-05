## Design


1. Create the AbsurdCI CR
2. The controller will init the statuses 
3. Controller will go through the steps
4. Controller will init a PVC

5. Every step is ran inside a POD where every commands ran inside an init container to preserve the order

    a. There will be additional init container which will prepare the working directory. It will clone the git repo. This will be the working directory of the containers.

    b. Then we will create additional init containers which will have the commands specified in the steps

    c. There will be a common container which may collect the logs and update the CR Statuses. This may be  a new pod also which will collect everything. In that case the last pod will update the CR

    d. Based on the CR updations the controller will get triggered again and will launch the next step in a new pod. 

    e. This process happens till all the steps are executed


6. We will use PVC to share data between multiple steps. Unless specified each steps will have same working directory. Make sure to use a storge type which support "ReadWriteMany" option. Else we need to delete pods that completed execution.

7. The pods may get deleted or maynot in the POC

8. The logs of execution can be aggragted to a file in the PVC and end of the steps we can use a pod to collect the evidence.





## CR status


1. PodExecutionContext{}
    * TotalNumberOfSteps
    * TotalNumberOfTasks
    * CurrentStep
    * NextStep
    * TotalNumberOfTasksCompleted
    * TotalNUmberOfStepsCompleted
    * TotalExecutionTime
    * Namespace
2. StepPodInfo
    * StepName
    * PodName
    * IsPodCreated
    * ConatinerNames []{}
        * ContainerName
        * CommandStatus
        * CommandLog
