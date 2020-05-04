#!/bin/bash

if [[ $# -lt 1 ]] ; then
  printHelp
  exit 0
else
  org=$1
  shift
fi

function printHelp(){
	echo "Usage: "
	echo "  ./peerManager.sh <org>"
	echo "    <org>"
	echo "      - 'init' - change the initial"
	echo "      - 'org1' - change working node to org1"
	echo "      - 'org2' - change working node to org2"
}

function setupBasicParams(){
	export PATH=${PWD}/bin:${PWD}:$PATH
	export FABRIC_CFG_PATH=$PWD/config/
}

function test(){
	echo $CORE_PEER_TLS_ENABLED
	echo $CORE_PEER_LOCALMSPID=
	echo $CORE_PEER_TLS_ROOTCERT_FILE
	echo $CORE_PEER_MSPCONFIGPATH
	echo $CORE_PEER_ADDRESS
	peer chaincode query -C mychannel -n fabcar -c '{"Args":["queryAllCars"]}'
}


function setupOrg(){
	org=$1
	if [ "${org}" == "init" ]; then
		setupBasicParams
	elif [ "${org}" == "org1" ]; then
		source ./configOrg/configOrg1.sh
	elif [ "${org}" == "org2" ]; then
		source ./configOrg/configOrg2.sh
	elif [ "${org}" == "test" ]; then
		test
	else
		echo "Don't regconize type"
	fi
}

setupOrg $org
