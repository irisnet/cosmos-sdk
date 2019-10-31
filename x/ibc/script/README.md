# TEST

## Install `iris` `gaia` `relayer`

Install `iris`

```bash
git clone https://github.com/irisnet/irishub.git
cd irishub
git checkout segue/ibc-handshake
go mod tidy && make install
```

Install `gaia`

```bash
git clone https://github.com/irisnet/gaia.git
cd gaia
git checkout segue/ibc-handshake
go mod tidy && make install
```

## Initiate

```bash
chmod 777 init.sh
chmod 777 handshake.sh
```

```bash
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
handshake.sh
```
