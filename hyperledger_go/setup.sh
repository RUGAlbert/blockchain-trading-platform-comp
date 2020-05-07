#!/bin/bash

if [[ $# -lt 1 ]] ; then
  printHelp
  exit 0
else
  MODE=$1
  shift
fi

function setWorkDir(){
	cd test-network
}

function removeWorkDir(){
	cd ..
}

function start(){
	setWorkDir
	#first setup network

	./network.sh up
	
	#start up chaincode
	./network.sh createChannel
	#./network.sh deployCC
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
	echo "      - 'deploy' - deploy contracts"
}

function deploy(){
	setWorkDir
	peer lifecycle chaincode queryinstalled -O json >&log.txt
	PACKAGE_VERSION=$(($(echo $(grep -Po '"version": ".*?[^\\]"' log.txt) | tr -dc '0-9')+1))
	./network.sh deployCC -l golang -v $PACKAGE_VERSION
	removeWorkDir
}

if [ "${MODE}" == "start" ]; then
  start
elif [ "${MODE}" == "stop" ]; then
  stop
elif [ "${MODE}" == "restart" ]; then
  restart
elif [ "${MODE}" == "deploy" ]; then
	deploy
else
  printHelp
  exit 1
fi
