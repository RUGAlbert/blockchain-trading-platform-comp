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
}

function installContracts(){
	source ./peerManager.sh init
	source ./peerManager.sh org1
	source ./configOrg/configCa.sh
	peer chaincode install -n evmcc -l golang -v 0 -p ./../../fabric-chaincode-evm/evmcc
	source ./peerManager.sh org2
	peer chaincode install -n evmcc -l golang -v 0 -p ./../../fabric-chaincode-evm/evmcc
	peer chaincode instantiate -n evmcc -v 0 -C mychannel -c '{"Args":[]}' -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA
}

if [ "${MODE}" == "start" ]; then
  start
elif [ "${MODE}" == "stop" ]; then
  stop
elif [ "${MODE}" == "restart" ]; then
  restart
elif [ "${MODE}" == "install" ]; then
	installContracts
else
  printHelp
  exit 1
fi
