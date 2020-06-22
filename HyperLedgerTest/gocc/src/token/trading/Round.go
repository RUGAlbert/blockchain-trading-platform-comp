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

// Init Implements the Init method
/* func (MP *MarketPlace) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Simply print a message
	fmt.Println("Init executed in history")

	// Return success
	return shim.Success(nil)
}

// Invoke method
func (MP *MarketPlace) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the function name and parameters
	funcName, args := stub.GetFunctionAndParameters()

	fmt.Println(args)

	if funcName == "SetNexLeader" {
		return SetNextLeader(stub)
	}

	return shim.Error("Bad Func Name!!!")
} */

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