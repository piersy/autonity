// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"fmt"
	"github.com/clearmatics/autonity/p2p/enode"
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/ethash"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/params"
)

// Tests that block headers can be retrieved from a remote chain based on user queries.
func TestGetBlockHeaders63(t *testing.T) { testGetBlockHeaders(t, 63) }
func TestGetBlockHeaders64(t *testing.T) { testGetBlockHeaders(t, 64) }

func testGetBlockHeaders(t *testing.T, protocol int) {
	p2pPeer := newTestP2PPeer("peer")
	pm, _ := newTestProtocolManagerMust(t, downloader.FullSync, downloader.MaxHashFetch+15, nil, nil, []string{p2pPeer.Info().Enode})
	peer, _ := newTestPeer(p2pPeer, protocol, pm, true)
	defer peer.close()

	// Create a "random" unknown hash for testing
	var unknown common.Hash
	for i := range unknown {
		unknown[i] = byte(i)
	}
	// Create a batch of tests for various scenarios
	limit := uint64(downloader.MaxHeaderFetch)
	tests := []struct {
		query  *getBlockHeadersData // The query to execute for header retrieval
		expect []common.Hash        // The hashes of the block whose headers are expected
		drop bool                   // Peer is dropped if peer is untrusted.
	}{
		// A single random block should be retrievable by hash and number too
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: pm.blockchain.GetBlockByNumber(limit / 2).Hash()}, Amount: 1},
			[]common.Hash{pm.blockchain.GetBlockByNumber(limit / 2).Hash()},
			false,
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Amount: 1},
			[]common.Hash{pm.blockchain.GetBlockByNumber(limit / 2).Hash()},
			false,
		},
		// Multiple headers should be retrievable in both directions
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Amount: 3},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 1).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 2).Hash(),
			},
			false,
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Amount: 3, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 1).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 2).Hash(),
			},
			false,
		},
		// Multiple headers with skip lists should be retrievable
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 4).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 8).Hash(),
			},
			false,
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 4).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 8).Hash(),
			},
			false,
		},
		// The chain endpoints should be retrievable
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 0}, Amount: 1},
			[]common.Hash{pm.blockchain.GetBlockByNumber(0).Hash()},
			false,
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64()}, Amount: 1},
			[]common.Hash{pm.blockchain.CurrentBlock().Hash()},
			false,
		},
		// Ensure protocol limits are honored
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() - 1}, Amount: limit + 10, Reverse: true},
			pm.blockchain.GetBlockHashesFromHash(pm.blockchain.CurrentBlock().Hash(), limit),
			false,
		},
		// Check that requesting more than available is handled gracefully
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() - 4}, Skip: 3, Amount: 3},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64() - 4).Hash(),
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64()).Hash(),
			},
			false,
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: 4}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(4).Hash(),
				pm.blockchain.GetBlockByNumber(0).Hash(),
			},
			false,
		},
		// Check that requesting more than available is handled gracefully, even if mid skip
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() - 4}, Skip: 2, Amount: 3},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64() - 4).Hash(),
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64() - 1).Hash(),
			},
			false,
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: 4}, Skip: 2, Amount: 3, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(4).Hash(),
				pm.blockchain.GetBlockByNumber(1).Hash(),
			},
			false,
		},
		// Check a corner case where requesting more can iterate past the endpoints
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 2}, Amount: 5, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(2).Hash(),
				pm.blockchain.GetBlockByNumber(1).Hash(),
				pm.blockchain.GetBlockByNumber(0).Hash(),
			},
			false,
		},
		// Check a corner case where skipping overflow loops back into the chain start
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: pm.blockchain.GetBlockByNumber(3).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64 - 1},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(3).Hash(),
			},
			false,
		},
		// Check a corner case where skipping overflow loops back to the same header
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: pm.blockchain.GetBlockByNumber(1).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(1).Hash(),
			},
			false,
		},
		// Check that non existing headers aren't returned
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: unknown}, Amount: 1},
			[]common.Hash{},
			false,
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() + 1}, Amount: 1},
			[]common.Hash{},
			false,
		},
		// Check that an untrusted peer should be dropped when it query block header data.
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: pm.blockchain.GetBlockByNumber(limit / 2).Hash()}, Amount: 1},
			[]common.Hash{pm.blockchain.GetBlockByNumber(limit / 2).Hash()},
			true,
		},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the headers to expect in the response
		headers := []*types.Header{}
		for _, hash := range tt.expect {
			headers = append(headers, pm.blockchain.GetBlockByHash(hash).Header())
		}

		if !tt.drop {
			// refresh untrusted peer list with white list
			whiteList := []*enode.Node{peer.Node()}
			pm.RefreshUntrustedPeers(whiteList)

			// Send the hash request and verify the response
			p2p.Send(peer.app, 0x03, tt.query)
			if err := p2p.ExpectMsg(peer.app, 0x04, headers); err != nil {
				t.Errorf("test %d: headers mismatch: %v", i, err)
			}
			// If the test used number origins, repeat with hashes as the too
			if tt.query.Origin.Hash == (common.Hash{}) {
				if origin := pm.blockchain.GetBlockByNumber(tt.query.Origin.Number); origin != nil {
					tt.query.Origin.Hash, tt.query.Origin.Number = origin.Hash(), 0

					p2p.Send(peer.app, 0x03, tt.query)
					if err := p2p.ExpectMsg(peer.app, 0x04, headers); err != nil {
						t.Errorf("test %d: headers mismatch: %v", i, err)
					}
				}
			}
		} else {
			pm.AddUntrustedPeer(crypto.PubkeyToAddress(*peer.Node().Pubkey()))
			// Send the hash request and verify the response
			p2p.Send(peer.app, 0x03, tt.query)
			// Verify that the remote peer is maintained or dropped
			if peers := pm.peers.Len(); peers != 0 {
				t.Fatalf("peer count mismatch: have %d, want %d", peers, 0)
			}
		}
	}
}

