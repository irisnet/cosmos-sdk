# Software Upgrade Implement

## Two types of upgrade
1. Add a new module
2. Change an old module 

## Main Idea
One module should only corresponds to a KVStore according to the `KVStoreKey`. There are two important design.
 
* If a KVStore can be read, it means that it can participate in apphash computing（completely different from the apphash computing）and its module can handle the `QueryMsgs` 
* If a KVStore can be written, it means that its module can handle the `msgs`. 

So we add a lifetime parameter to the  `KVStoreKey`  to control when the KVStore can be read or written. Lifetime parameter is only modified by the governance and it means that after the community agrees to pass a software-upgrade proposal, if the majority completes the upgrading process, the proposal could modify the lifetime to start the new version.

See below for details.
 
## KVStoreKey
Currently, the `KVStoreKey` in CosmosSDK v0.21.0 is

```go
type KVStoreKey struct {
  name string
}
```

Now we add `start` and `height` to `KVStoreKey`

```go
type KVStoreKey struct {
    name  string 
    start int64
    end   int64
}
```
When `start<= lastheight`,  `KVStore` can be read and written. When `lastheight>=end`, `KVStore` can't be written. `start` is always no more than `end`. 

All the `KVStoreKey` are stored in `Global Paramstore`.
`start` and `end` can only be changed by the governance.

## Design
### App Init
```go
	var app = &GaiaApp{
...
		keyMain:          sdk.NewKVStoreKey("main",0,math.MaxInt64),
		keyAccount:       sdk.NewKVStoreKey("acc",0,math.MaxInt64),
		keyStake:         sdk.NewKVStoreKey("stake",0,math.MaxInt64),
		keyGov:           sdk.NewKVStoreKey("gov",0,math.MaxInt64),
		keyParams:        sdk.NewKVStoreKey("params",0,math.MaxInt64),
		keyNew:           sdk.NewKVStoreKey("New",math.MaxInt64),//later verison
	}
	app.loadGlobalKey()// update the keys information
```

For genesis version, the `start` of all `KVStoreKey`  initialize to 0 like `main`,`acc` ... For later versions, if we add new module, the `start` initialize to $ 2^{64}$ like `keyNew`. `end` always initialize to $ 2^{64}$.

When application is initialized, the `KVStoreKey` in `Global Paramstore` will be assgined to all above keys according to the same `name`. If `key` doesn't exist in `Global Paramstore`，`Global Paramstore` will add this `KVStoreKey`.

### Gov

`SoftwareUpgradeProposal` is split into  `AddModuleProposal` and `ChangeModuleProposal`.

```go
type AddModuleProposal struct{
  name string 
  height int64
}

```
* `name`: the name of the module added
* `height`: the height of counting signal proportion

```go
type ChangeModuleProposal struct{
  oldname string
  newname string
  height int64
}
```
* `oldname`: the name of the module before changed
* `newname`: the name of the module after changed
* `height`: the height of counting signal proportion


```go
type MsgSwitch struct {
	ProposalID int64          `json:"proposalID"` // ID of the proposal
	validator  sdk.AccAddress `json:"depositer"`  // Address of the validator
}
```
This message is used to signal that my local software has been updated and ready to switch to the new version software.

## Upgrade Workflow

### Add a New Module
For example，we want to add a new module `ABC` and start to be used at  `height_ABC`.

1. (gov) Submit the `AddModuleProposal` to add `ABC`. Proposal contains the name `ABC` and `height_ABC`.
2. After the proposal is accepted, every validator start to download the new version, then run it and send the `MsgSwitch` message. Now the `start` in the KVStoreKey of the `ABC` module is the inital value($2^{64}$).So the `start` of `ABC` module is bigger than `lastheight`, it can't be read and written. 
3. When `lastheight` reach to the `height_ABC`, govenance will count if more than two-thirds of validators have updated the software according to the `MsgSwitch` messages. If `true`, the govenance will set the `start` in the KVStoreKey of the `ABC` module to the `height_ABC`  and the app switch to the new version of the software. But the validator who doesn't update to the version will report an error，because the old version don't have the `ABC` module. They will be slashed, as they are out of the consensus. They should download the new version to continue.

### Change an Old Module

 For example, we change the module `ABC` to `ABC*` at `height_ABC*`. `ABC*` is similiar with `ABC`, but `ABC*` is more complete. Our approach is to keep the old module handling the`msgs` before the `height_ABC*` and add the new module handling the `msgs` after the `hegiht_ABC`.

1. (gov) Submit the `ChangeModuleProposal` to change the module `ABC`. Proposal contains the old name `ABC`, the new name `ABC*` and the `height_ABC*` that we start to change the old module to the new module. 


2. After the proposal is accepted, every validator start to download new version, run it and send the `MsgSwitch`. 

 
3. When `lastheight` reach to the `height_ABC*`, govenance will count if more than two-thirds of validators have updated the software. If `true`, the govenance will do three things.   
   1. Set the `end` in the KVStoreKey of the `ABC` module to `height_ABC*`. (In fact, we can alse query the storage in old module `ABC`,but the module will not handle the message to change th ABCI status after the `height_ABC*`) 
   2. Start the new module `ABC*` by setting the `start` in the KVStoreKey of the `ABC*` module to the `height_ABC*`. The module `ABC*` will handle the message to change th ABCI status after the `height_ABC*`
   3. Transfer the useful storage from the old  to the new module.

4. Finally, the app switch to the new version of the software. 
Other thing is same as above.



 
