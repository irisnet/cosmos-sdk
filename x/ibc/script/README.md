# TEST

## Install `iris` `gaia` `relayer`

Install `iris`

```bash
git clone https://github.com/irisnet/irishub.git
cd irishub
git checkout segue/ibc-test
go mod tidy && make install
```

Install `gaia`

```bash
git clone https://github.com/irisnet/gaia.git
cd gaia
git checkout segue/ibc-test
go mod tidy && make install
```

replace your local cosmos-sdk in `go.mod`

```bash
replace github.com/cosmos/cosmos-sdk => /path-to-local/cosmos-sdk
```

## Initiate

```bash
chmod 777 init.sh
./init.sh
```

## Start `iris` and `gaia`

```bash
nohup iris --home ibc-iris/n0/iris start >ibc-iris.log &
nohup gaiad --home ibc-gaia/n0/gaiad start >ibc-gaia.log &

iris --home ibc-iris/n0/iris start
gaiad --home ibc-gaia/n0/gaiad start
```

## Handshake

Create `client`, `connection` handshake, and `channel` handshake

```bash
chmod 777 handshake.sh
handshake.sh
```
