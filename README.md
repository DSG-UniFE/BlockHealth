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

WARNING: We deployed the peers to three different servers each with an orderer.


# Crate channel

