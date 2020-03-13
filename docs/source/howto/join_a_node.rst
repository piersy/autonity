Join your node to an Autonity network
==================================

We want to join our node to an existing Autonity test network. The test network is hosted on Google Cloud Platform.

Procedure: initialise gcloud and log in, join an existing project, and do some configuration. Then we're ready to use the Autonity network to transact.

	.. Note::

		All commands use a terminal session

1. Initialise your Google Cloud Platform instance:

	.. Code:: bash
		
		gcloud init

	Returns:

	.. Code:: bash
		
		Welcome! This command will take you through the configuration of gcloud.

		Your current configuration has been set to: [default]

		You can skip diagnostics next time by using the following flag:
		  gcloud init --skip-diagnostics

		Network diagnostic detects and fixes local network connection issues.
		Checking network connection...done.
		Reachability Check passed.
		Network diagnostic passed (1/1 checks passed).


2. Log in:

	.. Code:: bash
		
		You must log in to continue. Would you like to log in (Y/n)?  y

		Your browser has been opened to visit:

		 https://accounts.google.com/o/oauth2/auth?code_challenge=D4iMIv7el6rO9sNduy87I1BTKBn-37GkpWOfLehTTtY&prompt=select_account&code_challenge_method=S256&access_type=offline&redirect_uri=http%3A%2F%2Flocalhost%3A8085%2F&response_type=code&client_id=32555940559.apps.googleusercontent.com&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcloud-platform+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fappengine.admin+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcompute+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Faccounts.reauth

		You are logged in as: [max.croasdale@clearmatics.com].

3. Join a project:

	.. Code:: bash
		
		Pick cloud project to use:
		 [1] clearmatics-net-404101
		 [2] devops-tools-114459
		 [3] platformx-3816
		 [4] playground-for-raj
		 [5] Create a new project
		Please enter numeric choice or text value (must exactly match list
		item):  3



4. Create a Namespace using 'kubectl create':

	.. Code:: bash

		kubectl create namespace max-0-test

	Returns:

	.. Code:: bash

		namespace/max-0-test created

