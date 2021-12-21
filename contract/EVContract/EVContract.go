/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// EVContract contract for managing CRUD for EVs
type EVContract struct {
	contractapi.Contract
}

// InitLedger creates the initial set of EVs in the ledger.
func (c *EVContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	EVUsers := []EV{
		{ID: "ID1", Model: "Tesla", Age: 3},
		{ID: "ID2", Model: "Tesla", Age: 1},
		{ID: "ID3", Model: "TEsla", Age: 36},
		{ID: "ID4", Model: "BMW", Age: 9},
		{ID: "ID5", Model: "Mercedes", Age: 12},
		{ID: "ID6", Model: "Aston", Age: 3},
		{ID: "ID7", Model: "Renault", Age: 3},
	}

	for _, User := range EVUsers {
		err := c.CreateEVUser(ctx, User.ID, User.Model, User.Age, "CSO1")
		if err != nil {
			return err
		}
	}

	return nil
}

//EVUserExists checks if a given EV exists
func (c *EVContract) EVUserExists(ctx contractapi.TransactionContextInterface, ID string) (bool, error) {
	evUser := EV{ID: ID}
	exists, err := evUser.LoadState(ctx)
	if err != nil {
		return false, fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return true, nil
	} else {
		return false, nil
	}
}

// CreateEVUser creates a new instance of EV
func (c *EVContract) CreateEVUser(ctx contractapi.TransactionContextInterface, EVID string, model string, age int, CSOID string) error {
	evUser := new(EV)
	evUser.ID = EVID
	exists, err := evUser.LoadState(ctx)

	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The EV %s already exists", EVID)
	}

	newEV := new(EV)
	newEV.Model = model
	newEV.ID = EVID
	newEV.Age = age

	// csoUser := new(CSO)
	// csoUser.ID = CSOID
	// exists, err = csoUser.LoadState(ctx)
	// if err != nil {
	// 	return fmt.Errorf("Could not read from world state. %s", err)
	// } else if !exists {
	// 	return fmt.Errorf("The given CSO %s does not exist", CSOID)
	// }

	// newEVData := new(EVData)
	// newEVData.ID = EVID
	// newEVData.Model = model
	// csoUser.EVDatas = append(csoUser.EVDatas, *newEVData)
	// err = csoUser.SaveState(ctx)
	// if err != nil {
	// 	return fmt.Errorf("Could not save to world state!. %s", err)
	// }

	return newEV.SaveState(ctx)
}

// ReadEVData retrieves an instance of EV from the world state
func (c *EVContract) ReadEVData(ctx contractapi.TransactionContextInterface, ID string) (*EV, error) {
	evUser := new(EV)
	evUser.ID = ID
	exists, err := evUser.LoadState(ctx)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The EV User %s does not exist", ID)
	}

	return evUser, nil
}

// DeleteEVUser deletes an EV from the world state
func (c *EVContract) DeleteEVUser(ctx contractapi.TransactionContextInterface, ID string) error {
	evUser := new(EV)
	evUser.ID = ID
	exists, err := evUser.LoadState(ctx)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The EV User %s does not exist", ID)
	}

	return ctx.GetStub().DelState(ID)
}

// UpdateEVData is only called by the CSOContract (Invoke Chaincode)
func (c *EVContract) UpdateEVData(ctx contractapi.TransactionContextInterface, EVID string, CSOID string, PowerCharge float64, PowerDischarge float64, Temperature float64) error {
	fmt.Println(EVID, CSOID)
	evUser := new(EV)
	evUser.ID = EVID
	exists, err := evUser.LoadState(ctx)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The EV User %s does not exist", EVID)
	}

	evUser.PowerCharge = PowerCharge
	evUser.PowerDischarge = PowerDischarge
	evUser.Temperature = Temperature

	return evUser.SaveState(ctx)
}

// QueryAll returns a JSON of all the EVs on the blockchain
func (c *EVContract) QueryAll(ctx contractapi.TransactionContextInterface) ([]*EV, error) {
	it, err := ctx.GetStub().GetStateByPartialCompositeKey(EVObjectType, []string{})
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	}
	defer it.Close()

	i := 1
	var EVData []*EV
	for it.HasNext() {
		fmt.Printf("Ran, %v time", i)
		response, err := it.Next()
		if err != nil {
			return nil, fmt.Errorf("unable to get the next element: %s", err.Error())
		}

		var singleEV EV
		if err := json.Unmarshal(response.Value, &singleEV); err != nil {
			return nil, fmt.Errorf("unable to parse the response: %s", err.Error())
		}
		EVData = append(EVData, &singleEV)
		i++
	}
	return EVData, nil
}

//QueryByFields allows users to query with optional fields
//Performs a rich query on couchDB with indexing (paramters can be tuned in the future)
//Parameters and the selectors can be tuend accordingly:
// model -> model of car
// age   -> age of car
// op    -> operator for comparison (i.e. $eq, $gt, $gte, $lt, $lte)
// QueryByFields(Tesla, $gt, 2) will return all Teslas that have age greater than 2
func (c *EVContract) QueryByFields(ctx contractapi.TransactionContextInterface, model string, op string, age int) ([]*EV, error) {
	queryString := fmt.Sprintf(`{"selector":{"model": "%s", "age": {"%s": %v}}, "use_index": ["_design/indexEVDoc", "indexEV"]}`, model, op, age)
	//queryString := fmt.Sprintf(`{"selector":{"model":"%s", "age": %v}}`, model, age)
	println(queryString)
	it, err := ctx.GetStub().GetQueryResult(queryString)
	println("iterator:")
	println(it)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	}
	defer it.Close()

	var EVS []*EV
	for it.HasNext() {
		queryResult, err := it.Next()
		if err != nil {
			return nil, err
		}
		var ev EV
		err = json.Unmarshal(queryResult.Value, &ev)
		if err != nil {
			return nil, err
		}
		EVS = append(EVS, &ev)
	}

	return EVS, nil
}
