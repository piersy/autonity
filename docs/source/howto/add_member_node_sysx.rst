Add a member's node to System X
===============================

To add a newly deployed node to System X, a System Operator must send a transaction to the Autonity Protocol method responsible for adding new nodes. That transaction must be signed by an authorised account i.e. Governance Operator and the account must have enough fee tokens to pay for the Gas.

Prerequisites
------------

- enode values received from the member's newly deployed node
- MetaMask is installed in the default browser (or any other web wallet of choice)
- Governance Operator and Treasury Operator addresses and private keys generated during the deployment of System X

Steps
-----

1. Open a terminal and run the following command to forward a deployed and running node to localhost:

  .. Code:: bash

    kubectl -n <node namespace> port-forward svc/autonity-node-0 8545:8545

2. Set up a new network in Metamask using the following information:


 

  .. table:: *Set up Metamask*
     :widths: auto

     ========================= ========================== 
      Address                        Value        
     ========================= ========================== 
     URL                        http://127.0.0.1
     ChainID                    1489  
     ========================= ========================== 

3. Locate the 2 files including the addresses and private keys of the Treasury Operator and Governance Operator.

4. Import the Treasury Operator account inside MetaMask using the Treasury Operator `account address and private key`. The balance of the Treasury Operator should be displayed inside MetaMask.

5. Import the Governance Operator account inside MetaMask using the Governance Operator's account `address and private key`. The Governance Operator's account should show 0. 

6. Send a transaction using MetaMask from the Treasury Operator to the Governance Operator with the following parameters:


  .. table:: *Send a transaction using Metamask*
     :widths: auto

     ========================= ========================== 
      Transaction                        Value        
     ========================= ========================== 
     Amount                     177 ETH
     Gas Price                  1,000 Gwei  
     Gas Limit                  21,000
     ========================= ========================== 


  Switch to the Governance Operator's account where the balance should be 177 ETH.  The Governance Operator account now has the required transaction fees to submit a transaction to add the new node to the network.

7. Run the following command to get the Autonity contract address:

  .. Code:: bash

    `curl -X POST -H "Content-Type: application/json" --data \
    '{"jsonrpc":"2.0","method":"tendermint_getContractAddress","params":[],"id":1}' \
    http://localhost:8545`

8. Clone the following repository:

  .. Code:: bash

    git clone git@github.com:clearmatics/governance-operator.git

  and change the current directory:

  .. Code:: bash

    cd governance-operator


9. Edit `config.json` and replace `address and private key` in the file with those of the Governance Operator:

  .. Code:: bash

      {
        "uri": "http://127.0.0.1:8545",
        "account": {
          "address": "0xe06e3238D28a7a661416CB65c8cecfFE47daB296",
          "privateKey": "ceb09f113efe9b65967ea8699b3ca6ae85aa1ac244546b3b5edc3675c1ee659b"
        },
        "gas": 200000,
        "gasPrice": 10000000000000
      }

10. Run the following command to send the transaction to the Autonity smart contract to add the new node and assign it a Validator role:

  .. Code:: bash 

    docker run -ti --rm -v $(pwd)/config.json:/governance-operator/config.json --net=host clearmatics/governance-operator \
    addValidator ${contract_addr} ${validator_addr} ${stake} ${enode}

The node is now added to System X.