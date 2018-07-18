# Software Upgrade Implement

## Two types of upgrade
1. Add the new module
2. Change the old module 

## KVstoreKey
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
When `start<= lastheight`,this `KVstore` not only can be read and written (the module can handle the `msgs`), but also participate in `CommitID` computing. When `lastheight>=end`, this `KVstore` can't be written and its module doesn't handle the `msg`. `start` is always no more than `end`. 

For genesis version, the `start` of all `KVstoreKey`  initialize to 0. For later versions, if we add new module, the `start` of it initialize to $ 2^{64}$. `end` always initialize to $ 2^{64}$.

`start` and `end` can only be changed by the governance.

## Design
### App Init
```go
	var app = &GaiaApp{
...
		keyMain:          sdk.NewKVStoreKey("main",0),
		keyAccount:       sdk.NewKVStoreKey("acc",0),
		keyStake:         sdk.NewKVStoreKey("stake",0),
		keyGov:           sdk.NewKVStoreKey("gov",0),
		keyParams:        sdk.NewKVStoreKey("params",0),
		keyNew: sdk.NewKVStoreKey("New",math.MaxInt64),//later verison
	}
	app.loadGlobalKey()// update the keys information
```

Then the `keys` in `Global Paramstore` will be assgined to all above keys according to the same `name`. If `key` doesn't exist in `Global Paramstore`，`Global Paramstore` will add this `key`.

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

## Update Workflow

### Add New Modules
For example，we want to add a new module `ABC` and start to be used at  `height_ABC`.

1. (gov) Submit the `AddModuleProposal` to add `ABC`. Proposal contains the name `ABC` and `height_ABC`.
2. After the proposal is accepted, every validator start to download the new version, then run it and send the `MsgSwitch` message. Now the `start` in the KVStoreKey of the `ABC` module is the inital value($2^{64}$).So the `start` of `ABC` module is bigger than `lastheight`, it can't be read and written. 
3. When `lastheight` reach to the `height_ABC`, govenance will count if more than two-thirds of validators have updated the software according to the `MsgSwitch` messages. If `true`, the govenance will set the `start` in the KVStoreKey of the `ABC` module to the `height_ABC`  and the app switch to the new version of the software. But the validator who doesn't update to the version will report an error，because the old version don't have the `ABC` module. They will be slashed, as they are out of the consensus. They should download the new version to continue.

### Change the Module

In fact, our method is  to keep the old module but replace the old module with the new module. So previous `msgs`  are handled by the former and subsequent `msgs` are handled by the latter. For example, we change the module `ABC` to `ABC*` at `height_ABC*`

1. (gov) Submit the `ChangeModuleProposal` to change the module `ABC`. Proposal contains the old name `ABC`, the new name `ABC*` and the `height` that we can start to change the old module to the new module. 


2. After the proposal is accepted, every validator start to download new version, run it and send the `MsgSwitch`. But the `start` of `ABC` module is bigger than `lastheight`,so it can't be read and written.
 
3. When `lastheight` reach to the `height_ABC*`, govenance will count if more than two-thirds of validators have updated the software. If `true`, the govenance will do three things.   
   1. Delete the old module `ABC` by setting the `end` of `ABC` to `height_ABC*`. (In fact, we can alse query the storage in old module `ABC`,but the module will not handle the message to change th ABCI status after the `height_ABC*`) 
   2. Start the new module `ABC*` by setting the `start` of `ABC*` to the `height_ABC*`. The module `ABC*` will handle the message to change th ABCI status after the `height_ABC*`
   3. Transfer the useful storage from the old  to the new module.

4. Finally, the app switch to the new version of the software. But the validator which don't update the version will report an error，because the old version don't have the `ABC*`. They will be slashed, because they are out of the consensus.




 
