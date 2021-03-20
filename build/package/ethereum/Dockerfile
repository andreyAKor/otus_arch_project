FROM ethereum/client-go:latest

ENTRYPOINT ["geth", "--syncmode=fast", "--rpc", "--rpcport=8545", "--rpcaddr=0.0.0.0", "--rpccorsdomain=*", "--rpcapi=eth,web3,personal", "--rpcvhosts=*", "--port=30303", "--allow-insecure-unlock", "--ropsten", "--maxpeers=10"]
