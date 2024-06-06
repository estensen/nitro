package arbtest

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/offchainlabs/nitro/solgen/go/precompilesgen"
)

func TestPriceOracleUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	builder := NewNodeBuilder(ctx).DefaultConfig(t, true)
	cleanup := builder.Build(t)
	defer cleanup()

	arbHi, err := precompilesgen.NewArbPriceFeed(common.HexToAddress("0x11a"), builder.L2.Client)
	Require(t, err)

	callOpts := &bind.CallOpts{Context: ctx}

	// Create tx that will trigger block production
	builder.L2Info.GenerateAccount("User")
	tx := builder.L2Info.PrepareTx("Faucet", "User", builder.L2Info.TransferGas, big.NewInt(1e12), nil)

	err = builder.L2.Client.SendTransaction(ctx, tx)
	Require(t, err)

	_, err = builder.L2.EnsureTxSucceeded(tx)
	Require(t, err)

	userAuth := builder.L2Info.GetDefaultTransactOpts("User", ctx)

	userAuth.GasLimit = 1_000_000

	receipt, err := WaitForTx(ctx, builder.L2.Client, tx.Hash(), time.Second*10)
	Require(t, err)
	fmt.Println("Transaction mined in block:", receipt.BlockNumber)

	num, err := arbHi.GetLatestBtcPrice(callOpts)
	Require(t, err)
	if num == 0 {
		t.Error("Expected price to be not zero")
	}
}
