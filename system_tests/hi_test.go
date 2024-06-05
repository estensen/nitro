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

func TestHi(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the node with default configuration
	builder := NewNodeBuilder(ctx).DefaultConfig(t, true)
	cleanup := builder.Build(t)
	defer cleanup()

	// Initialize the ArbPriceFeed precompile
	arbHi, err := precompilesgen.NewArbHi(common.HexToAddress("0x11a"), builder.L2.Client)
	Require(t, err)

	callOpts := &bind.CallOpts{Context: ctx}

	builder.L2Info.GenerateAccount("User")
	builder.L2.TransferBalance(t, "Faucet", "User", big.NewInt(1_000_000_000_000_000), builder.L2Info)
	userAuth := builder.L2Info.GetDefaultTransactOpts("User", ctx)

	userAuth.GasLimit = 1_000_000

	tx, err := arbHi.SetNumber(&userAuth, 42)
	Require(t, err)

	receipt, err := WaitForTx(ctx, builder.L2.Client, tx.Hash(), time.Second*10)
	Require(t, err)
	fmt.Println("Transaction mined in block:", receipt.BlockNumber)

	num, err := arbHi.GetNumber(callOpts)
	Require(t, err)
	if num != 42 {
		t.Errorf("Expected num to be 42, was: %d", num)
	}
}
