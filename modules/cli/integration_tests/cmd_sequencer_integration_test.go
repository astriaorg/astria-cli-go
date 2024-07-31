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
const TestFromAddress = "astria1rsxyjrcm255ds9euthjx6yc3vrjt9sxrm9cfgm"
const TestTo = "astria1xnlvg0rle2u6auane79t4p27g8hxnj36ja960z"
const TestToPubKey = "88787e29db8d5247c6adfac9909b56e6b2705c3120b2e3885e8ec8aa416a10f1"
const TransferAmount = 535353
const SequencerURL = "http://127.0.0.1:26657"
const SequencerChainID = "sequencer-test-chain-0"

const TestBinPath = "../../../bin/astria-go-testy"

// NOTE - the ordering of these tests matters

func TestAddAndRemoveFeeAssets(t *testing.T) {
	key := fmt.Sprintf("--privkey=%s", TestFromPrivKey)
	// test synchronously first

	// add a fee asset
	addFeeAssetCmdSync := exec.Command(TestBinPath, "sequencer", "sudo", "fee-asset", "add", "bananas", key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID)
	_, err := addFeeAssetCmdSync.CombinedOutput()
	assert.NoError(t, err)

	// remove a fee asset
	removeFeeAssetCmdSync := exec.Command(TestBinPath, "sequencer", "sudo", "fee-asset", "remove", "bananas", key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID)
	_, err = removeFeeAssetCmdSync.CombinedOutput()
	assert.NoError(t, err)

	// NOTE - test synchronously
	// add a fee asset
	testAssetName := "testAsset"
	addFeeAssetCmd := exec.Command(TestBinPath, "sequencer", "sudo", "fee-asset", "add", testAssetName, key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID, "--async")
	_, err = addFeeAssetCmd.CombinedOutput()
	assert.NoError(t, err)

	// remove a fee asset
	removeFeeAssetCmd := exec.Command(TestBinPath, "sequencer", "sudo", "fee-asset", "remove", testAssetName, key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID, "--async")
	_, err = removeFeeAssetCmd.CombinedOutput()
	assert.NoError(t, err)

	// wait a bit for async tx to be processed
	time.Sleep(2 * time.Second)
}

func TestTransferAndGetNonce(t *testing.T) {
	// get initial blockheight
	getBlockHeightCmd := exec.Command(TestBinPath, "sequencer", "blockheight", "--json", "--sequencer-url", SequencerURL)
	// FIXME - using Output (vs CombinedOutput) only returns Stdout, but using CombinedOutput returns all logging as well,
	//  which can break json marshalling. how can we use CombinedOutput so that we get stderr info in case something fails?
	//  would rather not do any crazy string filtering/matching, but might be a good idea for the test.
	//  could also use jq in the test? probably best option.
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
	getNonceCmd := exec.Command(TestBinPath, "sequencer", "nonce", TestFromAddress, "--json", "--sequencer-url", SequencerURL)
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
	getBalanceCmd := exec.Command(TestBinPath, "sequencer", "balances", TestTo, "--json", "--sequencer-url", SequencerURL)
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
	transferCmd := exec.Command(TestBinPath, "sequencer", "transfer", amtStr, TestTo, key, "--sequencer-chain-id", SequencerChainID, "--sequencer-url", SequencerURL)
	transferOutput, err := transferCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to transfer: %s, %v", transferOutput, err)
	}

	// get blockheight after transfer
	getBlockHeightAfterCmd := exec.Command(TestBinPath, "sequencer", "blockheight", "--json", "--sequencer-url", SequencerURL)
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
	getNonceAfterCmd := exec.Command(TestBinPath, "sequencer", "nonce", TestFromAddress, "--json", "--sequencer-url", SequencerURL)
	nonceAfterOutput, err := getNonceAfterCmd.CombinedOutput()
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
	getBalanceAfterCmd := exec.Command(TestBinPath, "sequencer", "balances", TestTo, "--json", "--sequencer-url", SequencerURL)
	balanceAfterOutput, err := getBalanceAfterCmd.CombinedOutput()
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
	transferCmd := exec.Command(TestBinPath, "sequencer", "transfer", "53", TestTo, key, secondKey, "--sequencer-url", SequencerURL)
	_, err := transferCmd.CombinedOutput()
	assert.Error(t, err)

	// test that we get error when no type of key passed in
	transferCmd = exec.Command(TestBinPath, "sequencer", "transfer", "53", TestTo, "--sequencer-url", SequencerURL)
	_, err = transferCmd.CombinedOutput()
	assert.Error(t, err)
}

