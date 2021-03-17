FROM kylemanna/bitcoind:1.0.0

COPY ./configs/bitcoin.conf /bitcoin/.bitcoin/bitcoin.conf

#RUN mkdir -p /var/log/otus_arch_project

ENTRYPOINT ["bitcoind", "-conf=/bitcoin/.bitcoin/bitcoin.conf", "-zapwallettxes"]
