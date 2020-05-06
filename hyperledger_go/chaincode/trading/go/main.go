package main

import (
	"fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	
	smartContract := new(SmartContract)
	simpleContract := new(SimpleContract)
	simpleContract.Name = "SimpleContract"
	smartContract.Name = "SmartContract"

	chaincode, err := contractapi.NewChaincode(smartContract, simpleContract)

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}
	chaincode.DefaultContract = smartContract.GetName()
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
