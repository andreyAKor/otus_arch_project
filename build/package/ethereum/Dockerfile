FROM ethereum/client-go:latest

RUN mkdir -p /var/log/otus_arch_project

ENTRYPOINT ["geth", "--syncmode=fast", "--rpc", "--rpcport=8545", "--rpcaddr=0.0.0.0", "--rpccorsdomain=*", "--rpcapi=eth,web3,personal", "--port=30303", "--allow-insecure-unlock", "--ropsten", "--maxpeers=10"]
