package arbtest

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/offchainlabs/nitro/solgen/go/precompilesgen"
)

func TestPriceOracleUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the node with default configuration
	builder := NewNodeBuilder(ctx).DefaultConfig(t, true)
	cleanup := builder.Build(t)
	defer cleanup()

	// Initialize the ArbPriceFeed precompile
	arbPriceFeed, err := precompilesgen.NewArbPriceFeed(common.HexToAddress("0x11a"), builder.L2.Client)
	Require(t, err)

	// Generate the User account
	builder.L2Info.GenerateAccount("User")
	// Fund the User account
	builder.L2.TransferBalance(t, "Faucet", "User", common.Big1, builder.L2Info)

	// Send a dummy transaction to trigger the pre-block hook
	userAuth := builder.L2Info.GetDefaultTransactOpts("User", ctx)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    userAuth.Nonce.Uint64(),
		GasPrice: userAuth.GasPrice,
		Gas:      userAuth.GasLimit,
		To:       &common.Address{},
		Value:    big.NewInt(0),
		Data:     nil,
	})
	builder.L2.Client.SendTransaction(ctx, tx)

	// Wait for the transaction to be mined
	receipt, err := WaitForTx(ctx, builder.L2.Client, tx.Hash(), time.Second*10)
	Require(t, err)
	fmt.Println("Transaction mined in block:", receipt.BlockNumber)

	// Now try to retrieve the BTC price
	callOpts := &bind.CallOpts{Context: ctx, BlockNumber: receipt.BlockNumber}
	price, err := arbPriceFeed.GetLatestBtcPrice(callOpts)
	Require(t, err)
	if price == 0 {
		t.Fatal("Expected BTC price to be modified")
	}
}
