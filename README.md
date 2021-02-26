# BlockHealth
BlockHealth is a solution for active surveillance and sharing of personal health data with not only tamper proofing and data protection guarantees, but also privacy-preserving (data can only be accessed by legitimate staff and, if required, it can be also deleted).  To this purpose, BlockHealth does not insert in the Blockchain the health data, but rather only their secure hash. Actual health data must be exclusively stored in a private database owned by (and under the control of) each different actors.

# Create BlockHealth Network
First of all you need to create cryptographic material (x509 certs and signing keys) for yours various network actors modifying the crypto-config.yaml file and using the cryptogen tool to generate them.
Now you must create the configuration artifacts modifying configtx.yaml file that contains the definitions for the sample network and using the configtxgen tool. In this case, the files that you'll create are: 
- orderer genesis block,
- channel configuration transaction,
- and three anchor peer transactions - one for each Peer Organization.

The profiles used are:
- "OrgsCompanies" to create channel configuration with the three organizations;
- "MultiNodeEtcdRaft" to create genesis block of network configuration. 

Remember: Pay specific attention to the “Profiles” section at the bottom of configtx.yaml file.

After that, you can run the network. You will notice that 3 peers and 3 CAs (one for each organization) and 3 Raft orderers for the quorum will be started. The image used is v2.3.1 for all peers and latest for CAs.

WARNING: We deployed the peers to three different servers with an orderer each.


# Crate channel
To create the channel, we used the CLI fabric-tool (latest version) and we will explain how to use this to create the infrastructure.

The CLI to work, you need to preface four environment variables given below. These informations are essential to indicate the cli to which peer it must connect.
Example in this case:

CompanyA
- CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/companyA.blockhealth.com/users/Admin@companyA.blockhealth.com/msp
- CORE_PEER_ADDRESS=peer0.companyA.blockhealth.com:7051
- CORE_PEER_LOCALMSPID="CompanyAMSP"
- CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/companyA.blockhealth.com/peers/peer0.companyA.blockhealth.com/tls/ca.crt

Now we can exec the command below to create a channel:
$peer channel create -o orderer.blockhealth.com:7050 -c $CHANNEL_NAME -f $PATH_TO/channel.tx --tls --cafile $PATH_TO/tlsca.blockhealth.com-cert.pem

After that we can join CompanyA peer to the channel:
$peer channel join -b $CHANNEL_NAME.block

These operations must be performed for all organizations.


#Anchor peer update
At this point it is essential to update the anchor peers like as:

$peer channel update -o orderer.blockhealth.com:7050 -c $CHANNEL_NAME -f $PATH_TO/CompanyAMSPanchors.tx --tls --cafile $PATH_TO/tlsca.blockhealth.com-cert.pem



#Install the chaincode

First, you need to install the Go (in this case) chaincode on every peer that will execute and endorse your transactions.
The members of the channel need to agree the chaincode definition that establishes chaincode governance.

You need to package the chaincode before it can be installed on peers. 
$peer lifecycle chaincode package $CHAINCODE.tar.gz --path $PATH_TO/chaincode/test --lang golang --label $CHAINCODE_1.0

After that, you need to provide a chaincode package label as a description of the chaincode. Then you can approve the chaincode definition:
$peer lifecycle chaincode approveformyorg --channelID $CHANNEL_NAME --name $CHAINCODE --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $PATH_TO/tlsca.blockhealth.com-cert.pem

Now we provided a --signature-policy argument to the command above to set the chaincode endorsement policy. In this case, the policy will require an endorsement from a peer belonging to CompanyA AND CompanyB AND CompanyC (i.e. three endorsements).

Since all channel members have approved the definition, you can now commit it to the channel like as below:
$peer lifecycle chaincode commit -o orderer.blockhealth.com:7050 --channelID $CHANNEL_NAME --name $CHAINCODE --version 1.0 --sequence 1 --tls --cafile "$PATH_TO/tlsca.blockhealth.com-cert.pem" --peerAddresses peer0.companyA.blockhealth.com:7051 --tlsRootCertFiles "$PATH_TO/ca.crt" --peerAddresses peer0.companyB.blockhealth.com:8051 --tlsRootCertFiles "$PATH_TO/ca.crt"      --peerAddresses peer0.companyC.blockhealth.com:9051 --tlsRootCertFiles "$PATH_TO/ca.crt"

Now you can invoke the chaincode.

WARNING: To resolve the logical host names, we configured the /etc/hosts file
