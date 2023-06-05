#!/usr/bin/env node
const { exec } = require('child_process');
const { exit } = require('process');
const { json } = require('stream/consumers');

const utilc = require("./get-cr-spec-fields")
const args = process.argv.slice(2);

const nameSpace = args[0]

const stepName = args[1]

const baseImage = args[2]

const crName = args[3]



// Define the kubectl command and arguments
const command = `kubectl get pod ${baseImage} -n ${nameSpace} --output=json`;

const commandsRan = []

// Execute the kubectl command
exec(command, async(error, stdout, stderr) => {
  if (error) {
    console.error('Error executing kubectl command:', error);
    return;
  }

  if (stderr) {
    console.log(stderr)

    //todo fetch log
    exit(1)
  } else{


  const podInformation = JSON.parse(stdout)

while(true){

  console.log(podInformation.status.phase)
   if (podInformation.status.phase !== "Pending" || podInformation.status.phase !== "Running"){

  console.log(podInformation.status.containerStatuses)

const allContainerStatuses = []

// podInformation.status.containerStatuses.forEach((container)=>{

// allContainerStatuses.push(container)
// })

const commandRanObject = {
    containerName:"",
    commandStatus: "",
    containerStatus: "",
    commandLog: ""
}

containerNames=podInformation.status.initContainerStatuses.map(async(container)=>{
      let log = await utilc.getLogOfContainer(baseImage,container.name,nameSpace)
       return {
          containerName:container.name,
          containerStatus:JSON.stringify(container.state),
          commandStatus:"",
          commandLog:log
        }
    })

   containers= await Promise.allSettled(containerNames)

   containers = containers.map((container)=>{delete container.status; return container.value})

   console.log(containers)
    crSpec = await utilc.getCrSpecFields(crName)

    crSpec.status.astepPodCreationInfo[crSpec.status.apodExecutionContextInfo["currentStepName"]].containerNames = containers

    crSpec.status.astepPodCreationInfo[crSpec.status.apodExecutionContextInfo["currentStepName"]].podStatus = podInformation.status.phase

    statsus = JSON.stringify(crSpec.status)
    console.log(statsus)
  const patchCommand= `kubectl patch AbsurdCI ${crName} --subresource='status' -p '{"status": ${statsus}}' --type='merge'`

  exec(patchCommand, (error, stdout, stderr) => {
    if (error) {
      // console.error('Error executing kubectl command:', error);
      return;
    }
    // console.log(stdout,stderr)
  })
  break;
} else {
    console.log("Running")
}
  }

  }
});