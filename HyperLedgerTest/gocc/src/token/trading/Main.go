package main

/**
 * Demonstrates the use of CID
 **/
 import (
	// For printing messages on console
	"fmt"
	// April 2020, Updated to Fabric 2.0 Shim
	"github.com/hyperledger/fabric-chaincode-go/shim"
	peer "github.com/hyperledger/fabric-protos-go/peer"
)

// Init Implements the Init method
func (MP *MarketPlace) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Simply print a message
	fmt.Println("Init executed in history")

	// Return success
	return shim.Success(nil)
}

// Invoke method
func (MP *MarketPlace) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the function name and parameters
	funcName, args := stub.GetFunctionAndParameters()
	var response peer.Response
	switch funcName {
		case "publishOffer":
			response = MP.publishOffer(stub, args)
		case "publishBid":
			response = MP.publishBid(stub, args)
		case "getOffersAndBids":
			response = MP.getOffersAndBids(stub)
		case "getCurrentOffersAndBids":
			response = MP.getCurrentOffersAndBids(stub)
		case "setNextLeader":
			response = SetNextLeader(stub)
		case "setNextRound":
			response = SetNextRound(stub)
		case "addMatch":
			response = AddMatch(stub, args)
		case "getMatchesOfUser":
			response = GetMatchesOfUser(stub)
		case "getMatch":
			response = GetMatch(stub, args)
		case "acceptMatch":
			response = AcceptMatch(stub, args)
		case "rejectMatch":
			response = RejectMatch(stub, args)
		case "claimVolume":
			response = ClaimVolume(stub, args)
		case "confirmVolume":
			response = ConfirmVolume(stub, args)
		case "claimPayment":
			response = ClaimPayment(stub, args)
		case "confirmPayment":
			response = ConfirmPayment(stub, args)
		case "register":
			response = RegisterTrader(stub)
		case "issueTicket":
			response = IssueTicket(stub, args)
		case "isLeader":
			ans := IsLeader(stub)
			if ans {
				response = shim.Success([]byte("true"))
			} else {
				response = shim.Success([]byte("false"))
			}
		default:
			response = shim.Error("Bad Func Name!!!")
	}

	return response
}

// Chaincode registers with the Shim on startup
func main() {
	fmt.Printf("Started Chaincode. token/trading\n")
	err := shim.Start(new(MarketPlace))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
