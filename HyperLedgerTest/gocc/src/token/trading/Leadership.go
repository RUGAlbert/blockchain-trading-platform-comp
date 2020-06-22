package main

/**
 * Demonstrates the use of CID
 **/
import (
	// For printing messages on console
	"fmt"
	"reflect"
	//"time"
	// April 2020, Updated to Fabric 2.0 Shim
	"github.com/hyperledger/fabric-chaincode-go/shim"
	
	peer "github.com/hyperledger/fabric-protos-go/peer"

	// Client Identity Library
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	//"github.com/hyperledger/fabric-protos-go/ledger/queryresult"

	// Standard go crypto package
	//"crypto/x509"
	"errors"
	"strconv"
	"encoding/json"
	"math/rand"
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

func IsLeader(stub shim.ChaincodeStubInterface) bool {
	id, _ := cid.GetID(stub)
	currentId, _ := stub.GetState("Leader")
	return reflect.DeepEqual(currentId,[]byte(id))
}

func SetNextLeader(stub shim.ChaincodeStubInterface) peer.Response {
	/*if !IsLeader(stub){
		fmt.Printf("Is not leader\n")
		return shim.Error("Caller is not leader")
	}*/
	id, err := GetRandomLeaderWithTicket(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	sts := stub.PutState("Leader",[]byte(id))
	fmt.Printf("sts: %s\n",sts);
	return shim.Success([]byte("Selected next leader"))
}

/* func getRandomLeader(stub shim.ChaincodeStubInterface) (string, error) {
	times, _ := stub.GetTxTimestamp()

	rand.Seed(int64(times.Nanos))
	countAsBytes, _ := stub.GetState("TraderCount")
	count := 0
	if countAsBytes == nil {
		return "", errors.New("There is no registerd trader")
	}
	count, _ = strconv.Atoi(string(countAsBytes))
	random := rand.Intn(count)
	traderAsBytes, _ := stub.GetState("Trader:"+strconv.Itoa(random))

	var newLeader TraderData
	_ = json.Unmarshal(traderAsBytes, &newLeader)
	fmt.Printf("qry: %s\n",newLeader.Id);
	return newLeader.Id, nil
} */

func getTickets(stub shim.ChaincodeStubInterface) ([]string, error) {
	ticketsAsBytes, _ := stub.GetState("Tickets")
	var tickets []string
	if ticketsAsBytes == nil {
		return tickets, errors.New("There are no tickets")
	}
	json.Unmarshal(ticketsAsBytes, &tickets)
	return tickets, nil
}

func GetRandomLeaderWithTicket(stub shim.ChaincodeStubInterface) (string, error) {
	tickets, err := getTickets(stub)
	if err != nil {
		return "", err
	}
	index := rand.Intn(len(tickets))
	newLeaderId := tickets[index]
	leader, _ := GetTrader(stub, newLeaderId)
	if leader.TicketAmount > 1 {
		removeLeaderFromTickets(stub, tickets, index)
	}
	return newLeaderId, nil
}

func removeLeaderFromTickets(stub shim.ChaincodeStubInterface, tickets []string, index int){
	tickets = remove(tickets, index)
	saveTickets(stub, tickets)
}

func remove(s []string, i int) []string {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func saveTickets(stub shim.ChaincodeStubInterface, tickets []string){
	ticketsAsBytes, _ := json.Marshal(tickets)
	stub.PutState("Tickets", ticketsAsBytes)
}

func IssueTicket(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	traderId := args[0]
	trader, err := GetTrader(stub, traderId)
	if err != nil {
		return shim.Error(err.Error())
	}
	trader.TicketAmount += 1
	tickets, _ := getTickets(stub)
	tickets = append(tickets, traderId)
	SaveTrader(stub, trader, traderId)
	saveTickets(stub, tickets)

	return shim.Success([]byte("Issued Ticket"))
}

func IssueTickets(stub shim.ChaincodeStubInterface, args []string) peer.Response{
	count, _ := strconv.Atoi(args[1])
	for i := 0; i < count; i++ {
		IssueTicket(stub, args)
	}
	return shim.Success([]byte("Issued Tickets"))
}

//will be removed later
func SetFirstLeader(stub shim.ChaincodeStubInterface) peer.Response {
	id, _ := cid.GetID(stub)
	fmt.Printf("id: %s\n",id);
	stub.PutState("Leader",[]byte(id))
	return shim.Success([]byte("Selected next leader"))
}