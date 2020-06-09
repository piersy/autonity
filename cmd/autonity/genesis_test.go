// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/clearmatics/autonity/cmd/gengen/gengen"
	"github.com/clearmatics/autonity/common/math"
	"github.com/stretchr/testify/require"
)

// Tests that initializing Autonity with a custom genesis block and chain definitions
// work properly.
func TestCustomGenesis(t *testing.T) {
	// Create a temporary data directory to use and inspect later
	datadir := tmpdir(t)
	defer os.RemoveAll(datadir)

	// See doc in genesis package for description of the string format for a user.
	genesis, _, err := gengen.NewGenesis(1, []string{"1e12,v,1,:6789"}, nil)
	require.NoError(t, err)

	// Initialize the data directory with the custom genesis block
	genesisFile, err := os.Create(filepath.Join(datadir, "genesis.json"))
	require.NoError(t, err)
	err = json.NewEncoder(genesisFile).Encode(genesis)
	require.NoError(t, err)

	runAutonity(t, "--datadir", datadir, "init", genesisFile.Name()).WaitExit()
	query := "eth.getBlock(0).nonce"
	expectedResult, err := math.HexOrDecimal64(genesis.Nonce).MarshalText()
	require.NoError(t, err)

	// Query the custom genesis block
	autonity := runAutonity(t,
		"--datadir", datadir, "--maxpeers", "0", "--port", "0",
		"--nodiscover", "--nat", "none", "--ipcdisable",
		"--exec", query, "console")
	autonity.ExpectRegexp(string(expectedResult))
	autonity.ExpectExit()
}
