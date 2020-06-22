/**
 * Demonstrates the use of Gateway Network & Contract classes
 */

// Needed for reading the connection profile as JS object
const fs = require('fs');
// Used for parsing the connection profile YAML file
const yaml = require('js-yaml');
// Import gateway class
const { Gateway, FileSystemWallet, DefaultEventHandlerStrategies, Transaction  } = require('fabric-network');

// Constants for profile
const CONNECTION_PROFILE_PATH = '../profiles/dev-connection.yaml'
// Path to the wallet
const FILESYSTEM_WALLET_PATH = './user-wallet'
// Identity context used
const USER_ID = 'mary@acme.com'
// Channel name
const NETWORK_NAME = 'airlinechannel'
// Chaincode
const CONTRACT_ID = "trading"



 // Get the requested action
let action='exec'
if (process.argv.length > 2){
    action = process.argv[2];
}

if(action == 'exec'){
    exec();
    return;
} else if (action == 'insert'){
    let numb = 100;
    if (process.argv.length > 3){
        numb = parseInt(process.argv[3]);
    }
    injectRandomBidsAndOffers(numb);
    return;
}
/**
 * Executes the functions for query & invoke
 */
async function exec() {

    // 2. Setup the gateway object
    var gateway = await setupGateway(USER_ID);
    // 3. Get the network
    let network = await gateway.getNetwork(NETWORK_NAME);
    // 5. Get the contract
    const contract = await network.getContract(CONTRACT_ID);
    // 6. Query the chaincode
    const allData = await getOffersAndBids(contract);
    var bidsAndOffers = splitOffersAndBids(allData);
    bids = bidsAndOffers['bids'];
    offers = bidsAndOffers['offers'];
    matches = MatchMaking(bids, offers);
    promises = [];
    for (let i = 0; i < matches.length; i++) {
        var m = matches[i];
        console.log("Doing " + i.toString());
        promises.push(submitData(contract, "addMatch", [m.offerAddress, m.bidAddress, m.volume.toString(), m.unitPrice.toString()]));
    }
    Promise.all(promises).then(() => {
        console.log("Submitted all matches");
        gateway.disconnect();
    });
}

async function injectRandomBidsAndOffers(amount){
    var gateway = await setupGateway(USER_ID);
    let network = await gateway.getNetwork(NETWORK_NAME);
    const contract = await network.getContract(CONTRACT_ID);
    let response = await contract.submitTransaction('setNextRound');
    await insertRandomAsUser("john@acme.com", NETWORK_NAME, CONTRACT_ID, "publishBid",amount/2);
    await insertRandomAsUser("anil@acme.com", NETWORK_NAME, CONTRACT_ID, "publishOffer",amount/2);
    var data = await getOffersAndBids(contract);
    if(data.length != amount){
        console.log("Something went wrong");
        console.log("Lenght " + data.length.toString());
    }
    gateway.disconnect();
}

async function insertRandomAsUser(userID, networkID, contractID, funcName, amount){
    let promises = [];
    var gatewayBid = await setupGateway(userID);
    let networkBid = await gatewayBid.getNetwork(networkID);
    const contractBid = await networkBid.getContract(contractID);
    //let response = await contractBid.submitTransaction('register');
    //const contractOffer = await getFastContract("anil@acme.com",NETWORK_NAME, CONTRACT_ID);
    for (let i = 0; i < amount; i++) {
        data = generateRandomData();
        promises.push(submitData(contractBid, funcName, data));
        //await submitData(contractBid, funcName, data);
    }

    await Promise.all(promises).then(() => {
        console.log("Done all for " + userID + " using " + funcName);
        gatewayBid.disconnect();
    });
}

function randomDate(start, end) {
    return new Date(start.getTime() + Math.random() * (end.getTime() - start.getTime()));
}

function ISODateString(d){
    function pad(n){return n<10 ? '0'+n : n}

    return d.getUTCFullYear()+'-'
         + pad(d.getUTCMonth()+1)+'-'
         + pad(d.getUTCDate())+'T'
         + pad(d.getUTCHours())+':'
         + pad(d.getUTCMinutes())+':'
         + pad(d.getUTCSeconds())+'+00:00'
}
   

function generateRandomData(){
    //"publishBid","1","2","2012-11-01T22:08:41+00:00","1"
    var minVol = (Math.floor(Math.random() * 100));
    var maxVol = (minVol + Math.floor(Math.random() * 100)).toString();
    var startDate = new Date();
    var date = ISODateString(randomDate(startDate, new Date(startDate.getTime()+86400*1000)));
    var price = (Math.floor(Math.random() * 10)).toString();
    return [minVol.toString(), maxVol, date, price]
}

/**
 * Submit the transaction
 * @param {object} contract 
 */
