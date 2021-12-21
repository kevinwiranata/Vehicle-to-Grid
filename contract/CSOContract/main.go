/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	myAssetContract := new(CSOContract)
	myAssetContract.Info.Version = "0.0.1"
	myAssetContract.Info.Description = "My Smart Contract"
	myAssetContract.Info.License = new(metadata.LicenseMetadata)
	myAssetContract.Info.License.Name = "Apache-2.0"
	myAssetContract.Info.Contact = new(metadata.ContactMetadata)
	myAssetContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(myAssetContract)
	chaincode.Info.Title = "Test2 chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from MyAssetContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