5. Deploy the node using Helm install by replacing the values for name and namespace:


	.. Code:: bash

			mc@clearmatics-3:~$ helm install --name max-0-test --namespace max-0-test charts-ose.clearmatics.com/autonity -f /Users/mc/work/test-network-deployment/genesis.yaml

	.. Note::

			name and namespace can be the same


	
	Returns:

	.. Code:: bash
	
		NAME:   max-0-test
		LAST DEPLOYED: Wed Feb 19 12:42:54 2020
		NAMESPACE: max-0-test
		STATUS: DEPLOYED

		RESOURCES:
		==> v1/ConfigMap
		NAME              DATA  AGE
		autonity-node-0   2     3s
		autonity-tests    1     3s
		genesis           0     3s
		genesis-template  1     3s
		nginx-conf        1     3s

		==> v1/Job
		NAME                             COMPLETIONS  DURATION  AGE
		init-job02-genesis-configurator  0/1          3s        3s

		==> v1/Pod(related)
		NAME                                   READY  STATUS    RESTARTS  AGE
		autonity-node-0-5d4777b6ff-gq9gw       0/2    Init:0/1  0         3s
		init-job02-genesis-configurator-92vb8  1/1    Running   0         3s

		==> v1/Role
		NAME           AGE
		genesis-write  3s
		secrets-write  3s

		==> v1/RoleBinding
		NAME           AGE
		genesis-write  3s
		secrets-write  3s

		==> v1/Secret
		NAME             TYPE    DATA  AGE
		autonity-node-0  Opaque  2     3s

		==> v1/Service
		NAME                 TYPE          CLUSTER-IP    EXTERNAL-IP  PORT(S)                     AGE
		autonity-node-0      ClusterIP     10.7.246.83   <none>       8545/TCP,8546/TCP,9200/TCP  3s
		p2p-autonity-node-0  LoadBalancer  10.7.252.206  <pending>    30303:30541/TCP             3s

		==> v1/ServiceAccount
		NAME                           SECRETS  AGE
		autonity-genesis-configurator  1        3s
		autonity-keys-generator        1        3s

		==> v1beta1/Deployment
		NAME             READY  UP-TO-DATE  AVAILABLE  AGE
		autonity-node-0  0/1    1           0          3s


		NOTES:
		==== Autonity ====

		# To get autonity autonity-node-0 account password type:
			    kubectl -n max-0-test get secrets autonity-node-0 -o 'go-template={{index .data "password"}}' | base64 --decode; echo ""

		# Get private key of autonity-node-0
				    kubectl -n max-0-test get secrets autonity-node-0 -o 'go-template={{index .data "private_key"}}' | base64 --decode; echo ""

		# Get address for autonity-node-0
				    kubectl -n max-0-test get configmap autonity-node-0 -o jsonpath='{.data.address}'

		# Get genesis.json
				    kubectl -n max-0-test get configmaps genesis -o jsonpath='{.data.genesis}'

		# Export genesis.yaml
		# sudo snap install jq yq
				    kubectl -n max-0-test get configmaps genesis -o jsonpath='{.data.genesis}' |jq '{genesis: .}' |yq r -
		# Forward rpcapi autonity-node-0 to localhost
				    kubectl -n max-0-test port-forward svc/autonity-node-0 8545:8545
		# Forward wsapi autonity-node-0 to localhost
				    kubectl -n max-0-test port-forward svc/autonity-node-0 8546:8546

		### Get enode:
		# It can take a time to wait until Public IP will allocated

			IP=$(kubectl -n max-0-test get svc p2p-autonity-node-0 -o jsonpath="{.status.loadBalancer.ingress[*].ip}"); \
			PUB_KEY=$(kubectl -n max-0-test get configmap autonity-node-0 -o jsonpath="{.data.pub_key}"); \
			PORT=$(kubectl -n max-0-test get svc p2p-autonity-node-0 -o jsonpath="{.spec.ports[0].port}"); \
			echo "enode://"${PUB_KEY}\@${IP}\:${PORT}

		###

		### HTTP(s)-RPC ###
		# Get last block number
			curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://localhost:8545

		# Get Autonity Contract Address
			 curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"tendermint_getContractAddress","params":[],"id":1}' http://localhost:8545


6. Get the node's enode value by executing the code:

	.. Code:: bash

		IP=$(kubectl -n max-0-test get svc p2p-autonity-node-0 -o jsonpath="{.status.loadBalancer.ingress[*].ip}"); \
		PUB_KEY=$(kubectl -n max-0-test get configmap autonity-node-0 -o jsonpath="{.data.pub_key}"); \
		PORT=$(kubectl -n max-0-test get svc p2p-autonity-node-0 -o jsonpath="{.spec.ports[0].port}"); \
		echo "enode://"${PUB_KEY}\@${IP}\:${PORT}

	Returns:

	.. Code:: bash

		enode://8777257ab3dcb4fb4f7aa30501432d68ffa209af1f184c024be23f08668cb9c38b08b13fad09e266bbf02c2238bc730bbd4929657fc9aadb1e9344716a611a8d@35.242.181.196:30303

	Send this `enode` value securely to the System Operator. The System Operator will add your node to the network and issue stake.

7. Set up port forwarding to interact with the network:

	.. Code:: bash

		mc@clearmatics-3:~$ kubectl -n max-0-test port-forward svc/autonity-node-0 8545:8545

	Returns:

	.. Code:: bash

		Forwarding from 127.0.0.1:8545 -> 8080
		Forwarding from [::1]:8545 -> 8080
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		^Cmc@clearmatics-3:~$ kubectl -n max-0-test port-forward svc/autonity-node-0 8545:8545
		Forwarding from 127.0.0.1:8545 -> 8080
		Forwarding from [::1]:8545 -> 8080
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545
		Handling connection for 8545

When you have received confirmation from the System Operator that your node has been added to the network you can transact on the network.
