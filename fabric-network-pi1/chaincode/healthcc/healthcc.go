package chaincode

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing patient containers
type SmartContract struct {
    contractapi.Contract
}

// PatientContainer describes details of the asset
// Fields in alphabetical order for determinism across languages
type PatientContainer struct {
    CID string `json:"CID"`
    UID string `json:"UID"`
}

// CreateOrUpdateContainer issues a new container or updates an existing container in the world state with given details.
func (s *SmartContract) CreateOrUpdateContainer(ctx contractapi.TransactionContextInterface, uid string, cid string) error {
    asset := PatientContainer{CID: cid, UID: uid}
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }
    return ctx.GetStub().PutState(uid, assetJSON)
}

// ReadContainer returns the container stored in the world state with given id.
func (s *SmartContract) ReadContainer(ctx contractapi.TransactionContextInterface, uid string) (*PatientContainer, error) {
    assetJSON, err := ctx.GetStub().GetState(uid)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("the asset %s does not exist", uid)
    }
    var asset PatientContainer
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, err
    }
    return &asset, nil
}

// DeleteContainer deletes a given container from the world state.
func (s *SmartContract) DeleteContainer(ctx contractapi.TransactionContextInterface, uid string) error {
    exists, err := s.ContainerExists(ctx, uid)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", uid)
    }
    return ctx.GetStub().DelState(uid)
}

// ContainerExists returns true when container with given ID exists in world state.
func (s *SmartContract) ContainerExists(ctx contractapi.TransactionContextInterface, uid string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(uid)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }
    return assetJSON != nil, nil
}

// ReadContainerHistory returns all versions of a container stored in the world state with given ID.
func (s *SmartContract) ReadContainerHistory(ctx contractapi.TransactionContextInterface, uid string) ([]PatientContainer, error) {
    historyIterator, err := ctx.GetStub().GetHistoryForKey(uid)
    if err != nil {
        return nil, err
    }
    if historyIterator == nil {
        return nil, fmt.Errorf("the asset %s does not exist", uid)
    }
    defer historyIterator.Close()
    var assets []PatientContainer
    for historyIterator.HasNext() {
        modification, err := historyIterator.Next()
        if err != nil {
            return nil, err
        }
        fmt.Println("Returning information about", string(modification.Value))
        var asset PatientContainer
        err = json.Unmarshal(modification.Value, &asset)
        if err != nil {
            return nil, err
        }
        assets = append(assets, asset)
    }
    return assets, nil
}

// GetAllContainers returns all containers found in world state.
func (s *SmartContract) GetAllContainers(ctx contractapi.TransactionContextInterface) ([]PatientContainer, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()
    var assets []PatientContainer
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }
        var asset PatientContainer
        err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
            return nil, err
        }
        assets = append(assets, asset)
    }
    return assets, nil
}

// GetAllContainersHistory returns all versions of all containers found in world state.
func (s *SmartContract) GetAllContainersHistory(ctx contractapi.TransactionContextInterface) ([]PatientContainer, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()
    var assets []PatientContainer
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }
        var asset PatientContainer
        err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
            return nil, err
        }
        assetid := asset.UID
        historyIterator, err := ctx.GetStub().GetHistoryForKey(assetid)
        if err != nil {
            return nil, err
        }
        defer historyIterator.Close()
        var assethistory []PatientContainer
        for historyIterator.HasNext() {
            modification, err := historyIterator.Next()
            if err != nil {
                return nil, err
            }
            var assetversion PatientContainer
            err = json.Unmarshal(modification.Value, &assetversion)
            if err != nil {
                return nil, err
            }
            assethistory = append(assethistory, assetversion)
        }
        assets = append(assets, assethistory...)
    }
    return assets, nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        fmt.Printf("Error create chaincode: %s", err)
        return
    }
    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting chaincode: %s", err)
    }
}