func TestRemoveAndAddIBCRelayer(t *testing.T) {
	key := fmt.Sprintf("--privkey=%s", TestFromPrivKey)

	// remove an address from the existing IBC relayer set
	removeIBCRelayerCmd := exec.Command(TestBinPath, "sequencer", "sudo", "ibc-relayer", "remove", TestTo, key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID)
	_, err := removeIBCRelayerCmd.CombinedOutput()
	assert.NoError(t, err)

	// add same address back to the IBC relayer set
	addIBCRelayerCmd := exec.Command(TestBinPath, "sequencer", "sudo", "ibc-relayer", "add", TestTo, key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID)
	_, err = addIBCRelayerCmd.CombinedOutput()
	assert.NoError(t, err)

	// test asynchronously
	// remove an address from the existing IBC relayer set
	removeIBCRelayerCmdAsync := exec.Command(TestBinPath, "sequencer", "sudo", "ibc-relayer", "remove", TestTo, key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID, "--async")
	_, err = removeIBCRelayerCmdAsync.CombinedOutput()
	assert.NoError(t, err)

	// add same address back to the IBC relayer set
	addIBCRelayerCmdAsync := exec.Command(TestBinPath, "sequencer", "sudo", "ibc-relayer", "add", TestTo, key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID, "--async")
	_, err = addIBCRelayerCmdAsync.CombinedOutput()
	assert.NoError(t, err)

	// wait a bit for async tx to be processed
	time.Sleep(2 * time.Second)
}

func TestValidatorUpdate(t *testing.T) {
	// update the validator power
	key := fmt.Sprintf("--privkey=%s", TestFromPrivKey)
	validatorUpdateCmd := exec.Command(TestBinPath, "sequencer", "sudo", "validator-update", TestToPubKey, "100", key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID)
	_, err := validatorUpdateCmd.CombinedOutput()
	assert.NoError(t, err)

	// revert the validator power
	validatorUpdateCmd = exec.Command(TestBinPath, "sequencer", "sudo", "validator-update", TestToPubKey, "10", key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID)
	_, err = validatorUpdateCmd.CombinedOutput()
	assert.NoError(t, err)

	// test async
	// update the validator power
	validatorUpdateCmdAsync := exec.Command(TestBinPath, "sequencer", "sudo", "validator-update", TestToPubKey, "100", key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID, "--async")
	_, err = validatorUpdateCmdAsync.CombinedOutput()
	assert.NoError(t, err)

	// revert the validator power
	validatorUpdateCmdAsync = exec.Command(TestBinPath, "sequencer", "sudo", "validator-update", TestToPubKey, "10", key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID, "--async")
	_, err = validatorUpdateCmdAsync.CombinedOutput()
	assert.NoError(t, err)

	// wait a bit for async tx to be processed
	time.Sleep(2 * time.Second)
}

func TestUpdateSudoAddress(t *testing.T) {
	key := fmt.Sprintf("--privkey=%s", TestFromPrivKey)

	// change the sudo address
	addressChangeCmd := exec.Command(TestBinPath, "sequencer", "sudo", "sudo-address-change", TestTo, key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID)
	_, err := addressChangeCmd.CombinedOutput()
	assert.NoError(t, err)

	// async
	// change the sudo address
	addressChangeCmdAsync := exec.Command(TestBinPath, "sequencer", "sudo", "sudo-address-change", TestTo, key, "--sequencer-url", SequencerURL, "--sequencer-chain-id", SequencerChainID, "--async")
	_, err = addressChangeCmdAsync.CombinedOutput()
	assert.NoError(t, err)
}

func TestGetBlock(t *testing.T) {
	// get initial blockheight
	getBlockHeightCmd := exec.Command(TestBinPath, "sequencer", "blockheight", "--json", "--sequencer-url", SequencerURL)
	blockHeightOutput, err := getBlockHeightCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get blockheight: %s, %v", blockHeightOutput, err)
	}
	var blockHeight sequencer.BlockheightResponse
	err = json.Unmarshal(blockHeightOutput, &blockHeight)
	if err != nil {
		t.Fatalf("Failed to unmarshal blockheight json output: %v", err)
	}
	initialBlockHeight := blockHeight.Blockheight

	// get a block
	if initialBlockHeight > 0 {
		getBlockCmd := exec.Command(TestBinPath, "sequencer", "block", "1", "--json", "--sequencer-url", SequencerURL)
		_, err := getBlockCmd.CombinedOutput()
		assert.NoError(t, err)
	} else {
		t.Fatalf("Blockheight is 0, cannot get block")
	}
}
