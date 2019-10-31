rm -rf ibc-iris && rm -rf ibc-gaia

echo 12345678 | iris testnet -o ibc-iris --v 1 --chain-id iris --node-dir-prefix n
echo 12345678 | gaiad testnet -o ibc-gaia --v 1 --chain-id cosmos --node-dir-prefix n

sed -i '' 's/"leveldb"/"goleveldb"/g' ibc-iris/n0/iris/config/config.toml
sed -i '' 's/"leveldb"/"goleveldb"/g' ibc-gaia/n0/gaiad/config/config.toml
sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26556"#g' ibc-gaia/n0/gaiad/config/config.toml
sed -i '' 's#"tcp://0.0.0.0:26657"#"tcp://0.0.0.0:26557"#g' ibc-gaia/n0/gaiad/config/config.toml
sed -i '' 's#"localhost:6060"#"localhost:6061"#g' ibc-gaia/n0/gaiad/config/config.toml
sed -i '' 's#"tcp://127.0.0.1:26658"#"tcp://127.0.0.1:26558"#g' ibc-gaia/n0/gaiad/config/config.toml

sed -i '' 's/n0token/uiris/' ibc-iris/n0/iris/config/genesis.json
sed -i '' 's/n0token/uatom/' ibc-gaia/n0/gaiad/config/genesis.json

sed -i '' 's/timeout_commit = "5s"/timeout_commit = "2s"/' ibc-iris/n0/iris/config/config.toml
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "2s"/' ibc-gaia/n0/gaiad/config/config.toml

iriscli config --home ibc-iris/n0/iriscli/ chain-id iris
gaiacli config --home ibc-gaia/n0/gaiacli/ chain-id cosmos
iriscli config --home ibc-iris/n0/iriscli/ output json
gaiacli config --home ibc-gaia/n0/gaiacli/ output json
iriscli config --home ibc-iris/n0/iriscli/ node http://localhost:26657
gaiacli config --home ibc-gaia/n0/gaiacli/ node http://localhost:26557
