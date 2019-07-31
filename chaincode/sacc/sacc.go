// 외부 모듈 추가
package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 체인코드 구조체 포함
type Asset struct {

}

// ChaincodeStubInterface: PutState, GetState

// 1. Init 함수 구현
// Instantiate(배포), Upgrade시 호출
// Peer의 컨테이너에 ChainCode가 저장되고 Channel에 그 ChainCode가 Install되면
// orderer에 의해 Channel에 배포되고 Channel에 있는 다른 Peer들의 컨테이너에
// ChainCode가 저장된다.
func (t *Asset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("error");
	}
	
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

// 2. Invoke 함수 구현
// submitTx, evaluateTx 호출 시
// 
func (t *Asset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "get" {
		result, err = get(stub, args)
	} else if fn == "getAllKeys" {
		result, err = getAllKeys(stub)
	} else {
		return shim.Error("Not supported chaincode function.")
	}

	if err != nil {
		return shim.Error(err.Error())
	}
	// endorsor로 보내고, endorsor가 application(or cli)로 보냄
	return shim.Success([]byte(result))
}

// set 함수 선언

func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value");
	}
	
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expection a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

func getAllKeys(stub shim.ChaincodeStubInterface) (string, error) {

	iter, err := stub.GetStateByRange("a", "z")
	if err != nil {
		return "", fmt.Errorf("Failed to get all keys with error: %s", err)
	}
	defer iter.Close()

	var buffer string
	buffer = "["
	comma := false
	for iter.HasNext() {
		res, err := iter.Next()
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}
		if comma == true {
			buffer += ","
		}
		buffer += "{\"key\":"
		buffer += "\""
		buffer += res.Key
		buffer += "\", \"Value\":\""
		buffer += string(res.Value)
		buffer += "\"}"
		comma = true
	}
	buffer += "]"
	fmt.Println(buffer)
	return string(buffer), nil
}

// 메인
func main() {
	if err := shim.Start(new(Asset)); err != nil {
		fmt.Printf("Error starting Asset chaincode: %s", err)
	}
}