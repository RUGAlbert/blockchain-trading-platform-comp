#!/bin/bash

function printHelp(){
	echo "Usage: "
	echo "  source ./setup.sh <Mode>"
	echo "    <Mode>"
	echo "      - 'set-env' - set env variables"
	echo "      - 'start' - starts docker containers containing the network"
	echo "      - 'deploy' - deploy contracts"
	echo "      - 'update' - updates contracts"
	echo "      - 'register' - registers contracts"
	echo "      - 'benchmark' - benchmarks it"
	echo "      - 'configUpdate' - updates the config after changing"
	echo "      - 'all' - start, deploy and register, all set and done for the benchmark"
}

function setGatewayDir(){
	cd app/sdk/gateway/
}

function removeGatewayDir(){
	cd ..
	cd ..
	cd ..
}

function set-env(){
	export PATH=$PATH:${PWD}/network/bin
	export GOPATH=${PWD}/gocc
}

function start(){
	source reset-chain-env.sh
	source set-env.sh acme
	source set-ca-msp.sh admin
	
	ca-dev-init.sh
}

function deploy(){
	set-chain-env.sh -n trading -v 1.0  -p token/trading -c '{"Args":["Init"]}' -I false
	chain.sh install -p
	s=$(peer lifecycle chaincode queryinstalled | sed -n '/Package ID: /,/, Label/p' | sed -e 's/Package ID: \(.*\), Label/\1/')
	export INSTALLED_MAX_PACKAGE_ID="${s%:*}"
	chain.sh instantiate
}

function update(){
	. cc.env.sh
	ver=$(($CC2_SEQUENCE+1))
	set-chain-env.sh -s $ver
	deploy
}

function register(){
	#register
	set-chain-env.sh  -i '{"Args": ["register"]}'
	
	. set-ca-msp.sh  mary
	chain.sh invoke
	
	#set the leader, since only one is registerd on chain, only one can be the leader
	set-chain-env.sh  -i '{"Args": ["setNextLeader"]}'
	chain.sh invoke
	
	#register also other users
	set-chain-env.sh  -i '{"Args": ["register"]}'

	. set-ca-msp.sh  john
	chain.sh invoke

	. set-ca-msp.sh  anil
	chain.sh invoke

	. set-ca-msp.sh  admin
	chain.sh invoke
	
	#give some tickets to some of the users
	set-chain-env.sh  -i '{"Args": ["issueTicket", "eDUwOTo6Q049YW5pbCxPVT11c2VyK09VPWFjbWUsTz1IeXBlcmxlZGdlcixTVD1Ob3J0aCBDYXJvbGluYSxDPVVTOjpDTj1yb290LmNhc2VydmVyLE9VPVN1cHBvcnQsTz1BY21lLFNUPU5ldyBKZXJzZXksQz1VUw=="]}'
	chain.sh invoke
	set-chain-env.sh  -i '{"Args": ["issueTicket", "eDUwOTo6Q049YWNtZS1hZG1pbixPVT1jbGllbnQrT1U9YWNtZSxPPUh5cGVybGVkZ2VyLFNUPU5vcnRoIENhcm9saW5hLEM9VVM6OkNOPXJvb3QuY2FzZXJ2ZXIsT1U9U3VwcG9ydCxPPUFjbWUsU1Q9TmV3IEplcnNleSxDPVVT"]}'
	chain.sh invoke
	
}

function execInsert(){
	sleep 1
	amount=$1
	time node gatewayL.js insert $amount
}


function benchmark(){
	setGatewayDir
	for amount in 2 10 100 200 300 400 500 600 700 800 900 1000
	do
		echo ""
		echo "########################################"
		echo "Doing tests with $amount bids and offers"
		echo "########################################"
		echo ""
		for time in 1 2 3 4 5
		do
			execInsert $amount
		done
	done
	removeGatewayDir
}

function configUpdate(){
	cd network/caserver/
	./init.sh
	cd ..
	cd ..
	setGatewayDir
	node walletL.js add acme admin
	node walletL.js add acme anil
	node walletL.js add acme john
	node walletL.js add acme mary
	removeGatewayDir
	
}

[[ "${BASH_SOURCE[0]}" == "${0}" ]] && printHelp && exit

if [[ $# -lt 1 ]] ; then
  printHelp
  return
else
  MODE=$1
  shift
fi

if [ "${MODE}" == "set-env" ]; then
	set-env
elif [ "${MODE}" == "start" ]; then
	start
elif [ "${MODE}" == "deploy" ]; then
	deploy
elif [ "${MODE}" == "update" ]; then
	update
elif [ "${MODE}" == "register" ]; then
	register
elif [ "${MODE}" == "benchmark" ]; then
	#. ./setup.sh benchmark >> results/par_10.txt 2>&1
	benchmark
elif [ "${MODE}" == "configUpdate" ]; then
	configUpdate
elif [ "${MODE}" == "all" ]; then
	start
	deploy
	register
else
  printHelp
  return
fi
