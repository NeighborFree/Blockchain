/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
}

func (cc *Chaincode) publishNFT(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	// Make sure a valid operator is provided
	// Retrieve info needed for the update procedure
	startTime := time.Now()

	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("fcn: ", fcn)
	fmt.Println("params: ", params)
	compositeKeyName := params[0] //"Account"
	// 함수이름 , [0]유저 Account , [1] 파일 이름, [2] 파일의 해시값
	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey1, compositeErr := stub.CreateCompositeKey(compositeKeyName, []string{"NFT", params[0], params[1], params[2]})

	if compositeErr != nil {
		fmt.Println("Error : compositeErr " + compositeErr.Error())
	}
	fmt.Println("**** Create **** CreateCompositeKey1 is " + compositeKey1 + " ****")

	// Save the composite key index
	compositePutErr := stub.PutState(compositeKey1, []byte{0x00})

	if compositePutErr != nil {
		fmt.Println("Error : compositePutErr " + compositePutErr.Error())
	}
	elapsedTime := time.Since(startTime)

	fmt.Printf("PUT 실행시간: %f\n", elapsedTime.Seconds())
	return shim.Success([]byte("publishNFT"))
}
func (cc *Chaincode) GetNFT(stub shim.ChaincodeStubInterface, params []string) sc.Response { //저장된 인증서
	startTime := time.Now()

	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("fcn: ", fcn)
	fmt.Println("params: ", params)

	//Unmarshal Sensorchain Metadata

	approvalIterator, err := stub.GetStateByPartialCompositeKey(params[0], []string{params[1]}) //SensorChainNameInit
	if err != nil {

	}
	defer approvalIterator.Close()

	for approvalIterator.HasNext() {
		responseRange, nextErr := approvalIterator.Next()
		if nextErr != nil {
		}
		fmt.Println("Current " + responseRange.Key)
		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)

		// Split the composite key into its component parts
		if splitKeyErr != nil {
		}
		fmt.Println("result " + keyParts[0] + " " + keyParts[1] + " " + keyParts[2] + " " + keyParts[3])

	}
	elapsedTime := time.Since(startTime)

	fmt.Printf("GET 실행시간: %f\n", elapsedTime.Seconds())
	return shim.Success([]byte("GetNFT"))
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Init()", fcn, params)
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Invoke()", fcn, params)
	switch fcn {
	case "setEnrollmentSensorChain":
		return cc.publishNFT(stub, params)

	case "getSensorChainByCompositeKey":
		return cc.GetNFT(stub, params)

	default:
		return sc.Response{Status: 404, Message: "404 Not Found", Payload: nil}
	}

}
