Deploy an Autonity network
==========================

Deploy an Autonity node to a Kubernetes cluster using Helm package manager. For this example, the Kubernetes cluster is hosted on Google Cloud Platform (GCP). There are three steps:

**Step 1:** Connect to Google Cloud Platform

**Step 2:** Deploy an Autonity network to a Kubernetes cluster on GCP

**Step 3:** Define secure communication between nodes


Prerequisite
------------


Install the Google SDK
***************

1. Install Python v2.7.0 on your system:

	.. code:: bash

	 sudo apt install python2.7

2. Run:

	.. code:: bash
 
 		curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-272.0.0-linux-x86_64.tar.gz

 		As part of the gcloud install.

3.   Extract the SDK to a new folder:

		.. code:: bash

 			tar zxvf [ARCHIVE_FILE] google-cloud-sdk 

4.  Remove the [ARCHIVE_FILE]: 

		.. code:: bash

			rm [ARCHIVE_FILE]

5.  Install the gcloud SDK with: 

		.. code:: bash

			./google-cloud-sdk/install.sh

		
	Opt out, (n) and proceed. 
			
	Enter (y) to update your $PATH and enable shell command completion.

6. Initialize the SDK: 

		.. code:: bash

			gcloud init

	You'll see a message like this:
	
	.. code:: bash

	    To continue, you must log in. Would you like to log in (Y/n)? Y
	    -Accept the option to log in using your Google user account. 
	    -A browser opens, log in to your Google user account when prompted and click 'Allow' to grant permission to access Google Cloud Platform resources.

Step 1: Connect to Google Cloud Platform
----------------------------------------

	.. NOTE:: The GCP infrastructure uses GSuite for IAM. This allows fine-grained access control levels to the infrastructure for internal Clearmatics engineers and staff, without needing to manage other directory services or accounts. Products external to GCP can either have a mapping to service accounts in GSuite, or will require a separate exercise, which is outside of the scope of this document

1. Create an access configuration:

	.. code:: bash

		gcloud config configurations create <local name> --account <your user>@clearmatics.com

	Replace the following in the command above:

		- <local name>: any name for local use to define the access to GCP

		- <your user>@clearmatics.com: company email username

2. Login to GCP using the credentials you created in Step 1:

	.. code:: bash
		
		gcloud auth login

	The browser opens automatically. Navigate to:

	.. code:: bash

		https://console.cloud.google.com 


	to access Google Cloud Console. Follow the steps to gain access to Google Cloud Platform.

3. To select a project from the available projects inside GCP, run:

	.. code:: bash

		gcloud projects list

	then, set the project to work under by running:

	.. code:: bash

		gcloud config set project <project name>

	Replace the <project name> with the name of the project you selected from the list returned above.

4. Check that the account added in Step 1 has been added to the project and is activated by running:

	.. code:: bash

		gcloud config configurations list

	The returned result includes:

		**NAME**: local name entered in Step 1

		**IS_ACTIVE**: True if the account has been setup properly

		**ACCOUNT**: the email username used in Step 1

		**PROJECT**: The project selected in Step 3

5. Select a cluster to deploy the network under by getting a list of available clusters by running:

	.. code:: bash

		gcloud container clusters list

6. Get credentials and generate config for kubectl:

	.. code:: bash

		gcloud container clusters get-credentials <cluster> --region <region>

	Replace the following in the above command:

	* **cluster** - one of the clusters names returned in Step 5

	* **region** - the location of the cluster returned in Step 5

7. Check the access to the cluster by running:

	.. code:: bash

		kubectl version

	You'll see a message like this:

	.. code:: bash

		Client Version: version.Info{Major:"1", Minor:"16", GitVersion:"v1.16.2", GitCommit:"c97fe5036ef3df2967d086711e6c0c405941e14b", GitTreeState:"clean", BuildDate:"2019-10-15T19:18:23Z", GoVersion:"go1.12.10", Compiler:"gc", Platform:"linux/amd64"}
		Server Version: version.Info{Major:"1", Minor:"14+", GitVersion:"v1.14.6-gke.1", GitCommit:"61c30f98599ad5309185df308962054d9670bafa", GitTreeState:"clean", BuildDate:"2019-08-28T11:06:42Z", GoVersion:"go1.12.9b4", Compiler:"gc", Platform:"linux/amd64"}




