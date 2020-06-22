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

	// Client Identity Library
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"

	// Standard go crypto package
	//"crypto/x509"
)

func ExecuteQuery(stub shim.ChaincodeStubInterface, query string) (string, error) {
	var pagesize int32 = 20
	bookmark := ""
	var counter uint64
	var pageCounter = 0
	var hasMorePages = true

	 // variables to hold query iterator and metadata
	 var qryIterator 	shim.StateQueryIteratorInterface
	 var queryMetaData 	*peer.QueryResponseMetadata
	 var err		error
	 resultJSON := "["
	 // start the pagination read loop
	 lastBookmark := ""
	 for hasMorePages {
		 // execute the rich query
		 qryIterator, queryMetaData, err = stub.GetQueryResultWithPagination(query, pagesize,bookmark)
		 if err != nil {
			 fmt.Printf("GetQueryResultWithPagination Error=" + err.Error())
			 return "", err
		 }
		 var resultKV *queryresult.KV
		 // Result read loop only if we received a different bookmark
		 
		 if lastBookmark != queryMetaData.Bookmark {
			 
			 for qryIterator.HasNext() {
 
				 // Get the next element
				 resultKV, err = qryIterator.Next()
				 
				 // Increment Counter
				 counter++
				if(counter != 1){
					resultJSON += ",\n"
				}
				resultJSON += string(resultKV.GetValue())
			 }
 
						 // Increment Page Counter
			 pageCounter++
 
			 fmt.Printf("Processed Page: %d \n", pageCounter)
		 } 
 
		 
		 // Get start key for the next page
		 bookmark = queryMetaData.Bookmark
 
		 // boomark = blank indicates no more records
		 hasMorePages = (bookmark != "" && lastBookmark != bookmark && bookmark != "nil")
		 lastBookmark = bookmark
 
		 // Close the iterator
		 qryIterator.Close()
	 }
	 resultJSON += "]"
	 // Total processed documents
	 fmt.Printf("Processed  Documents: %d \n", counter)
 
	 // Return the result JSON
	 return resultJSON, nil
}