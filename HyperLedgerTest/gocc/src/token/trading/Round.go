package main

/**
 * Demonstrates the use of CID
 **/
import (
	// For printing messages on console
	"fmt"
	//"time"
	// April 2020, Updated to Fabric 2.0 Shim
	"github.com/hyperledger/fabric-chaincode-go/shim"
	
	peer "github.com/hyperledger/fabric-protos-go/peer"

	// Client Identity Library
	//"github.com/hyperledger/fabric-protos-go/ledger/queryresult"

	// Standard go crypto package
	//"crypto/x509"
	"strconv"
	//"encoding/json"
)

//Get the current round
func GetCurrentRound(stub shim.ChaincodeStubInterface) int {
	currentRound, _ := stub.GetState("Round")
	var val int
	if currentRound == nil {
		val = 0
	}else{
		val, _ = strconv.Atoi(string(currentRound))
	}

	return val
}

//go to the next round
func SetNextRound(stub shim.ChaincodeStubInterface) peer.Response {
	if !IsLeader(stub){
		fmt.Printf("Is not leader\n")
		return shim.Success([]byte("Error: Caller is not the leader."))
	}
	val := GetCurrentRound(stub)
	val += 1;
	stub.PutState("Round", []byte(strconv.Itoa(val)))
	return shim.Success([]byte("Next round found"))
}