Step 2: Deploy an Autonity network to a Kubernetes cluster
---------------------------------------------------------

The deployed network has the following structure:

		* 4 nodes with the same genesis file and a System Operator

		* The 4 nodes will be validator nodes (this is the minimum required for the network to run) and they can run from different domains/cloud platforms or from the same cloud platform

		* The System Operator has two accounts: 

			- a Governance Account
			- a Treasury Account

1. Download the following template genesis.yaml file from (https://raw.githubusercontent.com/clearmatics/charts-ose/master/stable/autonity/genesis.yaml).

	-  Install Docker

	-  Generate account addresses and private keys for the Governance Operator and Treasury Operator:

		.. code:: bash

			docker run --rm clearmatics/eth-keys-generator > Governance_Operator
		  	docker run --rm clearmatics/eth-keys-generator > Treasure_Operator

	- Open the genesis.yaml file and replace the Governance Operator's address at the address value of the "operator" attribute with the account address in the file name "Governance_Operator":

		.. code:: bash

			...
			...
			...
			minGasPrice: 10000000000000
			operator: '0xae223655126e514C9C80096d99765A98547247D3'
			users:
			...
			...
			...



	- Replace the Treasury Operator's address by replacing the address value for the "alloc" attribute with the account address in the Treasure_Operator:

		.. code:: bash

			...
		  	...
		  	...
		  	alloc: 0x981dfa463bE3c247Fe311a05EeD4c67265204417:
		   	balance: '0x200000000000000000000000000000000000000000000000000000000000000'
		  	...
		  	...
		  	...

	- Change the name of the subdomain name in the DNS values by changing the values for the "enode" attributes to have a subdomain for every node:

		.. code:: bash

				...
			  	...
			  	...
			  	users:
			  		\- enode: fqdn://<node name>.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com
			    	type: validator
			    	stake: 50000
			  		\- enode: fqdn://<node name>.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com
			    	type: validator
				    stake: 50000
				  \- enode: fqdn://<node name>.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com
				    type: validator
				    stake: 50000
				  \- enode: fqdn://<node name>.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com
				    type: validator
				    stake: 50000
				  \- enode: fqdn://<node name>.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com
				    type: participant
				    ...
				    ...
				    ...
2. (optional)  Create a namespace under the cluster in Kubernetes for better organisation of the cluster:

	.. code:: bash

	  kubectl create namespace <namespace name>

	.. note:: Create a namespace for every node. Redundant node names across namespaces is forbidden

3. Deploy the node in a specific cluster with a specific namespace:

	.. code:: bash

	  	helm install --name $"<node name>" --namespace $"<namespace name>" charts-ose.clearmatics.com/autonity -f $"<full path to gensis.yaml file>"


4. Repeat Step 1 for all nodes listed in the 'user' section of the genesis.yaml file.

	When you have finished you see a message similar to:


	.. code:: bash

		NAME:   val-3-se
		LAST DEPLOYED: Mon Nov  4 17:57:42 2019
		NAMESPACE: how-to-se-03
		STATUS: DEPLOYED

		RESOURCES:
		==> v1/ConfigMap
		NAME              DATA  AGE
		autonity-node-0   2     6s
		autonity-tests    1     6s
		genesis           0     6s
		genesis-template  1     6s
		nginx-conf        1     6s

		==> v1/Job
		NAME                             COMPLETIONS  DURATION  AGE
		init-job02-genesis-configurator  0/1          6s        6s

		==> v1/Pod(related)
		NAME                                   READY  STATUS    RESTARTS  AGE
		autonity-node-0-67b5fdf8b-xgwmj        0/2    Init:0/1  0         6s
		init-job02-genesis-configurator-8xsv8  1/1    Running   0         6s

		==> v1/Role
		NAME           AGE
		genesis-write  6s
		secrets-write  6s

		==> v1/RoleBinding
		NAME           AGE
		genesis-write  6s
		secrets-write  6s

		==> v1/Secret
		NAME             TYPE    DATA  AGE
		autonity-node-0  Opaque  2     6s

		==> v1/Service
		NAME                 TYPE          CLUSTER-IP    EXTERNAL-IP  PORT(S)                     AGE
		autonity-node-0      ClusterIP     10.7.245.103  <none>       8545/TCP,8546/TCP,9200/TCP  6s
		p2p-autonity-node-0  LoadBalancer  10.7.250.75   <pending>    30303:31079/TCP             6s

		==> v1/ServiceAccount
		NAME                           SECRETS  AGE
		autonity-genesis-configurator  1        6s
		autonity-keys-generator        1        6s

		==> v1beta1/Deployment
		NAME             READY  UP-TO-DATE  AVAILABLE  AGE
		autonity-node-0  0/1    1           0          6s


		NOTES:
		======

		To get autonity autonity-node-0 account password type:
		    kubectl -n how-to-se-03 get secrets autonity-node-0 -o 'go-template={{index .data "password"}}' | base64 --decode; echo ""

		Get private key of autonity-node-0
		    kubectl -n how-to-se-03 get secrets autonity-node-0 -o 'go-template={{index .data "private_key"}}' | base64 --decode; echo ""

		Get address for autonity-node-0
		    kubectl -n how-to-se-03 get configmap autonity-node-0 -o jsonpath='{.data.address}'

		Get genesis.json
		    kubectl -n how-to-se-03 get configmaps genesis -o jsonpath='{.data.genesis}'

		Export genesis.yaml
		sudo snap install jq yq
		    kubectl -n how-to-se-03 get configmaps genesis -o jsonpath='{.data.genesis}' |jq '{genesis: .}' |yq r -

		Forward rpcapi autonity-node-0 to localhost
		    kubectl -n how-to-se-03 port-forward svc/autonity-node-0 8545:8545
		Forward wsapi autonity-node-0 to localhost
		    kubectl -n how-to-se-03 port-forward svc/autonity-node-0 8546:8546

		Get enode
		*********

		It can be some time until a Public IP is allocated

		    IP=$(kubectl -n how-to-se-03 get svc p2p-autonity-node-0 -o jsonpath="{.status.loadBalancer.ingress[*].ip}"); \
		    PUB_KEY=$(kubectl -n how-to-se-03 get configmap autonity-node-0 -o jsonpath="{.data.pub_key}"); \
		    PORT=$(kubectl -n how-to-se-03 get svc p2p-autonity-node-0 -o jsonpath="{.spec.ports[0].port}"); \
		    echo "enode://"${PUB_KEY}\@${IP}\:${PORT}

		HTTP(s)-RPC
		***********

		Get last block number
		    curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://localhost:8545

		Get Autonity Contract Address
		    curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"tendermint_getContractAddress","params":[],"id":1}' http://localhost:8545





