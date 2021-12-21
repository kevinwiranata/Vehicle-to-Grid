package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CSOContract contract for managing CRUD for EVs
type CSOContract struct {
	contractapi.Contract
}

// CreateCSOUser creates a new instance of cSO
func (c *CSOContract) CreateCSOUser(ctx contractapi.TransactionContextInterface, CSOID string) error {
	csoUser := new(CSO)
	csoUser.ID = CSOID
	exists, err := csoUser.LoadState(ctx)

	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The CSO %s already existss", CSOID)
	}

	newCSO := new(CSO)
	newCSO.ID = CSOID
	newCSO.EVCount = 0
	newCSO.TotalPowerCharge = 0
	newCSO.TotalPowerDischarge = 0

	return newCSO.SaveState(ctx)
}

// UpdateEVData retrieves an EV from the world state and updates its value
func (c *CSOContract) TransactEV(ctx contractapi.TransactionContextInterface, EVID string, CSOID string, PowerCharge float64, PowerDischarge float64, Temperature float64) ([]byte, error) {
	invokeArgs := [][]byte{[]byte("UpdateEVData"), []byte(EVID), []byte(CSOID), []byte(fmt.Sprint(PowerCharge)), []byte(fmt.Sprint(PowerDischarge)), []byte(fmt.Sprint(Temperature))}
	response := ctx.GetStub().InvokeChaincode("EV", invokeArgs, "mychannel")
	if response.Status != shim.OK {
		return nil, fmt.Errorf("Error invoking EV Chaincode. %s", response.GetMessage())
	}
	return response.GetPayload(), nil
}
