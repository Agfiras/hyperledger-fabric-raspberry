package main

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing medical records
type SmartContract struct {
    contractapi.Contract
}

// MedicalRecord describes basic details of a medical record
type MedicalRecord struct {
    RecordID    string `json:"recordID"`
    PatientName string `json:"patientName"`
    Diagnosis   string `json:"diagnosis"`
    Treatment   string `json:"treatment"`
    Timestamp   string `json:"timestamp"`
}

// InitLedger adds a base set of medical records to the ledger (optional)
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    return nil
}

// CreateRecord issues a new medical record to the world state with given details
func (s *SmartContract) CreateRecord(ctx contractapi.TransactionContextInterface, recordID string, patientName string, diagnosis string, treatment string, timestamp string) error {
    exists, err := s.RecordExists(ctx, recordID)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("the record %s already exists", recordID)
    }

    record := MedicalRecord{
        RecordID:    recordID,
        PatientName: patientName,
        Diagnosis:   diagnosis,
        Treatment:   treatment,
        Timestamp:   timestamp,
    }
    recordJSON, err := json.Marshal(record)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(recordID, recordJSON)
}

// ReadRecord returns the medical record stored in the world state with given id
func (s *SmartContract) ReadRecord(ctx contractapi.TransactionContextInterface, recordID string) (*MedicalRecord, error) {
    recordJSON, err := ctx.GetStub().GetState(recordID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if recordJSON == nil {
        return nil, fmt.Errorf("the record %s does not exist", recordID)
    }

    var record MedicalRecord
    err = json.Unmarshal(recordJSON, &record)
    if err != nil {
        return nil, err
    }

    return &record, nil
}

// UpdateRecord updates an existing medical record in the world state with provided parameters
func (s *SmartContract) UpdateRecord(ctx contractapi.TransactionContextInterface, recordID string, patientName string, diagnosis string, treatment string, timestamp string) error {
    exists, err := s.RecordExists(ctx, recordID)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the record %s does not exist", recordID)
    }

    // overwriting original record with new record
    record := MedicalRecord{
        RecordID:    recordID,
        PatientName: patientName,
        Diagnosis:   diagnosis,
        Treatment:   treatment,
        Timestamp:   timestamp,
    }
    recordJSON, err := json.Marshal(record)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(recordID, recordJSON)
}

// DeleteRecord deletes a given medical record from the world state
func (s *SmartContract) DeleteRecord(ctx contractapi.TransactionContextInterface, recordID string) error {
    exists, err := s.RecordExists(ctx, recordID)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the record %s does not exist", recordID)
    }

    return ctx.GetStub().DelState(recordID)
}

// RecordExists returns true when record with given ID exists in world state
func (s *SmartContract) RecordExists(ctx contractapi.TransactionContextInterface, recordID string) (bool, error) {
    recordJSON, err := ctx.GetStub().GetState(recordID)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return recordJSON != nil, nil
}

// GetAllRecords returns all medical records found in world state
func (s *SmartContract) GetAllRecords(ctx contractapi.TransactionContextInterface) ([]*MedicalRecord, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var records []*MedicalRecord
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var record MedicalRecord
        err = json.Unmarshal(queryResponse.Value, &record)
        if err != nil {
            return nil, err
        }
        records = append(records, &record)
    }

    return records, nil
}