5. Confirm the nodes are deployed successfully on the cluster:

	.. code:: bash

		helm status <node name>

	If there's a problem during deployment, you'll see a message like:


	.. code:: bash

		LAST DEPLOYED: Mon Nov  4 17:20:38 2019
		NAMESPACE: how-to-se
		STATUS: DEPLOYED

		RESOURCES:
		==> MISSING
		KIND                                                   NAME
		/v1, Resource=secrets                                  autonity-node-0
		/v1, Resource=configmaps                               autonity-tests
		/v1, Resource=configmaps                               autonity-node-0
		/v1, Resource=configmaps                               genesis
		/v1, Resource=configmaps                               genesis-template
		/v1, Resource=configmaps                               nginx-conf
		/v1, Resource=serviceaccounts                          autonity-genesis-configurator
		/v1, Resource=serviceaccounts                          autonity-keys-generator
		rbac.authorization.k8s.io/v1, Resource=roles           genesis-write
		rbac.authorization.k8s.io/v1, Resource=roles           secrets-write
		rbac.authorization.k8s.io/v1, Resource=rolebindings    secrets-write
		rbac.authorization.k8s.io/v1, Resource=rolebindings    genesis-write
		/v1, Resource=services                                 autonity-node-0
		/v1, Resource=services                                 p2p-autonity-node-0
		apps/v1beta1, Resource=deployments                     autonity-node-0
		batch/v1, Resource=jobs                                init-job02-genesis-configurator





