//go:build integration_tests
// +build integration_tests

package integration_tests

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os/exec"
	"testing"
	"time"

	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/astriaorg/go-sequencer-client/client"
	"github.com/stretchr/testify/assert"
)

const TestFromPrivKey = "2bd806c97f0e00af1a1fc3328fa763a9269723c8db8fac4f93af71db186d6e90"
const TestFromAddress = "1c0c490f1b5528d8173c5de46d131160e4b2c0c3"
const TestTo = "34fec43c7fcab9aef3b3cf8aba855e41ee69ca3a"
const TransferAmount = 535353

type balances []client.BalanceResponse

func TestTransferAndGetNonce(t *testing.T) {
	//setUp()
	//defer tearDown()

	// get initial nonce
	getNonceCmd := exec.Command("../bin/astria-go-testy", "sequencer", "get-nonce", TestFromAddress, "--json")
	nonceOutput, err := getNonceCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get nonce: %s, %v", nonceOutput, err)
	}
	var nonce sequencer.NonceResponse
	err = json.Unmarshal(nonceOutput, &nonce)
	if err != nil {
		t.Fatalf("Failed to marshal nonce json output: %v", err)
	}
	initialNonce := nonce.Nonce

	// get initial balance
	getBalanceCmd := exec.Command("../bin/astria-go-testy", "sequencer", "get-balances", TestTo, "--json")
	balanceOutput, err := getBalanceCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get balance: %s, %v", balanceOutput, err)
	}
	var toBalances balances
	err = json.Unmarshal(balanceOutput, &toBalances)
	if err != nil {
		t.Fatalf("Failed to marshal balance json output: %v", err)
	}
	initialBalance := toBalances[0].Balance

	// transfer
	key := fmt.Sprintf("--privkey=%s", TestFromPrivKey)
	amtStr := fmt.Sprintf("%d", TransferAmount)
	transferCmd := exec.Command("../bin/astria-go-testy", "sequencer", "transfer", amtStr, TestTo, key)
	transferOutput, err := transferCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to transfer: %s, %v", transferOutput, err)
	}

	// wait for transaction to be processed
	// FIXME - this could be flaky. can we check for the tx?
	time.Sleep(2 * time.Second)

	// get nonce after transfer
	getNonceAfterCmd := exec.Command("../bin/astria-go-testy", "sequencer", "get-nonce", TestFromAddress, "--json")
	nonceAfterOutput, err := getNonceAfterCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get nonce: %s, %v", nonceAfterOutput, err)
	}
	var nonceAfter sequencer.NonceResponse
	err = json.Unmarshal(nonceAfterOutput, &nonceAfter)
	if err != nil {
		t.Fatalf("Failed to marshal nonce json output: %v", err)
	}
	finalNonce := nonceAfter.Nonce
	expectedFinalNonce := initialNonce + 1
	assert.Equal(t, expectedFinalNonce, finalNonce)

	// get balance after transfer
	getBalanceAfterCmd := exec.Command("../bin/astria-go-testy", "sequencer", "get-balances", TestTo, "--json")
	balanceAfterOutput, err := getBalanceAfterCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get balance: %s, %v", balanceAfterOutput, err)
	}
	var toBalancesAfter balances
	err = json.Unmarshal(balanceAfterOutput, &toBalancesAfter)
	if err != nil {
		t.Fatalf("Failed to marshal balance json output: %v", err)
	}
	expectedFinalBalance := big.NewInt(0).Add(initialBalance, big.NewInt(TransferAmount))
	finalBalance := toBalancesAfter[0].Balance
	assert.Equal(t, expectedFinalBalance.String(), finalBalance.String())
}

// TODO - move setup and teardown here and out of the justfile.

//// build the cli with a unique name just for testing
//func setUp() {
//	wd, err := os.Getwd()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(wd)
//	err = os.Chdir(wd)
//	if err != nil {
//		panic(err)
//	}
//	c := exec.Command("go build -o bin/astria-go-testy")
//	o, err := c.CombinedOutput()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(o)
//}
//
//func tearDown() {
//	// TODO - cleanup testy binary?
//}
//
//func getBinPath() string {
//	e, err := os.Executable()
//	if err != nil {
//		panic(err)
//	}
//	path := path.Dir(e)
//	return path
//}
