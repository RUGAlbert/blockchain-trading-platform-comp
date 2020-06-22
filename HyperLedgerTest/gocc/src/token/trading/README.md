Client Identity Library
=======================
https://github.com/hyperledger/fabric-chaincode-go/tree/master/pkg/cid

// do this in the bin folder of network and gocc
export PATH=$PATH:${PWD}
export GOPATH=${PWD}

# Vendoring
cd $GOPATH/src/token/trading/
./govendor.sh

=============
Users in acme
=============
Please check out the attributes associated with the various users.
These identities have been discused in the lecture.

#  mary
"app.accounting.role=tradeapprover:ecert","department=accounting:ecert"
#  john
"app.accounting.role=accountant:ecert","department=accounting:ecert"
# anil
"department=logistics:ecert"

======================
Set up the environment
======================
# Make sure the context is admin otherwise the instantiate will Fail
# with endrosement error
source reset-chain-env.sh
. set-env.sh acme
. set-ca-msp.sh admin


ca-dev-init.sh
set-chain-env.sh -n trading -v 1.0  -p token/trading -c '{"Args":["Init"]}' -I false

# Generate the package & install
chain.sh install -p

# Instantiate the chaincode
chain.sh instantiate

==========================================
Query - publish offer
==========================================

set-chain-env.sh  -i '{"Args": ["publishBid","1","2","2012-11-01T22:08:41+00:00","1"]}'

. set-ca-msp.sh  mary
chain.sh invoke

. set-ca-msp.sh  john
chain.sh invoke

. set-ca-msp.sh  anil
chain.sh invoke

set-chain-env.sh  -i '{"Args": ["getCurrentOffersAndBids"]}'

set-chain-env.sh  -i '{"Args": ["setNextLeader"]}'

set-chain-env.sh  -i '{"Args": ["setNextRound"]}'

set-chain-env.sh  -i '{"Args": ["rejectMatch","br2ee4lgj8914v29a460"]}'

set-chain-env.sh  -i '{"Args": ["setFirstLeader"]}'
chain.sh invoke

set-chain-env.sh  -i '{"Args": ["register"]}'
chain.sh invoke

set-chain-env.sh  -i '{"Args": ["setNextRound"]}'
chain.sh invoke


. set-ca-msp.sh  admin
chain.sh install -p
set-chain-env.sh -s 
chain.sh instantiate