// Tests that block contents can be retrieved from a remote chain based on their hashes.
func TestGetBlockBodies63(t *testing.T) { testGetBlockBodies(t, 63) }
func TestGetBlockBodies64(t *testing.T) { testGetBlockBodies(t, 64) }

func testGetBlockBodies(t *testing.T, protocol int) {
	p2pPeer := newTestP2PPeer("peer")
	pm, _ := newTestProtocolManagerMust(t, downloader.FullSync, downloader.MaxBlockFetch+15, nil, nil, []string{p2pPeer.Info().Enode})
	peer, _ := newTestPeer(p2pPeer, protocol, pm, true)
	defer peer.close()

	// Create a batch of tests for various scenarios
	limit := downloader.MaxBlockFetch
	tests := []struct {
		random    int           // Number of blocks to fetch randomly from the chain
		explicit  []common.Hash // Explicitly requested blocks
		available []bool        // Availability of explicitly requested blocks
		expected  int           // Total number of existing blocks to expect
		drop      bool          // Peer is dropped if peer is untrusted.
	}{
		{1, nil, nil, 1, false},             // A single random block should be retrievable
		{10, nil, nil, 10, false},           // Multiple random blocks should be retrievable
		{limit, nil, nil, limit, false},     // The maximum possible blocks should be retrievable
		{limit + 1, nil, nil, limit, false}, // No more than the possible block count should be returned
		{0, []common.Hash{pm.blockchain.Genesis().Hash()}, []bool{true}, 1, false},      // The genesis block should be retrievable
		{0, []common.Hash{pm.blockchain.CurrentBlock().Hash()}, []bool{true}, 1, false}, // The chains head block should be retrievable
		{0, []common.Hash{{}}, []bool{false}, 0, false},                                 // A non existent block should not be returned

		// Existing and non-existing blocks interleaved should not cause problems
		{0, []common.Hash{
			{},
			pm.blockchain.GetBlockByNumber(1).Hash(),
			{},
			pm.blockchain.GetBlockByNumber(10).Hash(),
			{},
			pm.blockchain.GetBlockByNumber(100).Hash(),
			{},
		}, []bool{false, true, false, true, false, true, false}, 3, false},
		// Check that an untrusted peer should be dropped when it query block body data.
		{1, nil, nil, 1, true},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the hashes to request, and the response to expect
		hashes, seen := []common.Hash{}, make(map[int64]bool)
		bodies := []*blockBody{}

		for j := 0; j < tt.random; j++ {
			for {
				num := rand.Int63n(int64(pm.blockchain.CurrentBlock().NumberU64()))
				if !seen[num] {
					seen[num] = true

					block := pm.blockchain.GetBlockByNumber(uint64(num))
					hashes = append(hashes, block.Hash())
					if len(bodies) < tt.expected {
						bodies = append(bodies, &blockBody{Transactions: block.Transactions(), Uncles: block.Uncles()})
					}
					break
				}
			}
		}
		for j, hash := range tt.explicit {
			hashes = append(hashes, hash)
			if tt.available[j] && len(bodies) < tt.expected {
				block := pm.blockchain.GetBlockByHash(hash)
				bodies = append(bodies, &blockBody{Transactions: block.Transactions(), Uncles: block.Uncles()})
			}
		}

		if !tt.drop {
			// refresh untrusted peer list with white list
			whiteList := []*enode.Node{peer.Node()}
			pm.RefreshUntrustedPeers(whiteList)
			// Send the hash request and verify the response
			p2p.Send(peer.app, 0x05, hashes)
			if err := p2p.ExpectMsg(peer.app, 0x06, bodies); err != nil {
				t.Errorf("test %d: bodies mismatch: %v", i, err)
			}
		} else {
			pm.AddUntrustedPeer(crypto.PubkeyToAddress(*peer.Node().Pubkey()))
			// Send the hash request and verify the response
			p2p.Send(peer.app, 0x05, hashes)
			// Verify that the remote peer is maintained or dropped
			if peers := pm.peers.Len(); peers != 0 {
				t.Fatalf("peer count mismatch: have %d, want %d", peers, 0)
			}
		}
	}
}