async function submitData(contract, name, data){
    try{
        // Submit the transaction
        let response = await contract.submitTransaction(name, data[0],data[1],data[2],data[3]);
        //console.log("Submit Response=",response.toString());
    } catch(e){
        // fabric-network.TimeoutError
        console.log(e);
    }
}


/**
 * Queries the chaincode
 * @param {object} contract 
 */
async function queryContract(contract){
    try{
        // Query the chaincode
        let response = await contract.evaluateTransaction('getCurrentOffersAndBids');
        console.log(`Query Response=${response.toString()}`);
    } catch(e){
        console.log(e);
    }
}

/**
 * Queries the chaincode
 * @param {object} contract 
 */
async function getOffersAndBids(contract){
    var data;
    try{
        // Query the chaincode
        let response = await contract.evaluateTransaction('getCurrentOffersAndBids');
        //data = JSON.parse("[" + response.toString() + "]");
        data = JSON.parse(response.toString());
    } catch(e){
        console.log(e);
    }

    return data;
}

/**
 * Split offers and bids
 * @param {object} contract 
 */
function splitOffersAndBids(data){
    var offers = [];
    var bids = [];
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        if(obj.isOffer == true){
            offers.push(obj);
        }else{
            bids.push(obj);
        }
    }

    var offersAndBids = {};
    offersAndBids['offers'] = offers;
    offersAndBids['bids'] = bids;
    return offersAndBids;
}

function MatchMaking(bids, offers){
    var matches = [];
    for (let i = 0; i < bids.length; i++) {
        var bid = bids[i];
        if((bid.isMatched || false)){
            continue;
        }
        for (let j = 0; j < offers.length; j++) {
            var offer = offers[j];
            if((offer.isMatched || false)){
                continue;
            }

            var matchQuantity = calcQuantity(bid, offer);
            var matchPrice = calcPrice(bid, offer);

            if(matchQuantity > 0 && matchPrice > 0){
                offer.isMatched = true;
                bid.isMatched = true;
                var match = {};
                match["offerAddress"] = offer.owner;
                match["bidAddress"] = bid.owner;
                match["volume"] = matchQuantity;
                match["unitPrice"] = matchPrice;
                matches.push(match);
            }
        }
    }
    return matches;
}

function calcQuantity(bid, offer){
    if(bid.maxVolume < offer.minVolume || offer.maxVolume < bid.minVolume){
        return 0;
    }

    if(offer.maxVolume < bid.maxVolume){
        return offer.maxVolume;
    }else{
        return bid.maxVolume;
    }
}

//asumes that offering party is alright with more money and bidding party with less money
function calcPrice(bid, offer){
    if(bid.maxUnitPrice < offer.minUnitPrice){
        return 0;
    }

    return offer.minUnitPrice;
}

/**
 * Submit the transaction
 * @param {object} contract 
 */
async function submitTxnContract(contract){
    try{
        // Submit the transaction
        let response = await contract.submitTransaction('transfer', 'john','sam','2');
        console.log("Submit Response=",response.toString());
    } catch(e){
        // fabric-network.TimeoutError
        console.log(e);
    }
}

/**
 * Function for setting up the gateway
 * It does not actually connect to any peer/orderer
 */
async function setupGateway(userID) {

    var gateway = new Gateway()
    
    // 2.1 load the connection profile into a JS object
    let connectionProfile = yaml.safeLoad(fs.readFileSync(CONNECTION_PROFILE_PATH, 'utf8'));

    // 2.2 Need to setup the user credentials from wallet
    const wallet = new FileSystemWallet(FILESYSTEM_WALLET_PATH)

    // 2.3 Set up the connection options
    let connectionOptions = {
        identity: userID,
        wallet: wallet,
        discovery: { enabled: false, asLocalhost: true }
        /*** Uncomment lines below to disable commit listener on submit ****/
        // , eventHandlerOptions: {
        //     strategy: null
        // } 
    }

    // 2.4 Connect gateway to the network
    await gateway.connect(connectionProfile, connectionOptions)
    // console.log( gateway)
    return gateway
}



/**
 * Creates the transaction & uses the submit function
 * Solution to exercise
 * To execute this add the line in main() => submitTxnTransaction(contract)
 * @param {object} contract 
 */
// async function submitTxnTransaction(contract) {
//     // Provide the function name
//     let txn = contract.createTransaction('transfer')
    
//     // Get the name of the transaction
//     console.log(txn.getName())

//     // Get the txn ID
//     console.log(txn.getTransactionID())

//     // Submit the transaction
//     try{
//         let response = await txn.submit('john', 'sam', '2')
//         console.log("transaction.submit()=", response.toString())
//     } catch(e) {
//         console.log(e)
//     }
// }
