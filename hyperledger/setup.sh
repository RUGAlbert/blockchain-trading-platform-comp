#!/bin/bash

if [[ $# -lt 1 ]] ; then
  printHelp
  exit 0
else
  MODE=$1
  shift
fi

function setWorkDir(){
	cd Network
}

function removeWorkDir(){
	cd ..
}

function start(){
	setWorkDir
	//first setup network

	./network.sh up
	removeWorkDir
}

function stop(){
	setWorkDir
	./network.sh down
	removeWorkDir
}

function restart(){
	stop
	start
}


function printHelp(){
	echo "Usage: "
	echo "  ./setup.sh <Mode>"
	echo "    <Mode>"
	echo "      - 'start' - bring up fabric orderer and peer nodes and start nodes"
	echo "      - 'stop' - clear the network with docker-compose down"
	echo "      - 'restart' - restart the docker"
}

if [ "${MODE}" == "start" ]; then
  start
elif [ "${MODE}" == "stop" ]; then
  stop
elif [ "${MODE}" == "restart" ]; then
  stop
else
  printHelp
  exit 1
fi
