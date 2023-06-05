#!/usr/bin/env node
const { exec } = require('child_process');
const { exit } = require('process');
const { json } = require('stream/consumers');

const args = process.argv.slice(2);

const nameSpace = args[0]

const stepName = args[1]

const baseImage = args[2]

const crName = args[3]



// Define the kubectl command and arguments
const command = `kubectl get pod ${baseImage} -n ${nameSpace} --output=json`;

const commandsRan = []

// Execute the kubectl command
exec(command, (error, stdout, stderr) => {
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

   if (podInformation.status.phase !== "Pending"){

  console.log(podInformation.status.containerStatuses)

const allContainerStatuses = []

// podInformation.status.containerStatuses.forEach((container)=>{

// allContainerStatuses.push(container)
// })

const commandRanObject = {
    taskName: stepName,
    commandStatus: [],
    isFail: false   
}
podInformation.status.initContainerStatuses.forEach((container)=>{

        if (Object.keys(container.lastState).length !=0) {

            if (container.lastState.terminated.reason === "Error") {

                commandRanObject.isFail = true
                commandRanObject.commandStatus.push({commandName:container.name.split("-")[0],commandStatus:JSON.stringify(container.state)})
            }
        } else {
            commandRanObject.commandStatus.push({commandName:container.name.split("-")[0],commandStatus:JSON.stringify(container.state)})

        }


    })


const totalCommandsRan = commandRanObject.commandStatus.length

const totalCommandsObject = [commandRanObject]
const patchCommand= `kubectl patch AbsurdCI ${crName} --subresource='status' -p '{"status":{"totalCommandsRan":${totalCommandsRan},"commandsRan": ${JSON.stringify(totalCommandsObject)} }}' --type='merge'`

exec(patchCommand, (error, stdout, stderr) => {
    if (error) {
      console.error('Error executing kubectl command:', error);
      return;
    }
    console.log(stdout,stderr)
})
break;
} else {
    console.log("Running")
}
  }

  }
});