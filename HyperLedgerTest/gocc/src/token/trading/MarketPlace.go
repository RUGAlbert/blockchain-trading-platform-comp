package main

/**
 * Demonstrates the use of CID
 **/
import (
	// For printing messages on console
	"fmt"
	"time"
	// April 2020, Updated to Fabric 2.0 Shim
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/rs/xid"
	peer "github.com/hyperledger/fabric-protos-go/peer"

	// Client Identity Library
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"

	// Standard go crypto package
	//"crypto/x509"
	"strconv"
	"encoding/json"
)

// CidChaincode Represents our chaincode object
type MarketPlace struct {
}

// 
type OfferOrBidData struct {
	DocType				string  	`json:"docType"`
	Round 				int 		`json:"round"`
	MinVolume			float64 	`json:"minVolume"`
	MaxVolume			float64 	`json:"maxVolume"`
	ExpirationDateTime	time.Time 	`json:"expirationDateTime"`
	MinUnitPrice 		float64 	`json:"minUnitPrice"`
	MaxUnitPrice 		float64 	`json:"maxUnitPrice"`
	IsOffer 			bool 		`json:"isOffer"`
	Owner     			string 		`json:"owner"`
}

// publishoffer adds an offers
// as arguments the following should be added 0 = minVol, 1 = maxVol, 2 = expirationDate, 3 = minUnitPrice
func (MP *MarketPlace) publishOffer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	
	return MP.insertData(stub, args[0], args[1],args[2],args[3],"0",true)
}

// publishBids adds a bid
// as arguments the following should be added 0 = minVol, 1 = maxVol, 2 = expirationDate, 3 = maxUnitPrice
func (MP *MarketPlace) publishBid(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	
	return MP.insertData(stub, args[0], args[1],args[2],"0",args[3],false)
}

//inserts data, an offer or bid
func (MP *MarketPlace)insertData(stub shim.ChaincodeStubInterface, minVol string, maxVol string, date string, minPrice string, maxPrice string, isOffer bool) peer.Response{
	//get unique id of user
	uid, _ := cid.GetID(stub)
	if !IsRegisterdTrader(stub){
		return shim.Error("You are not a registerd trader, please register first")
	}
	minVolume, _ := strconv.ParseFloat(minVol, 64)
	maxVolume, _ := strconv.ParseFloat(maxVol, 64)
	layout := "2006-01-02T15:04:05Z07:00"
	expirationDateTimeConv, err := time.Parse(layout, date)

	if err != nil {
		fmt.Printf("Date parse error= %s",  err.Error())
		return shim.Error("Date parse error=" +  err.Error())
	}
	
	minUnitPrice, _ := strconv.ParseFloat(minPrice, 64)
	maxUnitPrice, _ := strconv.ParseFloat(maxPrice, 64)
	data := OfferOrBidData{DocType: "OfferOrBidData", MinVolume:minVolume, MaxVolume:maxVolume, ExpirationDateTime: expirationDateTimeConv, MinUnitPrice:minUnitPrice, MaxUnitPrice:maxUnitPrice, IsOffer:isOffer, Owner:uid, Round:GetCurrentRound(stub)}
	jsonData, _ := json.Marshal(data)
	id := xid.New().String()
	stub.PutState(id, jsonData)
	return shim.Success([]byte("Created offer or bid"))
}

//get all offers and bids
func (MP *MarketPlace)getOffersAndBids(stub shim.ChaincodeStubInterface) peer.Response {
	resultJSON, err := ExecuteQuery(stub, "{\"selector\": { \"docType\": \"OfferOrBidData\"}}")
	if err != nil {
		return shim.Error(err.Error())
	}
	 // Return the result JSON
	 return shim.Success([]byte(resultJSON))
}

//get all offers and bids in current round
func (MP *MarketPlace)getCurrentOffersAndBids(stub shim.ChaincodeStubInterface) peer.Response {
	query := "{\"selector\": { \"docType\": \"OfferOrBidData\",\"round\": "+ strconv.Itoa(GetCurrentRound(stub)) +"}}"
	fmt.Printf("%s\n",query)
	resultJSON, err := ExecuteQuery(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	 // Return the result JSON
	 return shim.Success([]byte(resultJSON))
}
