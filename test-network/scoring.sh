#!/bin/bash

# Initialize score
score=0
max_score=50

# Load Org1 peer context (assuming script is provided)
source ./scripts/setOrgPeerContext.sh 1

# Set environment paths
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

# Function to check if the peer has joined "mychannel"
check_peer_channel() {
    peer channel list > channel_list_output.txt 2>&1
    if grep -q "mychannel" channel_list_output.txt; then
        echo "Peer is part of mychannel."
        return 0
    else
        echo "Peer is not part of mychannel."
        return 1
    fi
}

# Check if Docker container for new peer is running
new_peer="peer1.org1.example.com"
docker ps | grep "$new_peer" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "Docker container for $new_peer is running."
    score=$((score + 10))
else
    echo "Docker container for $new_peer is not running."
fi

# Check if the peer has joined the mychannel
check_peer_channel
if [ $? -eq 0 ]; then
    score=$((score + 10))
else
    echo "Peer has not joined mychannel."
fi

# Check if chaincode is installed on the new peer
CHAINCODE_NAME="transactionlog"
peer lifecycle chaincode queryinstalled > installed_chaincodes.txt 2>&1
if grep -q "$CHAINCODE_NAME" installed_chaincodes.txt; then
    echo "Chaincode $CHAINCODE_NAME is installed on $new_peer."
    score=$((score + 10))
else
    echo "Chaincode $CHAINCODE_NAME is not installed on $new_peer."
fi

# Check if the logTransaction function works (Transaction Logging)
transaction_id="tx12345"
sender="Alice"
receiver="Bob"
amount=100


# Verify the transaction log with a query
echo "Querying the transaction log..."
QUERY_OUTPUT=$(peer chaincode query -C mychannel -n $CHAINCODE_NAME -c '{"Args":["QueryTransaction","'"$transaction_id"'"]}' 2>&1)
echo "Query output: $QUERY_OUTPUT"

if [[ $QUERY_OUTPUT == *"$transaction_id"* && $QUERY_OUTPUT == *"$sender"* && $QUERY_OUTPUT == *"$receiver"* && $QUERY_OUTPUT == *"$amount"* ]]; then
    echo "Transaction query successful."
    score=$((score + 20))
else
    echo "Transaction query failed."
fi

# Final score output
echo "Final Score: $score/$max_score"

# Exit
exit 0