// Tests that the node state database can be retrieved based on hashes.
func TestGetNodeData63(t *testing.T) { testGetNodeData(t, 63) }
func TestGetNodeData64(t *testing.T) { testGetNodeData(t, 64) }

func testGetNodeData(t *testing.T, protocol int) {
	tests := []struct{
		name string
		drop bool
	}{
		{"Handle get node data from trusted peer", false},
		{"Handle get node data from un-trusted peer", true},
	}

	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr := crypto.PubkeyToAddress(acc2Key.PublicKey)

	signer := types.HomesteadSigner{}
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *core.BlockGen) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some ether.
			tx, _ := types.SignTx(types.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(10000), params.TxGas, nil, nil), signer, testBankKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more ether to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := types.SignTx(types.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, testBankKey)
			tx2, _ := types.SignTx(types.NewTransaction(block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was mined by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetExtra([]byte("yeehaw"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			b2 := block.PrevBlock(1).Header()
			b2.Extra = []byte("foo")
			block.AddUncle(b2)
			b3 := block.PrevBlock(2).Header()
			b3.Extra = []byte("foo")
			block.AddUncle(b3)
		}
	}
	// Assemble the test environment
	p2pPeer := newTestP2PPeer("peer")
	pm, db := newTestProtocolManagerMust(t, downloader.FullSync, 4, generator, nil, []string{p2pPeer.Info().Enode})
	peer, _ := newTestPeer(p2pPeer, protocol, pm, true)
	defer peer.close()

	// Fetch for now the entire chain db
	hashes := []common.Hash{}

	it := db.NewIterator(nil, nil)
	for it.Next() {
		if key := it.Key(); len(key) == common.HashLength {
			hashes = append(hashes, common.BytesToHash(key))
		}
	}
	it.Release()

	for _, tt := range tests {
		if !tt.drop {
			// refresh untrusted peer list with white list
			whiteList := []*enode.Node{peer.Node()}
			pm.RefreshUntrustedPeers(whiteList)

			p2p.Send(peer.app, 0x0d, hashes)
			msg, err := peer.app.ReadMsg()
			if err != nil {
				t.Fatalf("failed to read node data response: %v", err)
			}
			if msg.Code != 0x0e {
				t.Fatalf("response packet code mismatch: have %x, want %x", msg.Code, 0x0c)
			}
			var data [][]byte
			if err := msg.Decode(&data); err != nil {
				t.Fatalf("failed to decode response node data: %v", err)
			}
			// Verify that all hashes correspond to the requested data, and reconstruct a state tree
			for i, want := range hashes {
				if hash := crypto.Keccak256Hash(data[i]); hash != want {
					t.Errorf("data hash mismatch: have %x, want %x", hash, want)
				}
			}
			statedb := rawdb.NewMemoryDatabase()
			for i := 0; i < len(data); i++ {
				statedb.Put(hashes[i].Bytes(), data[i])
			}
			accounts := []common.Address{testBank, acc1Addr, acc2Addr}
			for i := uint64(0); i <= pm.blockchain.CurrentBlock().NumberU64(); i++ {
				trie, _ := state.New(pm.blockchain.GetBlockByNumber(i).Root(), state.NewDatabase(statedb), nil)

				for j, acc := range accounts {
					state, _ := pm.blockchain.State()
					bw := state.GetBalance(acc)
					bh := trie.GetBalance(acc)

					if (bw != nil && bh == nil) || (bw == nil && bh != nil) {
						t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
					}
					if bw != nil && bh != nil && bw.Cmp(bw) != 0 {
						t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
					}
				}
			}
		} else {
			pm.AddUntrustedPeer(crypto.PubkeyToAddress(*peer.Node().Pubkey()))
			p2p.Send(peer.app, 0x0d, hashes)
			// Verify that the remote peer is maintained or dropped
			if peers := pm.peers.Len(); peers != 0 {
				t.Fatalf("peer count mismatch: have %d, want %d", peers, 0)
			}
		}
	}
}

// Tests that the transaction receipts can be retrieved based on hashes.
func TestGetReceipt63(t *testing.T) { testGetReceipt(t, 63) }
func TestGetReceipt64(t *testing.T) { testGetReceipt(t, 64) }

func testGetReceipt(t *testing.T, protocol int) {
	tests := []struct{
		name string
		drop bool
	}{
		{"Handle get transaction receipt from trusted peer", false},
		{"Handle get transaction receipt un-trusted peer", true},
	}

	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr := crypto.PubkeyToAddress(acc2Key.PublicKey)

	signer := types.HomesteadSigner{}
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *core.BlockGen) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some ether.
			tx, _ := types.SignTx(types.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(10000), params.TxGas, nil, nil), signer, testBankKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more ether to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := types.SignTx(types.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, testBankKey)
			tx2, _ := types.SignTx(types.NewTransaction(block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was mined by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetExtra([]byte("yeehaw"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			b2 := block.PrevBlock(1).Header()
			b2.Extra = []byte("foo")
			block.AddUncle(b2)
			b3 := block.PrevBlock(2).Header()
			b3.Extra = []byte("foo")
			block.AddUncle(b3)
		}
	}
	// Assemble the test environment
	p2pPeer := newTestP2PPeer("peer")
	pm, _ := newTestProtocolManagerMust(t, downloader.FullSync, 4, generator, nil, []string{p2pPeer.Info().Enode})
	peer, _ := newTestPeer(p2pPeer, protocol, pm, true)
	defer peer.close()

	// Collect the hashes to request, and the response to expect
	hashes, receipts := []common.Hash{}, []types.Receipts{}
	for i := uint64(0); i <= pm.blockchain.CurrentBlock().NumberU64(); i++ {
		block := pm.blockchain.GetBlockByNumber(i)

		hashes = append(hashes, block.Hash())
		receipts = append(receipts, pm.blockchain.GetReceiptsByHash(block.Hash()))
	}

	for _, tt := range tests {
		if !tt.drop {
			// refresh untrusted peer list with white list
			whiteList := []*enode.Node{peer.Node()}
			pm.RefreshUntrustedPeers(whiteList)
			// Send the hash request and verify the response
			p2p.Send(peer.app, 0x0f, hashes)
			if err := p2p.ExpectMsg(peer.app, 0x10, receipts); err != nil {
				t.Errorf("receipts mismatch: %v", err)
			}
		} else {
			pm.AddUntrustedPeer(crypto.PubkeyToAddress(*peer.Node().Pubkey()))
			// Send the hash request and verify the response
			p2p.Send(peer.app, 0x0f, hashes)
			// Verify that the remote peer is maintained or dropped
			if peers := pm.peers.Len(); peers != 0 {
				t.Fatalf("peer count mismatch: have %d, want %d", peers, 0)
			}
		}
	}
}

// Tests that post eth protocol handshake, clients perform a mutual checkpoint
// challenge to validate each other's chains. Hash mismatches, or missing ones
// during a fast sync should lead to the peer getting dropped.
func TestCheckpointChallenge(t *testing.T) {
	tests := []struct {
		syncmode   downloader.SyncMode
		checkpoint  bool
		timeout     bool
		empty       bool
		match       bool
		drop        bool
		trustedPeer bool
	}{
		// If checkpointing is not enabled locally, don't challenge and don't drop
		{downloader.FullSync, false, false, false, false, false, true},
		{downloader.FastSync, false, false, false, false, false, true},
		{downloader.FullSync, false, false, false, false, false, false},
		{downloader.FastSync, false, false, false, false, false, false},

		// If checkpointing is enabled locally and remote response is empty, only drop during fast sync
		{downloader.FullSync, true, false, true, false, false, true},
		{downloader.FastSync, true, false, true, false, true, true}, // Special case, fast sync, unsynced peer
		{downloader.FullSync, true, false, true, false, false, false},
		{downloader.FastSync, true, false, true, false, true, false}, // Special case, fast sync, unsynced peer

		// If checkpointing is enabled locally and remote response mismatches, always drop
		{downloader.FullSync, true, false, false, false, true, true},
		{downloader.FastSync, true, false, false, false, true, true},
		{downloader.FullSync, true, false, false, false, true, false},
		{downloader.FastSync, true, false, false, false, true, false},

		// If checkpointing is enabled locally and remote response matches, never drop
		{downloader.FullSync, true, false, false, true, false, true},
		{downloader.FastSync, true, false, false, true, false, true},
		{downloader.FullSync, true, false, false, true, false, false},
		{downloader.FastSync, true, false, false, true, false, false},

		// If checkpointing is enabled locally and remote times out, always drop
		{downloader.FullSync, true, true, false, true, true, true},
		{downloader.FastSync, true, true, false, true, true, true},
		{downloader.FullSync, true, true, false, true, true, false},
		{downloader.FastSync, true, true, false, true, true, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("sync %v checkpoint %v timeout %v empty %v match %v", tt.syncmode, tt.checkpoint, tt.timeout, tt.empty, tt.match), func(t *testing.T) {
			testCheckpointChallenge(t, tt.syncmode, tt.checkpoint, tt.timeout, tt.empty, tt.match, tt.drop, tt.trustedPeer)
		})
	}
}

func testCheckpointChallenge(t *testing.T, syncmode downloader.SyncMode, checkpoint bool, timeout bool, empty bool, match bool, drop bool, trustedPeer bool) {
	// Reduce the checkpoint handshake challenge timeout
	defer func(old time.Duration) { syncChallengeTimeout = old }(syncChallengeTimeout)
	syncChallengeTimeout = 250 * time.Millisecond

	// Initialize a chain and generate a fake CHT if checkpointing is enabled
	var (
		db     = rawdb.NewMemoryDatabase()
		config = new(params.ChainConfig)
	)
	p2pPeer := newTestP2PPeer("peer")
	config.AutonityContractConfig = &params.AutonityContractGenesis{
		Users: []params.User{
			{
				Enode: p2pPeer.Info().Enode,
				Type:  params.UserValidator,
				Stake: 1,
			},
		},
	}

	if err := config.AutonityContractConfig.Prepare(); err != nil {
		t.Fatal(err)
	}

	(&core.Genesis{
		Config:     config,
		Difficulty: big.NewInt(1),
	}).MustCommit(db) // Commit genesis block
	// If checkpointing is enabled, create and inject a fake CHT and the corresponding
	// chllenge response.
	var response *types.Header
	var cht *params.TrustedCheckpoint
	if checkpoint {
		index := uint64(rand.Intn(500))
		number := (index+1)*params.CHTFrequency - 1
		response = &types.Header{Number: big.NewInt(int64(number)), Extra: []byte("valid")}

		cht = &params.TrustedCheckpoint{
			SectionIndex: index,
			SectionHead:  response.Hash(),
		}
	}
	// Create a checkpoint aware protocol manager
	blockchain, err := core.NewBlockChain(db, nil, config, ethash.NewFaker(), vm.Config{}, nil, &core.TxSenderCacher{}, nil)
	if err != nil {
		t.Fatalf("failed to create new blockchain: %v", err)
	}
	pm, err := NewProtocolManager(config, cht, syncmode, DefaultConfig.NetworkId, new(event.TypeMux), &testTxPool{pool: make(map[common.Hash]*types.Transaction)}, ethash.NewFaker(), blockchain, db, 1, nil, nil)
	if err != nil {
		t.Fatalf("failed to start test protocol manager: %v", err)
	}
	pm.Start(1000)
	defer pm.Stop()

	// Connect a new peer and check that we receive the checkpoint challenge
	peer, _ := newTestPeer(p2pPeer, eth63, pm, true)
	defer peer.close()

	// the hand-shake and checkpoint checking should always works for node to get sync from untrusted peer before it get fully synced.
	if !trustedPeer {
		// refresh untrusted peer list with white list
		whiteList := []*enode.Node{peer.Node()}
		pm.RefreshUntrustedPeers(whiteList)
	} else {
		pm.AddUntrustedPeer(crypto.PubkeyToAddress(*peer.Node().Pubkey()))
	}

	if checkpoint {
		challenge := &getBlockHeadersData{
			Origin:  hashOrNumber{Number: response.Number.Uint64()},
			Amount:  1,
			Skip:    0,
			Reverse: false,
		}
		if err := p2p.ExpectMsg(peer.app, GetBlockHeadersMsg, challenge); err != nil {
			t.Fatalf("challenge mismatch: %v", err)
		}
		// Create a block to reply to the challenge if no timeout is simulated
		if !timeout {
			if empty {
				if err := p2p.Send(peer.app, BlockHeadersMsg, []*types.Header{}); err != nil {
					t.Fatalf("failed to answer challenge: %v", err)
				}
			} else if match {
				if err := p2p.Send(peer.app, BlockHeadersMsg, []*types.Header{response}); err != nil {
					t.Fatalf("failed to answer challenge: %v", err)
				}
			} else {
				if err := p2p.Send(peer.app, BlockHeadersMsg, []*types.Header{{Number: response.Number}}); err != nil {
					t.Fatalf("failed to answer challenge: %v", err)
				}
			}
		}
	}
	// Wait until the test timeout passes to ensure proper cleanup
	time.Sleep(syncChallengeTimeout + 300*time.Millisecond)

	// Verify that the remote peer is maintained or dropped
	if drop {
		if peers := pm.peers.Len(); peers != 0 {
			t.Fatalf("peer count mismatch: have %d, want %d", peers, 0)
		}
	} else {
		if peers := pm.peers.Len(); peers != 1 {
			t.Fatalf("peer count mismatch: have %d, want %d", peers, 1)
		}
	}
}

func TestBroadcastBlock(t *testing.T) {
	var tests = []struct {
		totalPeers        int
		broadcastExpected int
	}{
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 2},
		{5, 2},
		{9, 3},
		{12, 3},
		{16, 4},
		{26, 5},
		{100, 10},
	}
	for _, test := range tests {
		testBroadcastBlock(t, test.totalPeers, test.broadcastExpected)
	}
}

func testBroadcastBlock(t *testing.T, totalPeers, broadcastExpected int) {
	var (
		evmux  = new(event.TypeMux)
		pow    = ethash.NewFaker()
		db     = rawdb.NewMemoryDatabase()
		config = &params.ChainConfig{}
		gspec  = &core.Genesis{Config: config}
	)
	config.AutonityContractConfig = &params.AutonityContractGenesis{}

	p2pPeers := make([]*p2p.Peer, totalPeers)
	for i := 0; i < totalPeers; i++ {
		p2pPeers[i] = newTestP2PPeer(fmt.Sprintf("peer %d", i))
		config.AutonityContractConfig.Users = append(
			config.AutonityContractConfig.Users,
			params.User{
				Enode: p2pPeers[i].Info().Enode,
				Type:  params.UserValidator,
				Stake: 100,
			},
		)
	}
	if err := config.AutonityContractConfig.Prepare(); err != nil {
		t.Fatal(err)
	}
	gspec.Difficulty = big.NewInt(1)

	genesis := gspec.MustCommit(db)

	blockchain, err := core.NewBlockChain(db, nil, config, pow, vm.Config{}, nil, core.NewTxSenderCacher(), nil)
	if err != nil {
		t.Fatalf("failed to create new blockchain: %v", err)
	}
	pm, err := NewProtocolManager(config, nil, downloader.FullSync, DefaultConfig.NetworkId, evmux, &testTxPool{pool: make(map[common.Hash]*types.Transaction)}, pow, blockchain, db, 1, nil, nil)
	if err != nil {
		t.Fatalf("failed to start test protocol manager: %v", err)
	}

	pm.Start(1000)
	defer pm.Stop()
	var peers []*testPeer

	for i := 0; i < totalPeers; i++ {
		peer, errc := newTestPeer(p2pPeers[i], eth63, pm, true)
		go func() {
			for err := range errc {
				fmt.Println("testPeerErr", err)
			}
		}()
		defer peer.close()

		peers = append(peers, peer)
	}
	chain, _ := core.GenerateChain(gspec.Config, genesis, ethash.NewFaker(), db, 1, func(i int, gen *core.BlockGen) {})

	errCh := make(chan error, totalPeers)
	doneCh := make(chan struct{}, totalPeers)
	for _, peer := range peers {
		go func(p *testPeer) {
			if err := p2p.ExpectMsg(p.app, NewBlockMsg, &newBlockData{Block: chain[0], TD: new(big.Int).Add(genesis.Difficulty(), chain[0].Difficulty())}); err != nil {
				errCh <- err
			} else {
				doneCh <- struct{}{}
			}
		}(peer)
	}
	pm.BroadcastBlock(chain[0], true /*propagate*/)
	var received int
	for {
		select {
		case <-doneCh:
			received++

		case <-time.After(time.Second):
			if received != broadcastExpected {
				t.Errorf("broadcast count mismatch: have %d, want %d", received, broadcastExpected)
			}
			return

		case err = <-errCh:
			t.Fatalf("broadcast failed: %v", err)
		}
	}

}

// Tests that a propagated malformed block (uncles or transactions don't match
// with the hashes in the header) gets discarded and not broadcast forward.
func TestBroadcastMalformedBlock(t *testing.T) {
	// Create a live node to test propagation with
	var (
		engine = ethash.NewFaker()
		db     = rawdb.NewMemoryDatabase()
		config = &params.ChainConfig{}
		gspec  = &core.Genesis{Config: config}
	)
	config.AutonityContractConfig = &params.AutonityContractGenesis{}
	sourcePeer := newTestP2PPeer("source")
	sinkPeer := newTestP2PPeer("sink")
	config.AutonityContractConfig.Users = []params.User{{
		Enode: sourcePeer.Info().Enode,
		Type:  params.UserValidator,
		Stake: 1,
	}, {
		Enode: sinkPeer.Info().Enode,
		Type:  params.UserValidator,
		Stake: 1,
	}}

	if err := config.AutonityContractConfig.Prepare(); err != nil {
		t.Fatal(err)
	}
	gspec.Difficulty = big.NewInt(1)
	genesis := gspec.MustCommit(db)

	blockchain, err := core.NewBlockChain(db, nil, config, engine, vm.Config{}, nil, &core.TxSenderCacher{}, nil)
	if err != nil {
		t.Fatalf("failed to create new blockchain: %v", err)
	}
	pm, err := NewProtocolManager(config, nil, downloader.FullSync, DefaultConfig.NetworkId, new(event.TypeMux), new(testTxPool), engine, blockchain, db, 1, nil, nil)

	if err != nil {
		t.Fatalf("failed to start test protocol manager: %v", err)
	}
	pm.Start(2)
	defer pm.Stop()

	// Create two peers, one to send the malformed block with and one to check
	// propagation
	source, _ := newTestPeer(sourcePeer, eth63, pm, true)
	defer source.close()

	sink, _ := newTestPeer(sinkPeer, eth63, pm, true)
	defer sink.close()

	// Create various combinations of malformed blocks
	chain, _ := core.GenerateChain(gspec.Config, genesis, ethash.NewFaker(), db, 1, func(i int, gen *core.BlockGen) {})

	malformedUncles := chain[0].Header()
	malformedUncles.UncleHash[0]++
	malformedTransactions := chain[0].Header()
	malformedTransactions.TxHash[0]++
	malformedEverything := chain[0].Header()
	malformedEverything.UncleHash[0]++
	malformedEverything.TxHash[0]++

	// Keep listening to broadcasts and notify if any arrives
	notify := make(chan struct{}, 1)
	go func() {
		if _, err := sink.app.ReadMsg(); err == nil {
			notify <- struct{}{}
		}
	}()
	// Try to broadcast all malformations and ensure they all get discarded
	for _, header := range []*types.Header{malformedUncles, malformedTransactions, malformedEverything} {
		block := types.NewBlockWithHeader(header).WithBody(chain[0].Transactions(), chain[0].Uncles())
		if err := p2p.Send(source.app, NewBlockMsg, []interface{}{block, big.NewInt(131136)}); err != nil {
			t.Fatalf("failed to broadcast block: %v", err)
		}
		select {
		case <-notify:
			t.Fatalf("malformed block forwarded")
		case <-time.After(100 * time.Millisecond):
		}
	}
}
