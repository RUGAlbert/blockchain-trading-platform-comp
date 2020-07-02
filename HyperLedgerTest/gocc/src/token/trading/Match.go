package main

/**
 * Demonstrates the use of CID
 **/
import (
	// For printing messages on console
	"fmt"
	// April 2020, Updated to Fabric 2.0 Shim
	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/rs/xid"
	
	peer "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"

	"encoding/json"
	"reflect"

	// Standard go crypto package
	"strconv"
)

//All data used to save a match
type MatchData struct {
	Id				string  	`json:"id"`
	DocType			string  	`json:"docType"`
	OfferAddress 	string 		`json:"offerAddress"`
	BidAddress		string 		`json:"bidAddress"`
	Volume			float64 	`json:"volume"`
	UnitPrice		float64 	`json:"unitPrice"`
	AcceptedBid 	bool	 	`json:"acceptedBid"`
	AcceptedOffer 	bool	 	`json:"acceptedOffer"`
	Stage 			int 		`json:"stage"`
}

//It called with the following args: 0: offerAdress, 1: bidAddress, 2: volume, 3: unitPrice
func AddMatch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if !IsLeader(stub){
		return shim.Error("Caller is not leader")
	}

	offerAddress := args[0]
	bidAddress := args[1]
	volume, _ := strconv.ParseFloat(args[2], 64)
	unitPrice, _ := strconv.ParseFloat(args[3], 64)
	id := xid.New().String()
	data := MatchData{Id: id, DocType: "MatchData", OfferAddress:offerAddress, BidAddress:bidAddress, Volume: volume, UnitPrice:unitPrice, AcceptedBid:false, AcceptedOffer:false, Stage:0}
	jsonData, _ := json.Marshal(data)
	stub.PutState(id, jsonData)
	return shim.Success([]byte("Created match"))
}


//Get all matches of the user
func GetMatchesOfUser(stub shim.ChaincodeStubInterface) peer.Response {
	id, _ := cid.GetID(stub)
	qry := "{"
	qry += "\"selector\": {"
	qry += "\"docType\": \"MatchData\","
	qry += "\"$or\": ["
	qry += "{\"offerAddress\": \"" + id + "\"},"
	qry += "{\"bidAddress\": \"" + id + "\"}"
	qry += "]}}"
	resultJSON, err := ExecuteQuery(stub, qry)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(resultJSON))
}

//as argument the id of the accepted match to get a certain match
func GetMatch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	resultJSON, err := ExecuteQuery(stub, "{\"selector\": { \"docType\": \"MatchData\", \"id\":\"" + args[0] + "\"}}")
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(resultJSON))
}

//get match based on a string with an expected stage, otherwise it will return an error
func getMatchForUse(stub shim.ChaincodeStubInterface, id string, expectedStage int) (MatchData, error){
	var result MatchData

	matchAsBytes, err := stub.GetState(id)

	if err != nil {
		return result, err
	}

	if matchAsBytes == nil {
		return result, fmt.Errorf("No Matchdata found for this id")
	}

	json.Unmarshal(matchAsBytes, &result)

	if result.Stage != expectedStage {
		return result, fmt.Errorf("Match no longer pending")
	}

	return result, nil
}

//This will be used to accept a match, based on the id id of the match
func AcceptMatch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	id, _ := cid.GetID(stub)

	match, err := getMatchForUse(stub, args[0], 0)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	if reflect.DeepEqual([]byte(match.OfferAddress),[]byte(id)) {
		match.AcceptedOffer = true
	} else if reflect.DeepEqual([]byte(match.BidAddress),[]byte(id)) {
		match.AcceptedBid = true
	} else {
		return shim.Error("You are not part of this deal")
	}

	if match.AcceptedOffer == true && match.AcceptedBid == true {
		match.Stage = 1;
	}
	
	matchAsBytes, _ := json.Marshal(match)
	stub.PutState(args[0],matchAsBytes);

	return shim.Success([]byte("Accepted match"))
}

//reject a certain match, can only be done if not both parties have accepted it
func RejectMatch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	id, _ := cid.GetID(stub)

	match, err := getMatchForUse(stub, args[0], 0)
	if err != nil {
		return shim.Error(err.Error())
	}
	

	if reflect.DeepEqual([]byte(match.OfferAddress),[]byte(id)) {
		match.Stage = -1
	} else if reflect.DeepEqual([]byte(match.BidAddress),[]byte(id)) {
		match.Stage = -1
	} else {
		return shim.Error("You are not part of this deal")
	}

	matchAsBytes, _ := json.Marshal(match)
	stub.PutState(args[0],matchAsBytes);

	return shim.Success([]byte("Rejected match"))
}

//Claims a match and if it is allowed let is go to the next stage
func claimer(stub shim.ChaincodeStubInterface, currentStage int, isOffer bool, mid string) peer.Response{
	id, _ := cid.GetID(stub)

	match, err := getMatchForUse(stub, mid, currentStage)
	if err != nil {
		return shim.Error(err.Error())
	}
	

	if isOffer && !reflect.DeepEqual([]byte(match.OfferAddress),[]byte(id)) {
		return shim.Error("You are not the correct person to make this call")
	}
	if !isOffer && !reflect.DeepEqual([]byte(match.BidAddress),[]byte(id)) {
		return shim.Error("You are not the correct person to make this call")
	}


	match.Stage = currentStage + 1

	matchAsBytes, _ := json.Marshal(match)
	stub.PutState(mid,matchAsBytes);

	return shim.Success([]byte("Went to next stage"))
}

//Claim that the volume is send
func ClaimVolume(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	return claimer(stub, 1, true, args[0])
}

//Confirm that the volume is send
func ConfirmVolume(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	return claimer(stub, 2, false, args[0])
}

//Claim that the payment is send
func ClaimPayment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	return claimer(stub, 3, false, args[0])
}

//Confirm that the payment is send
func ConfirmPayment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	return claimer(stub, 4, true, args[0])
}

