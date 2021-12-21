/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// EVObjectType for composite key
const EVObjectType = "EV-Owner"

// EV stores data
type EV struct {
	ID             string  `json:"ID"`
	Model          string  `json:"model"`
	PowerCharge    float64 `json:"power_charge"`
	PowerDischarge float64 `json:"power_discharge"`
	Age            int     `json:"age"`
	Temperature    float64 `json:"temperature"`
	Energy         float64 `json:"energy"`
}

// ToCompositeKey returns a composite key based on the ID and accountType
func (c *EV) ToCompositeKey(ctx contractapi.TransactionContextInterface) (string,
	error) {
	attributes := []string{
		c.ID,
	}
	return ctx.GetStub().CreateCompositeKey(EVObjectType, attributes)
}

// ToLedgerValue creates a JSON-encoded account
func (c *EV) ToLedgerValue() ([]byte, error) {
	return json.Marshal(c)
}

// SaveState saves the accounts into the ledger
func (c *EV) SaveState(ctx contractapi.TransactionContextInterface) error {
	compositeKey, err := c.ToCompositeKey(ctx)
	if err != nil {
		message := fmt.Sprintf("Unable to create a composite key: %s", err.Error())
		return errors.New(message)
	}

	ledgerValue, err := c.ToLedgerValue()

	if err != nil {
		message := fmt.Sprintf("Unable to  compose a ledger value: %s", err.Error())
		return errors.New(message)
	}
	return ctx.GetStub().PutState(compositeKey, ledgerValue)
}

// LoadState loads the data from the ledger into the EV object if the data is found
// Returns false if an Account object wasn't found in the ledger; otherwise
//returns true.
func (c *EV) LoadState(ctx contractapi.TransactionContextInterface) (bool, error) {
	compositeKey, err := c.ToCompositeKey(ctx)
	if err != nil {
		message := fmt.Sprintf("Unable to create a composite key: %s", err.Error())
		return false, errors.New(message)
	}

	ledgerValue, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		message := fmt.Sprintf("Unable to  compose a ledger value: %s", err.Error())
		return false, errors.New(message)
	}

	if ledgerValue == nil {
		return false, nil
	}

	return true, json.Unmarshal(ledgerValue, &c)
}