Step 3: Define secure communication between nodes
-------------------------------------------------

To secure communications between nodes we need to define DNS records in GCP and add secure keys to the `genesis.yaml`.


.. Note:: Use the following steps to check the current status of the network. We'll also use these steps to make sure that the network is properly set up after we've completed the inititation steps

1. Get the information for the first Validator:

	.. code:: bash
		
		helm status <node name>

	This command returns the following:


	.. code:: bash

	  LAST DEPLOYED: Wed Nov  6 14:18:00 2019
	  NAMESPACE: how-to-se
	  STATUS: DEPLOYED

	  RESOURCES:
	  ==> v1/ConfigMap
	  NAME              DATA  AGE
	  autonity-node-0   2     14m
	  autonity-tests    1     14m
	  genesis           0     14m
	  genesis-template  1     14m
	  nginx-conf        1     14m

	  ==> v1/Job
	  NAME                             COMPLETIONS  DURATION  AGE
	  init-job02-genesis-configurator  0/1          14m       14m

	  ==> v1/Pod(related)
	  NAME                                   READY  STATUS    RESTARTS  AGE
	  autonity-node-0-67b5fdf8b-kqspj        0/2    Init:0/1  0         14m
	  init-job02-genesis-configurator-54qp6  1/1    Running   0         14m

	  ==> v1/Role
	  NAME           AGE
	  genesis-write  14m
	  secrets-write  14m

	  ==> v1/RoleBinding
	  NAME           AGE
	  genesis-write  14m
	  secrets-write  14m

	  ==> v1/Secret
	  NAME             TYPE    DATA  AGE
	  autonity-node-0  Opaque  2     14m

	  ==> v1/Service
	  NAME                 TYPE          CLUSTER-IP    EXTERNAL-IP    PORT(S)                     AGE
	  autonity-node-0      ClusterIP     10.7.246.119  <none>         8545/TCP,8546/TCP,9200/TCP  14m
	  p2p-autonity-node-0  LoadBalancer  10.7.251.42   35.230.150.24  30303:31449/TCP             14m

	  ==> v1/ServiceAccount
	  NAME                           SECRETS  AGE
	  autonity-genesis-configurator  1        14m
	  autonity-keys-generator        1        14m

	  ==> v1beta1/Deployment
	  NAME             READY  UP-TO-DATE  AVAILABLE  AGE
	  autonity-node-0  0/1    1           0          14m

	  NOTES:
	 

		To get autonity autonity-node-0 account password type:
		      kubectl -n how-to-se get secrets autonity-node-0 -o 'go-template={{index .data "password"}}' | base64 --decode; echo ""

		Get private key of autonity-node-0
		      kubectl -n how-to-se get secrets autonity-node-0 -o 'go-template={{index .data "private_key"}}' | base64 --decode; echo ""

		Get address for autonity-node-0
		      kubectl -n how-to-se get configmap autonity-node-0 -o jsonpath='{.data.address}'

		Get genesis.json
		      kubectl -n how-to-se get configmaps genesis -o jsonpath='{.data.genesis}'

		Export genesis.yaml
		sudo snap install jq yq
		      kubectl -n how-to-se get configmaps genesis -o jsonpath='{.data.genesis}' |jq '{genesis: .}' |yq r -

		Forward rpcapi autonity-node-0 to localhost
		      kubectl -n how-to-se port-forward svc/autonity-node-0 8545:8545
		Forward wsapi autonity-node-0 to localhost
		      kubectl -n how-to-se port-forward svc/autonity-node-0 8546:8546

		Get enode
		---------

		An IP address may take some time to be allocated:

		      IP=$(kubectl -n how-to-se get svc p2p-autonity-node-0 -o jsonpath="{.status.loadBalancer.ingress[*].ip}"); \
		      PUB_KEY=$(kubectl -n how-to-se get configmap autonity-node-0 -o jsonpath="{.data.pub_key}"); \
		      PORT=$(kubectl -n how-to-se get svc p2p-autonity-node-0 -o jsonpath="{.spec.ports[0].port}"); \
		      echo "enode://"${PUB_KEY}\@${IP}\:${PORT}


		HTTP(s)-RPC 
		Get last block number
		      curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://localhost:8545

		Get Autonity Contract Address
		      curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"tendermint_getContractAddress","params":[],"id":1}' http://localhost:8545


		As shown in the message above under the pod information:


		==> v1/Pod(related)
		NAME                                   READY  STATUS    RESTARTS  AGE
		autonity-node-0-67b5fdf8b-kqspj        0/2    Init:0/1  0         14m
		init-job02-genesis-configurator-54qp6  1/1    Running   0         14m


		There are 2 running services:
		- `autonity-node-0-67b5fdf8b-kqspj` which is the Autonity node and it is in `Init` status
		- `init-job02-genesis-configurator-54qp6` which is the peer discovery process trying to find the other nodes listed in the `gensis.yaml` file uploaded in the Step 2. This service is currently in `Running` status.



