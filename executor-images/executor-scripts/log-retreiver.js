#!/usr/bin/env node
const { exit } = require('process');
const util = require('util');
const exec = util.promisify(require('child_process').exec);

const args = process.argv.slice(2);

const nameSpace = args[0]

const stepName = args[1]

const baseImage = args[2]

const crName = args[3]


// const podInformation = JSON.parse(stdout)


async function retrieveLogs() {


const logCommand = `kubectl logs ${baseImage} -n ${nameSpace} --all-containers=true`

const {stdout,stderr} = await exec(logCommand)

if (stderr) {

    console.log("Error",stderr)
    return
}

const commandsLogs = JSON.stringify({"stepName":stdout})

console.log(commandsLogs)
const patchCommand= `kubectl patch AbsurdCI ${crName} --subresource='status' -p '{"status":{"commandsLogs":${commandsLogs}}}' --type='merge'`

await exec(patchCommand)



}


retrieveLogs().then(()=>console.log("completed")).catch(err=>console.log(err))