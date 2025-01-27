package rpc

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
)

const (
	ETHNetworkID = iota
	BSCNetworkID
	PLGNetworkID
	FTMNetworkID
)

// EVMIssuingMetadata keeps track of EVM issuing metadata types based on the EVM networkIDs.
var EVMIssuingMetadata = map[int]int{
	ETHNetworkID: metadata.IssuingETHRequestMeta,
	BSCNetworkID: metadata.IssuingBSCRequestMeta,
	PLGNetworkID: metadata.IssuingPLGRequestMeta,
	FTMNetworkID: metadata.IssuingFantomRequestMeta,
}

// EVMBurningMetadata keeps track of EVM burning metadata types based on the EVM networkIDs.
var EVMBurningMetadata = map[int]int{
	ETHNetworkID: metadata.BurningRequestMetaV2,
	BSCNetworkID: metadata.BurningPBSCRequestMeta,
	PLGNetworkID: metadata.BurningPLGRequestMeta,
	FTMNetworkID: metadata.BurningFantomRequestMeta,
}

var burnProofRPCMethod = map[int]string{
	ETHNetworkID: getBurnProof,
	BSCNetworkID: getBSCBurnProof,
	PLGNetworkID: getPLGBurnProof,
	FTMNetworkID: getFTMBurnProof,
}

// EVMNetworkNotFoundError returns an error indicating that the given EVM networkID is not supported.
func EVMNetworkNotFoundError(evmNetworkID int) error {
	return fmt.Errorf("EVMNetworkID %v not supported", evmNetworkID)
}

// GetBurnProof retrieves the burning proof of a transaction with the given target evmNetworkID.
// evmNetworkID can be one of the following:
//	- ETHNetworkID: the Ethereum network
//	- BSCNetworkID: the Binance Smart Chain network
//	- PLGNetworkID: the Polygon network
//	- FTMNetworkID: the Fantom network
// If set empty, evmNetworkID defaults to ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
func (server *RPCServer) GetBurnProof(txHash string, evmNetworkID ...int) ([]byte, error) {
	networkID := ETHNetworkID
	if len(evmNetworkID) > 0 {
		networkID = evmNetworkID[0]
	}

	if _, ok := burnProofRPCMethod[networkID]; !ok {
		return nil, EVMNetworkNotFoundError(networkID)
	}
	method := burnProofRPCMethod[networkID]
	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(method, params)
}

func (server *RPCServer) GetUnifiedBurnProof(txHash string) ([]byte, error) {
        tmpParams := make(map[string]interface{})
        tmpParams["TxReqID"] = txHash
	tmpParams["DataIndex"] = 0

	method := "bridgeaggGetBurnProof"
        params := make([]interface{}, 0)
        params = append(params, tmpParams)
        return server.SendQuery(method, params)
}

// GetBurnProofForSC retrieves the burning proof of a transaction for depositing to smart contracts.
func (server *RPCServer) GetBurnProofForSC(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(getBurnProofForDepositToSC, params)
}

// GetBurnPRVPeggingProof retrieves the burning prv pegging proof of a transaction.
func (server *RPCServer) GetBurnPRVPeggingProof(txHash string, evmNetworkIDs ...int) ([]byte, error) {
	method := getPRVERC20BurnProof
	if len(evmNetworkIDs) > 0 {
		switch evmNetworkIDs[0] {
		case BSCNetworkID:
			method = getPRVBEP20BurnProof
		case PLGNetworkID, FTMNetworkID:
			return nil, EVMNetworkNotFoundError(evmNetworkIDs[0])
		}
	}
	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(method, params)
}

// CheckShieldStatus checks the status of a decentralized shielding transaction.
func (server *RPCServer) CheckShieldStatus(txHash string) ([]byte, error) {
	tmpParams := make(map[string]interface{})
	tmpParams["TxReqID"] = txHash

	params := make([]interface{}, 0)
	params = append(params, tmpParams)
	return server.SendQuery(getBridgeReqWithStatus, params)
}

// GetAllBridgeTokens retrieves the list of bridge tokens in the network.
func (server *RPCServer) GetAllBridgeTokens() ([]byte, error) {
	return server.SendQuery(getAllBridgeTokens, nil)
}
