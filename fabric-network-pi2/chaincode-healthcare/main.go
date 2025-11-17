package main

import (
    "fmt"
    "log"
    
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
    assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
        log.Panicf("Error creating healthcare chaincode: %v", err)
    }

    if err := assetChaincode.Start(); err != nil {
        fmt.Printf("Error starting healthcare chaincode: %v", err)
    }
}
