# Hyperledger Fabric on Raspberry Pi

A distributed Hyperledger Fabric blockchain network implementation running on Raspberry Pi devices, designed for healthcare data management with performance benchmarking capabilities.

## Overview

This project demonstrates a multi-organization Hyperledger Fabric network deployed across multiple Raspberry Pi devices. The network implements a healthcare record management system with smart contracts for creating, reading, updating, and deleting medical records.

## Network Architecture

### Physical Infrastructure
- **Pi 1**: Orderer node (orderer1.example.com), Org1 Peer (peer0.org1.example.com), Org1 CA
- **Pi 2**: Org2 Peer (peer0.org2.example.com), Org2 CA
- **Pi 3**: Org3 Peer (peer0.org3.example.com)

### Network Components
- **Channel**: mychannel
- **Consensus**: Raft-based ordering service
- **Organizations**: 3 organizations (Org1, Org2, Org3)
- **Fabric Version**: 2.5
- **CA Version**: 1.5

## Directory Structure

```
hyperledger-fabric-raspberry/
├── fabric-network-pi1/          # Pi 1 configuration and deployment
│   ├── docker-compose-pi1.yaml  # Docker services for orderer and Org1 peer
│   ├── configtx.yaml            # Network channel configuration
│   ├── crypto-config.yaml       # Cryptographic material generation config
│   ├── chaincode/               # Smart contract implementations
│   │   └── healthcc/            # Healthcare chaincode (Go)
│   ├── crypto-config/           # PKI certificates and keys
│   ├── channel-artifacts/       # Channel genesis blocks and transactions
│   ├── caliper-healthcare/      # Performance benchmarking suite
│   │   ├── benchmarks/          # Benchmark configuration files
│   │   ├── networks/            # Network connection profiles
│   │   └── workload/            # Transaction workload modules
│   └── orderer1-data/           # Persistent orderer ledger data
│
├── fabric-network-pi2/          # Pi 2 configuration and deployment
│   ├── docker-compose-pi2.yaml  # Docker services for Org2 peer
│   ├── chaincode-healthcare/    # Healthcare smart contract (Go)
│   ├── crypto-config/           # Org2 cryptographic material
│   └── peer0-org2-data/         # Persistent peer ledger data
│
└── fabric-network-pi3/          # Pi 3 configuration and deployment
    ├── docker-compose.yaml      # Docker services for Org3 peer
    └── crypto-config/           # Org3 cryptographic material
```

## Smart Contracts

### Healthcare Chaincode
The network implements healthcare record management with the following functions:

- **CreateRecord**: Create a new medical record
- **ReadRecord**: Retrieve a medical record by ID
- **UpdateRecord**: Update an existing medical record
- **DeleteRecord**: Remove a medical record
- **RecordExists**: Check if a record exists
- **GetAllRecords**: Query all records in the ledger

### Data Model
```go
type MedicalRecord struct {
    RecordID    string
    PatientName string
    Diagnosis   string
    Treatment   string
    Timestamp   string
}
```

## Network Configuration

### IP Addresses
- **Pi 1 (Orderer + Org1)**: 100.89.132.94
- **Pi 2 (Org2)**: 100.81.64.92
- **Pi 3 (Org3)**: To be configured

### Port Mappings
#### Pi 1
- Orderer: 7050
- Org1 Peer: 7051
- Org1 CA: 7054
- Operations/Metrics: 9444-9448

#### Pi 2
- Org2 Peer: 8051
- Org2 CA: 8054
- Operations: 9449

## Prerequisites

- Raspberry Pi 3B+
- Ubuntu Server or Raspberry Pi OS (64-bit)
- Docker and Docker Compose installed
- At least 4GB RAM per device
- Network connectivity between all Pi devices
- Go 1.19+ (for chaincode development)
- Node.js 14+ (for Caliper benchmarking)

## Installation & Setup

### 1. Clone Repository
```bash
git clone https://github.com/Agfiras/hyperledger-fabric-raspberry.git
cd hyperledger-fabric-raspberry
```

### 2. Generate Cryptographic Material
On Pi 1:
```bash
cd fabric-network-pi1
cryptogen generate --config=./crypto-config.yaml
```

### 3. Generate Channel Artifacts
```bash
configtxgen -profile ChannelUsingRaft -outputBlock ./channel-artifacts/genesis.block -channelID system-channel
configtxgen -profile ChannelUsingRaft -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID mychannel
configtxgen -profile ChannelUsingRaft -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP
configtxgen -profile ChannelUsingRaft -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
```

### 4. Start Network Components

On Pi 1:
```bash
cd fabric-network-pi1
docker-compose -f docker-compose-pi1.yaml up -d
```

On Pi 2:
```bash
cd fabric-network-pi2
docker-compose -f docker-compose-pi2.yaml up -d
```

On Pi 3:
```bash
cd fabric-network-pi3
docker-compose up -d
```

### 5. Create and Join Channel
```bash
# Create channel
peer channel create -o orderer1.example.com:7050 -c mychannel -f ./channel-artifacts/mychannel.tx

# Join peers to channel
peer channel join -b mychannel.block
```

### 6. Deploy Chaincode
```bash
# Package chaincode
peer lifecycle chaincode package healthcc.tar.gz --path ./chaincode/healthcc --lang golang --label healthcc_1

# Install on each peer
peer lifecycle chaincode install healthcc.tar.gz

# Approve for each organization
peer lifecycle chaincode approveformyorg -o orderer1.example.com:7050 --channelID mychannel --name healthcc --version 1.0 --package-id <PACKAGE_ID> --sequence 1

# Commit chaincode
peer lifecycle chaincode commit -o orderer1.example.com:7050 --channelID mychannel --name healthcc --version 1.0 --sequence 1
```

## Performance Benchmarking

The project includes Hyperledger Caliper for performance testing with multiple configurations:

### Available Benchmarks
- 5 workers: `all-experiments-5workers.yaml`

### Run Benchmarks
```bash
cd fabric-network-pi1/caliper-healthcare
npm install
npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/all-experiments-10workers.yaml --caliper-flow-only-test
```

### Workload Modules
- **createRecord.js**: Creates new medical records
- **readRecord.js**: Queries existing medical records

## Network Management

### View Logs
```bash
docker logs -f <container_name>
```

### Stop Network
```bash
docker-compose down
```

### Clean Up
```bash
docker-compose down -v
docker system prune -a --volumes
```

## Troubleshooting

### Common Issues

1. **Connection Refused Errors**
   - Verify IP addresses in docker-compose files
   - Check firewall rules between Pi devices
   - Ensure all containers are running: `docker ps`

2. **Endorsement Policy Failures**
   - Confirm chaincode is installed on all required peers
   - Verify organization MSP configurations

3. **Out of Memory**
   - Increase swap space on Raspberry Pi
   - Reduce concurrent workers in Caliper benchmarks


## Technology Stack

- **Blockchain**: Hyperledger Fabric 2.5
- **Smart Contracts**: Go (Golang)
- **Deployment**: Docker & Docker Compose
- **Benchmarking**: Hyperledger Caliper 0.5
- **Consensus**: Raft
- **Database**: GolangDB (embedded)

## License

This project is for educational and research purposes.

## References

- [Hyperledger Fabric Documentation](https://hyperledger-fabric.readthedocs.io/)
- [Hyperledger Caliper Documentation](https://hyperledger.github.io/caliper/)
- [Fabric Samples](https://github.com/hyperledger/fabric-samples)
