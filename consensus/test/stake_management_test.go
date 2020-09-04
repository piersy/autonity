package test

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

/*
  In this file, it create 3 test cases which have similar work flow base on the local e2e test framework's main flow.

  First it setup an autontiy network by according to the genesis hook function, then from the specific chain height, it
  start to issue transactions via the transaction hook function specified for the target node, for example in the
  mintStakeHook, redeemStakeHook, and sendStakeHook, it issues transaction to call autonity contract via operator account
  to manage the stakes on the members.

  Then the test case verify the output from its finalAssert hook function on the specified height of the blockchain, for
  example, it checks the stake balance in different height to compare if the balance is expected.
*/

func TestStakeManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	stakeDelta := new(big.Int).SetUint64(50)
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	// mint stake hook
	mintStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			contract, err := autonityInstance(validator.rpcPort)
			if err != nil {
				return true, nil, err
			}
			defer contract.Close()

			txOpt, err := contract.transactionOpts(operatorKey)
			if err != nil {
				return true, nil, err
			}

			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			_, err = contract.MintStake(txOpt, *validatorsList[0].Address, stakeDelta)
			if err != nil {
				return true, nil, err
			}
		}
		return false, nil, nil
	}

	redeemStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			contract, err := autonityInstance(validator.rpcPort)
			if err != nil {
				return true, nil, err
			}
			defer contract.Close()

			txOpt, err := contract.transactionOpts(operatorKey)
			if err != nil {
				return true, nil, err
			}

			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			_, err = contract.RedeemStake(txOpt, *validatorsList[0].Address, stakeDelta)
			if err != nil {
				return true, nil, err
			}
		}
		return false, nil, nil
	}

	sendStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			contract, err := autonityInstance(validator.rpcPort)
			if err != nil {
				return true, nil, err
			}
			defer contract.Close()

			txOpt, err := contract.transactionOpts(validator.privateKey)
			if err != nil {
				return true, nil, err
			}

			senderAddress := crypto.PubkeyToAddress(validator.privateKey.PublicKey)
			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			toIndex := 0
			if senderAddress == *validatorsList[toIndex].Address {
				toIndex = 1
			}
			_, err = contract.Send(txOpt, *validatorsList[toIndex].Address, stakeDelta)
			if err != nil {
				return true, nil, err
			}
		}
		return false, nil, nil
	}

	// genesis hook
	genesisHook := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}
		return g
	}

	stakeCheckerHook := func(t *testing.T, validators map[string]*testNode) error {
		contract, err := autonityInstance(validators["VA"].rpcPort)
		if err != nil {
			return err
		}
		defer contract.Close()

		callOpt := contract.callOpts(3)
		initNetworkMetrics, err := contract.DumpEconomicsMetricData(callOpt)
		if err != nil {
			return err
		}

		callOpt.BlockNumber.SetUint64(validators["VA"].lastBlock)
		curNetworkMetrics, err := contract.DumpEconomicsMetricData(callOpt)
		if err != nil {
			return err
		}
		validatorsList := validators["VA"].service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		// check account stake balance.
		found := false
		for index, v := range initNetworkMetrics.Accounts {
			if v == *validatorsList[0].Address {
				found = true
				initBalance := initNetworkMetrics.Stakes[index]
				newBalance := curNetworkMetrics.Stakes[index]
				delta := new(big.Int).Abs(initBalance.Sub(initBalance, newBalance))
				assert.Equal(t, delta.Uint64(), stakeDelta.Uint64(), "stake balance is not expected")
				totalDelta := new(big.Int).Abs(initNetworkMetrics.Stakesupply.Sub(initNetworkMetrics.Stakesupply, curNetworkMetrics.Stakesupply))
				assert.Equal(t, totalDelta.Uint64(), stakeDelta.Uint64(), "stake total supply is not expected")
			}
		}
		assert.True(t, found, "cannot find wanted account from chain DB")
		return nil
	}

	stakeSendCheckerHook := func(t *testing.T, validators map[string]*testNode) error {
		contract, err := autonityInstance(validators["VA"].rpcPort)
		if err != nil {
			return err
		}
		defer contract.Close()

		callOpt := contract.callOpts(3)
		initNetworkMetrics, err := contract.DumpEconomicsMetricData(callOpt)
		if err != nil {
			return err
		}

		callOpt.BlockNumber.SetUint64(validators["VA"].lastBlock)
		curNetworkMetrics, err := contract.DumpEconomicsMetricData(callOpt)
		if err != nil {
			return err
		}

		validatorsList := validators["VA"].service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		senderAddr := crypto.PubkeyToAddress(validators["VA"].privateKey.PublicKey)
		toIndex := 0
		if senderAddr == *validatorsList[toIndex].Address {
			toIndex = 1
		}
		receiverAddr := *validatorsList[toIndex].Address

		// check account stake balance.
		senderPassed := false
		receiverPassed := false
		for index, v := range curNetworkMetrics.Accounts {
			if v == senderAddr {
				curBalance := curNetworkMetrics.Stakes[index].Uint64()
				initBalance := initNetworkMetrics.Stakes[index].Uint64()
				delta := initBalance - curBalance
				if delta == stakeDelta.Uint64() {
					senderPassed = true
					continue
				}
			}
			if v == receiverAddr {
				curBalance := curNetworkMetrics.Stakes[index].Uint64()
				initBalance := initNetworkMetrics.Stakes[index].Uint64()
				delta := curBalance - initBalance
				if delta == stakeDelta.Uint64() {
					receiverPassed = true
					continue
				}
			}
		}

		assert.Equal(t, senderPassed, true, "sender stake balance checking failed")
		assert.Equal(t, receiverPassed, true, "receiver stake balance checking failed")
		assert.Equal(t, initNetworkMetrics.Stakesupply.Uint64(), curNetworkMetrics.Stakesupply.Uint64(), "total stake supply is not expected")
		return nil
	}

	testCases := []*testCase{
		{
			name:          "stake management test mint stake",
			numValidators: 6,
			numBlocks:     20,
			txPerPeer:     1,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VA": mintStakeHook,
			},
			genesisHook: genesisHook,
			finalAssert: stakeCheckerHook,
		},
		{
			name:          "stake management test redeem stake",
			numValidators: 6,
			numBlocks:     20,
			txPerPeer:     1,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VA": redeemStakeHook,
			},
			genesisHook: genesisHook,
			finalAssert: stakeCheckerHook,
		},
		{
			name:          "stake management test send stake",
			numValidators: 6,
			numBlocks:     20,
			txPerPeer:     1,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VA": sendStakeHook,
			},
			genesisHook: genesisHook,
			finalAssert: stakeSendCheckerHook,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}