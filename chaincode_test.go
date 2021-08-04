/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type NFTMetadata struct { //NFT
	Name    string "json:\"name\""
	Hash    string "json:\"ipfshash\""
	Account string "json:\"account\""
}

func TestInit(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("chaincode", cc)
	res := stub.MockInit("1", [][]byte{[]byte("initFunc")})
	if res.Status != shim.OK {
		t.Error("Init failed", res.Status, res.Message)
	}
}
func balanceOf(stub *shim.MockStub) *shim.MockStub { //유저의 NFT 개수
	return stub

}

func safeTransferFrom(stub *shim.MockStub) *shim.MockStub { // 전달
	return stub

}

func checkOnNFTReceived(stub *shim.MockStub) *shim.MockStub { // NFT 전달 확인
	return stub

}

func publishNFT(stub *shim.MockStub) *shim.MockStub {
	// Make sure a valid operator is provided
	// Retrieve info needed for the update procedure
	txid := stub.GetTxID()
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("fcn: ", fcn)
	fmt.Println("params: ", params)
	compositeKeyName := params[0] //"Account"
	// 함수이름 , [0]유저 Account , [1] 파일 이름, [2] 파일의 해시값
	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey1, compositeErr := stub.CreateCompositeKey(compositeKeyName, []string{"NFT", params[0], params[1], params[2], txid})
	compositeKey2, compositeErr := stub.CreateCompositeKey(compositeKeyName, []string{"NFT", params[0], params[1], "FileHash2", txid})
	compositeKey3, compositeErr := stub.CreateCompositeKey(compositeKeyName, []string{"NFT", params[0], params[1], "FileHash3", txid})

	if compositeErr != nil {
		fmt.Println("Error : compositeErr " + compositeErr.Error())
	}
	fmt.Println("**** Create **** CreateCompositeKey1 is " + compositeKey1 + " ****")
	fmt.Println("**** Create **** CreateCompositeKey1 is " + compositeKey2 + " ****")

	// Save the composite key index
	stub.MockTransactionStart("testCompositeKey")
	compositePutErr := stub.PutState(compositeKey1, []byte{0x00})
	compositePutErr = stub.PutState(compositeKey2, []byte{0x00})
	compositePutErr = stub.PutState(compositeKey3, []byte{0x00})

	stub.MockTransactionEnd("testCompositeKey")

	if compositePutErr != nil {
		fmt.Println("Error : compositePutErr " + compositePutErr.Error())
	}
	return stub
}
func GetNFT(stub *shim.MockStub) *shim.MockStub { //저장된 인증서
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
	return stub
}
func TestInvoke(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("chaincode", cc)
	res := stub.MockInit("1", [][]byte{[]byte("initFunc")})
	if res.Status != shim.OK {
		t.Error("Init failed", res.Status, res.Message)
	}
	res = stub.MockInvoke("1", [][]byte{[]byte("publishNFT"), []byte("Account"), []byte("FileName"), []byte("FileHash1")})
	startTime := time.Now()

	publishNFT(stub)
	elapsedTime := time.Since(startTime)

	fmt.Printf("Invoke 실행시간: %f\n", elapsedTime.Seconds())

	res = stub.MockInvoke("1", [][]byte{[]byte("VerifyGenesisBlockByCompositeKey"), []byte("Account"), []byte("NFT")})
	if res.Status != shim.OK {
		t.Error("Invoke failed", res.Status, res.Message)
	}

	GetNFT(stub)

}