2. Check the status of the peer discovery:
	

	.. code:: bash

			kubectl -n <node namespace> logs <peer discovery service name>

	The service name in this example as shown above will be `init-job02-genesis-configurator-54qp6`


	This command returns the following message:

	.. code:: bash
			
		
			  2019-11-06 14:18:06 INFO     Trying to resolv peers: ['val-0-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com', 'val-1-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com', 'val-2-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com', 'val-3-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com', 'val-4-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com']
			  2019-11-06 14:18:06 INFO     Use Name Servers to resolve records: ['1.1.1.1', '8.8.8.8']
			  2019-11-06 14:18:06 WARNING  A record: None of DNS query names exist: val-0-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-0-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 WARNING  TXT record: None of DNS query names exist: val-0-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-0-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 WARNING  A record: None of DNS query names exist: val-1-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-1-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 WARNING  TXT record: None of DNS query names exist: val-1-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-1-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 WARNING  A record: None of DNS query names exist: val-2-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-2-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 WARNING  TXT record: None of DNS query names exist: val-2-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-2-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 WARNING  A record: None of DNS query names exist: val-3-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-3-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 WARNING  TXT record: None of DNS query names exist: val-3-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-3-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 WARNING  A record: None of DNS query names exist: val-4-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-4-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 WARNING  TXT record: None of DNS query names exist: val-4-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-4-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
			  2019-11-06 14:18:06 INFO     Fully resolved 0 fqdn records from 5


	As shown in the return message above, the node is trying to resolve the other network peers and the current status is listed in the following line:

	.. code:: bash

		2019-11-06 14:18:06 INFO     Fully resolved 0 fqdn records from 5

The above line mentions that at that point the node cannot resolve any of the DNSs mentioned in the `genesis.yaml` file.

