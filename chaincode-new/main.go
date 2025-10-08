package main

import (
    "fmt"
    "log"
    
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
    assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
        log.Panicf("Error creating patientcc chaincode: %v", err)
    }

    if err := assetChaincode.Start(); err != nil {
        fmt.Printf("Error starting patientcc chaincode: %v", err)
    }
}
