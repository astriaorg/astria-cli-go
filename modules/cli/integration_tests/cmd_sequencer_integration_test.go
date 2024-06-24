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

	"github.com/astriaorg/astria-cli-go/modules/cli/internal/sequencer"
	"github.com/stretchr/testify/assert"
)

const TestFromPrivKey = "2bd806c97f0e00af1a1fc3328fa763a9269723c8db8fac4f93af71db186d6e90"
const TestFromAddress = "1c0c490f1b5528d8173c5de46d131160e4b2c0c3"
const TestTo = "34fec43c7fcab9aef3b3cf8aba855e41ee69ca3a"
const TransferAmount = 535353

const TestBinPath = "../../../bin/astria-go-testy"

func TestCreateaccount(t *testing.T) {
	createaccountCmd := exec.Command(TestBinPath, "sequencer", "createaccount", "--insecure", "--json")
	createaccountOutput, err := createaccountCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create account: %s, %v", createaccountOutput, err)

	}
	var account sequencer.AccountJSON
	err = json.Unmarshal(createaccountOutput, &account)
	if err != nil {
		t.Fatalf("Failed to unmarshal account json output: %v", err)
	}
	assert.NotEmpty(t, account.Address, "Address should not be empty")
	assert.NotEmpty(t, account.PublicKey, "PublicKey should not be empty")
	assert.NotEmpty(t, account.PrivateKey, "PrivateKey should not be empty")
}

func TestTransferFlags(t *testing.T) {
	// test that we get error when too many flags passed in
	key := fmt.Sprintf("--privkey=%s", TestFromPrivKey)
	secondKey := fmt.Sprintf("--keyfile=/fake/file")
	transferCmd := exec.Command(TestBinPath, "sequencer", "transfer", "53", TestTo, key, secondKey, "--sequencer-url", "http://127.0.0.1:26657")
	_, err := transferCmd.CombinedOutput()
	assert.Error(t, err)

	// test that we get error when no type of key passed in
	transferCmd = exec.Command(TestBinPath, "sequencer", "transfer", "53", TestTo, "--sequencer-url", "http://127.0.0.1:26657")
	_, err = transferCmd.CombinedOutput()
	assert.Error(t, err)
}

func TestTransferAndGetNonce(t *testing.T) {
	// get initial blockheight
	getBlockHeightCmd := exec.Command(TestBinPath, "sequencer", "blockheight", "--json", "--sequencer-url", "http://127.0.0.1:26657")
	blockHeightOutput, err := getBlockHeightCmd.Output()
	if err != nil {
		t.Fatalf("Failed to get blockheight: %s, %v", blockHeightOutput, err)
	}
	var blockHeight sequencer.BlockheightResponse
	err = json.Unmarshal(blockHeightOutput, &blockHeight)
	if err != nil {
		t.Fatalf("Failed to unmarshal blockheight json output: %v", err)
	}
	initialBlockHeight := blockHeight.Blockheight

	// get initial nonce
	getNonceCmd := exec.Command(TestBinPath, "sequencer", "nonce", TestFromAddress, "--json", "--sequencer-url", "http://127.0.0.1:26657")
	nonceOutput, err := getNonceCmd.Output()
	if err != nil {
		t.Fatalf("Failed to get nonce: %s, %v", nonceOutput, err)
	}
	var nonce sequencer.NonceResponse
	err = json.Unmarshal(nonceOutput, &nonce)
	if err != nil {
		t.Fatalf("Failed to unmarshal nonce json output: %v", err)
	}
	initialNonce := nonce.Nonce

	// get initial balance
	getBalanceCmd := exec.Command(TestBinPath, "sequencer", "balances", TestTo, "--json", "--sequencer-url", "http://127.0.0.1:26657")
	balanceOutput, err := getBalanceCmd.Output()
	if err != nil {
		t.Fatalf("Failed to get balance: %s, %v", balanceOutput, err)
	}
	var toBalances sequencer.BalancesResponse
	err = json.Unmarshal(balanceOutput, &toBalances)
	if err != nil {
		t.Fatalf("Failed to unmarshal balance json output: %v", err)
	}
	initialBalance := toBalances[0].Balance

	// transfer
	key := fmt.Sprintf("--privkey=%s", TestFromPrivKey)
	amtStr := fmt.Sprintf("%d", TransferAmount)
	transferCmd := exec.Command(TestBinPath, "sequencer", "transfer", amtStr, TestTo, key, "--sequencer-chain-id", "sequencer-test-chain-0", "--sequencer-url", "http://127.0.0.1:26657")
	transferOutput, err := transferCmd.Output()
	if err != nil {
		t.Fatalf("Failed to transfer: %s, %v", transferOutput, err)
	}

	// wait for transaction to be processed
	// FIXME - this could be flaky. can we check for the tx?
	time.Sleep(2 * time.Second)

	// get blockheight after transfer
	getBlockHeightAfterCmd := exec.Command(TestBinPath, "sequencer", "blockheight", "--json", "--sequencer-url", "http://127.0.0.1:26657")
	blockHeightAfterOutput, err := getBlockHeightAfterCmd.Output()
	if err != nil {
		t.Fatalf("Failed to get blockheight: %s, %v", blockHeightAfterOutput, err)
	}
	var blockHeightAfter sequencer.BlockheightResponse
	err = json.Unmarshal(blockHeightAfterOutput, &blockHeightAfter)
	if err != nil {
		t.Fatalf("Failed to unmarshal blockheight json output: %v", err)
	}
	finalBlockHeight := blockHeightAfter.Blockheight
	assert.Greaterf(t, finalBlockHeight, initialBlockHeight, "Blockheight should increase")

	// get nonce after transfer
	getNonceAfterCmd := exec.Command(TestBinPath, "sequencer", "nonce", TestFromAddress, "--json", "--sequencer-url", "http://127.0.0.1:26657")
	nonceAfterOutput, err := getNonceAfterCmd.Output()
	if err != nil {
		t.Fatalf("Failed to get nonce: %s, %v", nonceAfterOutput, err)
	}
	var nonceAfter sequencer.NonceResponse
	err = json.Unmarshal(nonceAfterOutput, &nonceAfter)
	if err != nil {
		t.Fatalf("Failed to unmarshal nonce json output: %v", err)
	}
	finalNonce := nonceAfter.Nonce
	expectedFinalNonce := initialNonce + 1
	assert.Equal(t, expectedFinalNonce, finalNonce)

	// get balance after transfer
	getBalanceAfterCmd := exec.Command(TestBinPath, "sequencer", "balances", TestTo, "--json", "--sequencer-url", "http://127.0.0.1:26657")
	balanceAfterOutput, err := getBalanceAfterCmd.Output()
	if err != nil {
		t.Fatalf("Failed to get balance: %s, %v", balanceAfterOutput, err)
	}
	var toBalancesAfter sequencer.BalancesResponse
	err = json.Unmarshal(balanceAfterOutput, &toBalancesAfter)
	if err != nil {
		t.Fatalf("Failed to unmarshal balance json output: %v", err)
	}
	expectedFinalBalance := big.NewInt(0).Add(initialBalance, big.NewInt(TransferAmount))
	finalBalance := toBalancesAfter[0].Balance
	assert.Equal(t, expectedFinalBalance.String(), finalBalance.String())
}
