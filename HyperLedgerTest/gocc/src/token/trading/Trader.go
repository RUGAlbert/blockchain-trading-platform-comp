package main

import (
	// April 2020, Updated to Fabric 2.0 Shim
	"github.com/hyperledger/fabric-chaincode-go/shim"
	peer "github.com/hyperledger/fabric-protos-go/peer"

	// Client Identity Library
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"encoding/json"
	"strconv"
	"errors"
)

type TraderData struct {
	DocType			string  `json:"docType"`
	Id     			string 	`json:"id"`
	TicketAmount 	uint 	`json:"ticketAmount"`
}

func IsRegisterdTrader(stub shim.ChaincodeStubInterface) bool {
	id, _ := cid.GetID(stub)
	matchAsBytes, err := stub.GetState("Trader:" + string(id))
	return !(err != nil || matchAsBytes == nil)
}

func GetTrader(stub shim.ChaincodeStubInterface, id string) (TraderData, error) {
	traderAsBytes, err := stub.GetState("Trader:"+id)
	var trader TraderData
	if err != nil {
		return trader, err
	}
	if traderAsBytes == nil {
		return trader, errors.New("There is no trader known by this address")
	}
	
	json.Unmarshal(traderAsBytes, &trader)
	return trader, nil
}

func SaveTrader(stub shim.ChaincodeStubInterface, trader TraderData, uid string) {
	jsonData, _ := json.Marshal(trader)
	id := "Trader:"+uid
	stub.PutState(id, jsonData)
}

func RegisterTrader(stub shim.ChaincodeStubInterface) peer.Response {
	uid, _ := cid.GetID(stub)
	count := 0

	if IsRegisterdTrader(stub) {
		return shim.Error("Already registerd")
	}


	/* countAsBytes, _ := stub.GetState("TraderCount")
	if countAsBytes != nil {
		count, _ = strconv.Atoi(string(countAsBytes))
	}
	count += 1 */
	stub.PutState("TraderCount", []byte(strconv.Itoa(count)))

	data := TraderData{DocType: "TraderData", Id:uid, TicketAmount:1}
	
	//id := "Trader:"+strconv.Itoa(count)
	SaveTrader(stub, data, uid)
	//issue one ticket

	tickets, _ := getTickets(stub)
	tickets = append(tickets, uid)
	saveTickets(stub, tickets)
	return shim.Success([]byte("Registed trader"))
}