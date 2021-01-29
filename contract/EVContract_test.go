// /*
//  * SPDX-License-Identifier: Apache-2.0
//  */

package main

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"testing"

// 	"github.com/hyperledger/fabric-contract-api-go/contractapi"
// 	"github.com/hyperledger/fabric-chaincode-go/shim"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// const getStateError = "world state get error"

// type MockStub struct {
// 	shim.ChaincodeStubInterface
// 	mock.Mock
// }

// func (ms *MockStub) GetState(key string) ([]byte, error) {
// 	args := ms.Called(key)

// 	return args.Get(0).([]byte), args.Error(1)
// }

// func (ms *MockStub) PutState(key string, value []byte) error {
// 	args := ms.Called(key, value)

// 	return args.Error(0)
// }

// func (ms *MockStub) DelState(key string) error {
// 	args := ms.Called(key)

// 	return args.Error(0)
// }

// type MockContext struct {
// 	contractapi.TransactionContextInterface
// 	mock.Mock
// }

// func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
// 	args := mc.Called()

// 	return args.Get(0).(*MockStub)
// }

// func configureStub() (*MockContext, *MockStub) {
// 	var nilBytes []byte

// 	testMyAsset := new(MyAsset)
// 	testMyAsset.Value = "set value"
// 	myAssetBytes, _ := json.Marshal(testMyAsset)

// 	ms := new(MockStub)
// 	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
// 	ms.On("GetState", "missingkey").Return(nilBytes, nil)
// 	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
// 	ms.On("GetState", "myAssetkey").Return(myAssetBytes, nil)
// 	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
// 	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

// 	mc := new(MockContext)
// 	mc.On("GetStub").Return(ms)

// 	return mc, ms
// }

// func TestMyAssetExists(t *testing.T) {
// 	var exists bool
// 	var err error

// 	ctx, _ := configureStub()
// 	c := new(MyAssetContract)

// 	exists, err = c.MyAssetExists(ctx, "statebad")
// 	assert.EqualError(t, err, getStateError)
// 	assert.False(t, exists, "should return false on error")

// 	exists, err = c.MyAssetExists(ctx, "missingkey")
// 	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
// 	assert.False(t, exists, "should return false when no value for key in world state")

// 	exists, err = c.MyAssetExists(ctx, "existingkey")
// 	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
// 	assert.True(t, exists, "should return true when value for key in world state")
// }

// func TestCreateMyAsset(t *testing.T) {
// 	var err error

// 	ctx, stub := configureStub()
// 	c := new(MyAssetContract)

// 	err = c.CreateMyAsset(ctx, "statebad", "some value")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

// 	err = c.CreateMyAsset(ctx, "existingkey", "some value")
// 	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

// 	err = c.CreateMyAsset(ctx, "missingkey", "some value")
// 	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
// }

// func TestReadMyAsset(t *testing.T) {
// 	var myAsset *MyAsset
// 	var err error

// 	ctx, _ := configureStub()
// 	c := new(MyAssetContract)

// 	myAsset, err = c.ReadMyAsset(ctx, "statebad")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
// 	assert.Nil(t, myAsset, "should not return MyAsset when exists errors when reading")

// 	myAsset, err = c.ReadMyAsset(ctx, "missingkey")
// 	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
// 	assert.Nil(t, myAsset, "should not return MyAsset when key does not exist in world state when reading")

// 	myAsset, err = c.ReadMyAsset(ctx, "existingkey")
// 	assert.EqualError(t, err, "Could not unmarshal world state data to type MyAsset", "should error when data in key is not MyAsset")
// 	assert.Nil(t, myAsset, "should not return MyAsset when data in key is not of type MyAsset")

// 	myAsset, err = c.ReadMyAsset(ctx, "myAssetkey")
// 	expectedMyAsset := new(MyAsset)
// 	expectedMyAsset.Value = "set value"
// 	assert.Nil(t, err, "should not return error when MyAsset exists in world state when reading")
// 	assert.Equal(t, expectedMyAsset, myAsset, "should return deserialized MyAsset from world state")
// }

// func TestUpdateMyAsset(t *testing.T) {
// 	var err error

// 	ctx, stub := configureStub()
// 	c := new(MyAssetContract)

// 	err = c.UpdateMyAsset(ctx, "statebad", "new value")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

// 	err = c.UpdateMyAsset(ctx, "missingkey", "new value")
// 	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

// 	err = c.UpdateMyAsset(ctx, "myAssetkey", "new value")
// 	expectedMyAsset := new(MyAsset)
// 	expectedMyAsset.Value = "new value"
// 	expectedMyAssetBytes, _ := json.Marshal(expectedMyAsset)
// 	assert.Nil(t, err, "should not return error when MyAsset exists in world state when updating")
// 	stub.AssertCalled(t, "PutState", "myAssetkey", expectedMyAssetBytes)
// }

// func TestDeleteMyAsset(t *testing.T) {
// 	var err error

// 	ctx, stub := configureStub()
// 	c := new(MyAssetContract)

// 	err = c.DeleteMyAsset(ctx, "statebad")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

// 	err = c.DeleteMyAsset(ctx, "missingkey")
// 	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

// 	err = c.DeleteMyAsset(ctx, "myAssetkey")
// 	assert.Nil(t, err, "should not return error when MyAsset exists in world state when deleting")
// 	stub.AssertCalled(t, "DelState", "myAssetkey")
// }
