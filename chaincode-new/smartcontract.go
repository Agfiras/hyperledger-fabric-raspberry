package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing assets
type SmartContract struct {
	contractapi.Contract
}

// PatientContainer describes details of what makes up an asset
//Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type PatientContainer struct {
	CID             string `json:"CID"`
	UID             string `json:"UID"`
}

func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// CreateOrUpdateContainer issues a new container or updates an existing container in the world state with given details.
func (s *SmartContract) CreateOrUpdateContainer(ctx contractapi.TransactionContextInterface, uid string, cid string) error {

	asset := PatientContainer{
		CID:             cid,
		UID:             uid,
	}
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
func (s *SmartContract) ReadContainerHistory(ctx contractapi.TransactionContextInterface, uid string) ([]*PatientContainer, error) {
	historyIterator, err := ctx.GetStub().GetHistoryForKey(uid)
	if err != nil {
		return nil, err
	}
	if historyIterator == nil {
		return nil, fmt.Errorf("the asset %s does not exist", uid)
	}
	defer historyIterator.Close()
	
	var assets []*PatientContainer
	for historyIterator.HasNext() {
		modification, err := historyIterator.Next()
		if err != nil {
			return nil, err
		}
		//fmt.Println("Returning information about", string(modification.Value))
		
		var asset PatientContainer
		err = json.Unmarshal(modification.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}
	
	return assets, nil
}

// GetAllAssets returns all containers found in world state.
func (s *SmartContract) GetAllContainers(ctx contractapi.TransactionContextInterface) ([]*PatientContainer, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*PatientContainer
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
		assets = append(assets, &asset)
	}

	return assets, nil
}

// GetAllContainersHistory returns all versions of all containers found in world state.
func (s *SmartContract) GetAllContainersHistory(ctx contractapi.TransactionContextInterface) ([][]*PatientContainer, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets [][]*PatientContainer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		//an asset instance (asset) storing all results iterator contents (assets in the world state) one by one
		var asset PatientContainer
		//copy iterator next to asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		//id of current asset so that we can get its history
		assetid := asset.UID
		historyIterator, err := ctx.GetStub().GetHistoryForKey(assetid)
		if err != nil {
			return nil, err
		}
		defer historyIterator.Close()
		
		var assethistory []*PatientContainer
		for historyIterator.HasNext() {
			modification, err := historyIterator.Next()
			if err != nil {
				return nil, err
			}
			//an asset instance (assetversion) storing all history iterator contents (versions of one asset) one by one
			var assetversion PatientContainer
			//copy iterator next to assetversion
			err = json.Unmarshal(modification.Value, &assetversion)
			if err != nil {
				return nil, err
			}
			assethistory = append(assethistory, &assetversion)
		}		
		assets = append(assets, assethistory)	
	}

	return assets, nil
}