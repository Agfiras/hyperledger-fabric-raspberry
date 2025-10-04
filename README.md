```markdown
# Multi-Host Hyperledger Fabric Network on Raspberry Pi

A distributed Hyperledger Fabric 2.5 blockchain network running across two Raspberry Pi 3 Model B devices.

## Network Architecture

- **Pi1 (Raspberry Pi 1)**: Orderer + Peer0.Org1 + CA.Org1
- **Pi2 (Raspberry Pi 2)**: Peer0.Org2 + CA.Org2
- **Connectivity**: Tailscale VPN (100.x.x.x addresses)
- **Consensus**: etcdraft
- **TLS**: Enabled

## Prerequisites

- 2x Raspberry Pi 3 Model B (or higher)
- Docker and Docker Compose
- Hyperledger Fabric 2.5 binaries
- Tailscale VPN (or local network with proper DNS)

## Network Components

### Pi1 (172.20.10.11 / 100.69.213.86)
- `orderer1.example.com:7050` - Ordering service
- `peer0.org1.example.com:7051` - Org1 peer
- `ca.org1.example.com:7054` - Org1 Certificate Authority

### Pi2 (172.20.10.12 / 100.81.64.92)
- `peer0.org2.example.com:8051` - Org2 peer
- `ca.org2.example.com:8054` - Org2 Certificate Authority

## Setup Instructions

### 1. Generate Cryptographic Material

cd ~/fabric/fabric-network
cryptogen generate --config=./crypto-config.yaml

```

### 2. Generate Channel Configuration

```

configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channelc.tx -channelID channelc
configtxgen -profile TwoOrgsChannel -outputBlock ./channel-artifacts/channelc.block -channelID channelc

```

### 3. Copy Files to Pi2

```


# From Pi1

scp -r crypto-config pi@<pi2-ip>:~/fabric/fabric-network/
scp -r channel-artifacts pi@<pi2-ip>:~/fabric/fabric-network/

```

### 4. Start Network

**On Pi1:**
```

cd ~/fabric/fabric-network
sudo docker-compose -f docker-compose-pi1.yaml up -d

```

**On Pi2:**
```

cd ~/fabric/fabric-network
sudo docker-compose -f docker-compose-pi2.yaml up -d

```

### 5. Create and Join Channel

**Create channel (Pi1):**
```

sudo docker exec -it fabric-network-peer0.org1.example.com-1 bash
peer channel create -o orderer1.example.com:7050 -c channelc \
-f /etc/hyperledger/crypto/channel-artifacts/channelc.tx \
--outputBlock /etc/hyperledger/crypto/channel-artifacts/channelc.block \
--tls true \
--cafile /etc/hyperledger/crypto/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/ca.crt

```

**Join peers to channel:**
```


# Pi1 - Org1 peer

peer channel join -b /etc/hyperledger/crypto/channel-artifacts/channelc.block

# Pi2 - Org2 peer (in peer container)

peer channel join -b /etc/hyperledger/crypto/channel-artifacts/channelc.block

```

## Configuration Files

- `configtx.yaml` - Channel and network configuration
- `crypto-config.yaml` - Certificate generation configuration
- `docker-compose-pi1.yaml` - Pi1 container definitions
- `docker-compose-pi2.yaml` - Pi2 container definitions (if separate)
- `core.yaml` - Peer configuration (optional)

## Network DNS Configuration

Update `/etc/hosts` on both Pi devices:

**Pi1:**
```

127.0.0.1       orderer1.example.com peer0.org1.example.com ca.org1.example.com
100.81.64.92    peer0.org2.example.com ca.org2.example.com

```

**Pi2:**
```

100.69.213.86   orderer1.example.com peer0.org1.example.com ca.org1.example.com
127.0.0.1       peer0.org2.example.com ca.org2.example.com

```

## Verification

```


# Check channel membership

peer channel list

# Check channel info

peer channel getinfo -c channelc

# Check gossip connectivity

sudo docker logs fabric-network-peer0.org2.example.com-1 | grep "Membership view"

```

## Troubleshooting

- **Gossip timeout errors**: Check /etc/hosts DNS entries
- **TLS certificate errors**: Regenerate crypto-config
- **Port conflicts**: Ensure ports 7050, 7051, 8051, 9444-9449 are available
```