3. Create the required DNSs through the GCP for all 5 nodes by setting the following information for each DNS:
  a. Create an A record which includes the node IP address
  b. Create a Txt record which includes the node port number and the node's public key value

	  To get the values, run:

		.. code:: bash

		 	IP=$(kubectl -n how-to-se get svc p2p-autonity-node-0 -o jsonpath="{.status.loadBalancer.ingress[*].ip}"); \
		    PUB_KEY=$(kubectl -n how-to-se get configmap autonity-node-0 -o jsonpath="{.data.pub_key}"); \
		    PORT=$(kubectl -n how-to-se get svc p2p-autonity-node-0 -o jsonpath="{.spec.ports[0].port}"); \
		    echo "enode://"${PUB_KEY}\@${IP}\:${PORT}


		Which returns:

		.. code:: bash
			
		  enode://fcd5c05d98846325f5578f825ed05fbd96ef073b8d45e88eb3e9cc298b92326d5c2c4d3b0492862ce9f142ba04a109ee726449f97f009e24be4e898b000dad62@35.230.150.24:30303


		The above message includes all the required information for the DNS creation as follows:


		    .. list-table:: 
		       :widths: 5 30 30 10
		       :header-rows: 1

		       * - Type
		         - Name
		         - Value
		         - TTL
		       * - A
		         - [node name].[cloud DNS domain]
		         - [node ip address]
		         - 1 min
		       * - TXT
		         - [node name].[cloud DNS domain]
		         - "p=[port number]; k=[node public key]"
		         - 1 min 


	.. Important:: Repeat this step for all 5 nodes

	When you've created all the node DNS records, the nodes will be resolved and they'll be able to communicate. The log process started in step 2 above will return the following message:

		.. code:: bash
			
			  2019-11-06 15:00:27 INFO     Fully resolved 5 fqdn records from 5
			  2019-11-06 15:00:27 INFO     All fqdn records was resolved successfully
			  2019-11-06 15:00:28 INFO     Generated genesis was written successfully to ConfigMap genesis
			  {
			    "alloc": {
			      "0x981dfa463bE3c247Fe311a05EeD4c67265204417": {
			        "balance": "0x200000000000000000000000000000000000000000000000000000000000000"
			      }
			    },
			    "coinbase": "0x0000000000000000000000000000000000000000",
			    "config": {
			      "autonityContract": {
			        "abi": "",
			        "bytecode": "",
			        "deployer": "0x0000000000000000000000000000000000000002",
			        "minGasPrice": 10000000000000,
			        "operator": "0xae223655126e514C9C80096d99765A98547247D3",
			        "users": [
			          {
			            "enode": "enode://fcd5c05d98846325f5578f825ed05fbd96ef073b8d45e88eb3e9cc298b92326d5c2c4d3b0492862ce9f142ba04a109ee726449f97f009e24be4e898b000dad62@35.230.150.24:30303",
			            "stake": 50000,
			            "type": "validator"
			          },
			          {
			            "enode": "enode://3f75bdc77aca9506b4a9e40bf1192a1e970a6d9fd3982466c36fc5a40a612e926faec7babfef8c94c56f47106bbc7e64d092e58ca9742602c2bbdb3142b8ef1b@35.246.127.20:30303",
			            "stake": 50000,
			            "type": "validator"
			          },
			          {
			            "enode": "enode://ceb3c1510539cbb8b292f9a0eb0daf31eff96ac4e29ea6cbde7c3a8371f34f28d2e60dda974d0a6a660b37b824a6711bd49dbc89094123835b3fba7702e8b4ee@35.246.42.226:30303",
			            "stake": 50000,
			            "type": "validator"
			          },
			          {
			            "enode": "enode://ba1e7faf46460311861bad1f3140e6e657f5a130ba77c369ccb870c886411da071de0cb886b45342f4842500bc017159ebafc82aff6c631bdc11a8b41db479b4@35.234.155.208:30303",
			            "stake": 50000,
			            "type": "validator"
			          },
			          {
			            "enode": "enode://06610fa67585744ed1d67a5f4213568944206fdcf51cf5dd540d549de5494a32c3c514de0ee01827a3315090eefc467c055714886b3d4b911a7b104e4d7340aa@34.89.46.114:30303",
			            "type": "participant"
			          }
			        ]
			      },
			      "byzantiumBlock": 0,
			      "chainId": 1489,
			      "constantinopleBlock": 0,
			      "eip150Block": 0,
			      "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
			      "eip155Block": 0,
			      "eip158Block": 0,
			      "homesteadBlock": 0,
			      "petersburgBlock": 0,
			      "tendermint": {
			        "block-period": 1,
			        "policy": 0
			      }
			    },
			    "difficulty": "0x1",
			    "gasLimit": "0x5F5E100",
			    "gasUsed": "0x0",
			    "mixHash": "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
			    "nonce": "0x0",
			    "number": "0x0",
			    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
			    "timestamp": "0x0"
			  }

	The returned message states that all nodes has been resolved and returns the updated network gensis file that can be used later to deploy any new nodes.


4. Verify the network active status by running the following commands. 

	For any node:

	  a. Run 

	  	.. code:: bash

	  		helm status <node name>

	  b. Copy and run the command to `# Forward rpcapi autonity-node-0 to localhost` 

		  .. code:: bash

		  		kubectl -n how-to-se-03 port-forward svc/autonity-node-0 8545:8545)

	  c. Open a terminal window and run the following command to get the last block number

	    	.. Code:: bash

	       		curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://localhost:8545

	  d. Repeat step c and compare the returned block number. The block number should be different to verify that the network is running and blocks are being mined.

