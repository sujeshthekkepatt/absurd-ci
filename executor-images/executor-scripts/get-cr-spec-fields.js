const util = require('util');
const exec = util.promisify(require('child_process').exec);
async function getCrSpecFields(crName) {

const command = `kubectl get AbsurdCI ${crName} --output=json`;



const {stdout,stderr} = await exec(command)

if (stderr) {

    console.log("Error",stderr)
    return null
}


const spec = JSON.parse(stdout)



spec.status.apodExecutionContextInfo.totalNumberOfStepsCompleted = spec.status.apodExecutionContextInfo.totalNumberOfStepsCompleted + 1 

return spec
}


async function getLogOfContainer(podName,containerName,nameSpace) {

    const command = `kubectl logs ${podName} -c ${containerName} -n ${nameSpace}`;
    
    
    
    const {stdout,stderr} = await exec(command)
    
    if (stderr) {
    
        console.log("Error",stderr)
        return null
    }
    
    
    return JSON.stringify(stdout)
    }



// getCrSpecFields("absurdci-sujesh").catch((err)=>console.log(err))

// getLogOfContainer("ci-pipeline-test-1","init-working-dir","default").catch((err)=>console.log(err))


module.exports = {
    getCrSpecFields,
    getLogOfContainer
    
}