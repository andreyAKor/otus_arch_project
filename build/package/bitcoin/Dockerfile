FROM kylemanna/bitcoind:latest

COPY ./configs/bitcoin.conf /bitcoin/.bitcoin/bitcoin.conf

ENTRYPOINT ["bitcoind", "-conf=/bitcoin/.bitcoin/bitcoin.conf"]
