package evm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/version"
	"github.com/cosmos/ethermint/x/evm/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
)

// Supported endpoints
const (
	QueryBalance = "balance"
	QueryBlockNumber = "blockNumber"
	QueryStorage = "storage"
	QueryCode = "code"
)

// TODO: Implement querier to route RPC methods. Unable to access RPC otherwise
// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, path, keeper)
		case QueryBlockNumber:
			return queryBlockNumber(ctx, keeper)
		case QueryStorage:
			return queryStorage(ctx, path, keeper)
		case QueryCode:
			return queryCode(ctx, path, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint")
		}
	}
}

func queryProtocolVersion(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	vers := version.ProtocolVersion

	res, err := codec.MarshalJSONIndent(keeper.cdc, vers)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

//func querySyncing(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper ) ([]byte, error) {
//	// TODO: Implement
//	status := true
//
//	res, err := codec.MarshalJSONIndent(keeper.cdc, status)
//	if err != nil {
//		panic("could not marshal result to JSON")
//	}
//
//	return res, nil
//}
//
//func queryCoinbase(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper ) ([]byte, error) {
//	// TODO: Implement
//	status := true
//
//	res, err := codec.MarshalJSONIndent(keeper.cdc, status)
//	if err != nil {
//		panic("could not marshal result to JSON")
//	}
//
//	return res, nil
//}

func queryBalance(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := ethcmn.BytesToAddress([]byte(path[1]))
	balance := keeper.GetBalance(ctx, addr)
	bRes := types.QueryResBalance{ Balance: balance}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryBlockNumber(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	num := ctx.BlockHeight()
	bnRes := types.QueryResBlockNumber{ Number: big.NewInt(num)}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bnRes)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryStorage(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := ethcmn.BytesToAddress([]byte(path[1]))
	key := ethcmn.BytesToHash([]byte(path[2]))
	val := keeper.GetState(ctx, addr, key)
	bRes := types.QueryResStorage{ Value: val}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryCode(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := ethcmn.BytesToAddress([]byte(path[1]))
	code := keeper.GetCode(ctx, addr)
	cRes := types.QueryResCode{ Code: code}
	res, err := codec.MarshalJSONIndent(keeper.cdc, cRes)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}