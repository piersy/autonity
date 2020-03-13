Generate a Genesis file
=======================

The genesis file defines the first block in the chain and the first block defines which chain you want to join. You will modify this file by adding the Governance Operator address and other user enode addresses (that you receive from the System Operator). You will then be able to communicate with the network.

1. Download the following template genesis.yaml file from (https://raw.githubusercontent.com/clearmatics/charts-ose/master/stable/autonity/genesis.yaml)


2. Update the `genesis.yaml` file as follows:

- Replace the Governance Operator account address by copying the address from the provided genesis file under the `alloc` attribute:

	.. code:: bash

		{
		  .
		  .
		  .
		  "alloc": {
		    "0xea9D7B5BD9A35F00617e78E36AbE9d7d0efde3Fc": {
		      "balance": "0x200000000000000000000000000000000000000000000000000000000000000"
		    }
		  }
		  .
		  .
		  .
		} 



- Find the `users` section and replace the 'fqdn' values with the enodes URLs:

 .. code:: bash

	 "users": [
	        {
	          "enode": "fqdn://a74ec6a3ad49323a87ad75651ccac6217c48d940f93c05d82b2a6a4bd772e9f3444e631cefed86f8d3816d6a9a4f77bfd946fdf795a79b5e5981c90d21058529@35.246.127.20:30303",
	          "stake": 50000,
	          "type": "validator"
	        },
	        {
	          "enode": "fqdn://a4ecd9d7055fe95ebd8ee3f1efde9ee2b5c723f6d15b7807bd94c7f8ba022717e8fbaeb1a08c79cd0023310f10cc585ff5305c9d5a0e7f92fb3384274733f21a@35.230.150.24:30303",
	          "stake": 50000,
	          "type": "validator"
	        },
	        {
	          "enode": "fqdn://9f613f43ac99a7fed5c86491f66852932016dbdb21201fada697b2b5adc1bb76f9a543adb5839f022070f9a458d91cd979bc518c732e0848aad7ecbf832fb711@35.246.42.226:30303",
	          "stake": 50000,
	          "type": "validator"
	        },
	        {
	          "enode": "fqdn://441492fbe0e61baefa5549c4443c085b5223129ad9e10b0c0cfce46647f9b82b696cbe0a53ab5ea25145bc25cfd87529ee5e2ac3d421b7375c3f13ce59979fd8@35.234.155.208:30303",
	          "stake": 50000,
	          "type": "validator"
	        },
	        {
	          "enode": "fqdn://a4c1d42bd3da35d51f3903ebda4ca759cdda6bafc4c06ae3487452f89e0a3da2d7aaa4f321744da45547c437b041ca0590f4cc977daa9c0e2d1d3d50e9983fce@34.89.46.114:30303",
	          "type": "participant"
	        }
	      ]
	  .
	  .
	  .


Example Genesis file
---------------------

.. Code:: bash

	---
	 Options for genesis block. Must be the same for all validator nodes (optional value only for new networks)
	{
	  "config": {
	    "homesteadBlock": 0,
	    "eip150Block": 0,
	    "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	    "eip155Block": 0,
	    "eip158Block": 0,
	    "byzantiumBlock": 0,
	    "constantinopleBlock": 0,
	    "petersburgBlock": 0,
	    "tendermint": {
	      "policy": 0,
	      "block-period": 1
	    },
	    "chainId": 1488,
	    "autonityContract": {
	      "deployer": "0x0000000000000000000000000000000000000002",
	      "bytecode": "",
	      "abi": "",
	      "minGasPrice": 5000,
	      "operator": "0x0000000000000000000000000000000000000003",
	      "users": [
	        {
	          "enode": "enode://b568df46303a9ac470f89fa39caf53f3c9aa9fec6c31eebb8a4af3938e893f1b4f572b509dc6fdaaf6fabb813818bf60c9e62140d182216591365ad1c9d20713@validator-1:30303",
	          "type": "validator",
	          "stake": 10000
	        },
	        {
	          "enode": "enode://b568df46303a9ac470f89fa39caf53f3c9aa9fec6c31eebb8a4af3938e893f1b4f572b509dc6fdaaf6fabb813818bf60c9e62140d182216591365ad1c9d20713@validator-1:30303",
	          "type": "validator",
	          "stake": 10000
	        },
	        {
	          "enode": "enode://b568df46303a9ac470f89fa39caf53f3c9aa9fec6c31eebb8a4af3938e893f1b4f572b509dc6fdaaf6fabb813818bf60c9e62140d182216591365ad1c9d20713@validator-1:30303",
	          "type": "validator",
	          "stake": 10000
	        },
	        {
	          "enode": "enode://b568df46303a9ac470f89fa39caf53f3c9aa9fec6c31eebb8a4af3938e893f1b4f572b509dc6fdaaf6fabb813818bf60c9e62140d182216591365ad1c9d20713@validator-1:30303",
	          "type": "validator",
	          "stake": 10000
	        },
	        {
	          "address": "0x12334",
	          "type": "stakeholder",
	          "stake": 500000
	        },
	        {
	          "enode": "enode://f2b0b3e2d957c9637e472532fe261e6357f2ee058483d562408dcc93dd5f7d28a21053dc74397c684b39910eb9b7a730a8499036530d0f931c9720405f22b329@observer-0:30303",
	          "type": "participant"
	        }
	      ]
	    }
	  },
	  "nonce": "0x0",
	  "timestamp": "0x0",
	  "gasLimit": "0xffffffff",
	  "difficulty": "0x1",
	  "coinbase": "0x0000000000000000000000000000000000000000",
	  "number": "0x0",
	  "gasUsed": "0x0",
	  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	  "mixHash": "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
	  "alloc": {
	    "0x25EC79340dbABad2E9bFDc433FB4d6D32EA8760D": {
	      "balance": "0x200000000000000000000000000000000000000000000000000000000000000"
	    },
	    "0xCAB6B71b1786781882Fba3Cd170333f8C1E59a69": {
	      "balance": "0x200000000000000000000000000000000000000000000000000000000000000"
	    },
	    "0xaa29166CF676a91e001568dee93acb629BB1e859": {
	      "balance": "0x200000000000000000000000000000000000000000000000000000000000000"
	    },
	    "0x21774D78e191278033A3244C90F3809abf941657": {
	      "balance": "0x200000000000000000000000000000000000000000000000000000000000000"
	    }
	  }
	}




Description
-----------

- nonce: Unused value by default (PoW only), can be left 0 or random
- timestamp: timestamp of the genesis block, specify the date when the network starts mining. Can be left 0 for immediate start
- difficulty: Must be equal to "0x1" with BFT consensus
- coinbase: Can be left random
- number: Must be equal to 0
- gasUsed: Must be equal to 0
- parentHash: Must be equal to "0x0000000000000000000000000000000000000000000000000000000000000000".
- mixHash : Must be equal to "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365", identify a BFT block
- alloc : Initial native token allocation, *alloc.ADDRESS.balance = AMOUNT* will allocate *AMOUNT* to the account *ADDRESS*. (wei)

Config
-------

- homesteadBlock, eip150Block, eip155Block, eip158Block, byzantiumBlock, constantinopleBlock, petersburgBlock : Leave to 0 if no special requirements, specify at which block legacy ethereum hard-fork should occur
- tendermint.policy : Proposer selection mechanism, leave it to 0 for round-robin (default)
- tendermint.block-period : Minimal time between two consecutive blocks, default to 1 sec
- chainId : Used for tx signature generature, use a random number if you wishes to deploy several network with the same config to avoid replay attacks

Config.autonityContract
------------------------ 

- deployer: address deploying the Autonity Contract, use a random one if no special needs
- abi: abi of the Autonity Contract, leave empty for default
- bytecode : EVM bytecode for Autonity Contract, leave empty for default
- minGasPrice : Initial minimum gas price
- operator : Address of the operator account
 
An object _user_ belonging to autonityContract.users contains 4 fields: address, enode, type and stake.

- user.type must be defined and be either equal to _participant_, _stakeholder_ or _validator_
- user.enode _or_ user.address must be defined
- If both user.enode and user.address are defined, then the derived address from user.enode must be equal to user.address
- If user.enode is defined then the full node associated will be allowed to join the network
- user.stake must be nil or equal to 0 for users of type _participant_
- If user.type is _validator_ then user.enode must be defined


