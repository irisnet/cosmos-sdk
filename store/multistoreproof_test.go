package store

import (
	"testing"
	"encoding/hex"
	"github.com/tendermint/iavl"
	"github.com/stretchr/testify/assert"
)

func TestVerifyProofForMultiStore(t *testing.T) {
	appHash,_ := hex.DecodeString("77c14d4f0f4ddbd2bf03f51fc81aa385f51f0fd3")

	storeName := "acc"
	substoreRootHash,_ := hex.DecodeString("ea5d468431015c2cd6295e9a0bb1fc0e49033828")

	var multiStoreCommitInfo []iavl.SubstoreCommitID

	multiStoreCommitInfo=append(multiStoreCommitInfo,iavl.SubstoreCommitID{
		Name:"ibc",
		Version:963,
		CommitHash:nil,
	})

	stakeRootHash,_ := hex.DecodeString("6f104e1d5884a9afd49ae102b54175264c351571")
	multiStoreCommitInfo=append(multiStoreCommitInfo,iavl.SubstoreCommitID{
		Name:"stake",
		Version:963,
		CommitHash:stakeRootHash,
	})

	slashingRootHash,_ := hex.DecodeString("7400e19b1eb05e19e108dec7eb6692b0a2e6dd43")
	multiStoreCommitInfo=append(multiStoreCommitInfo,iavl.SubstoreCommitID{
		Name:"slashing",
		Version:963,
		CommitHash:slashingRootHash,
	})

	govRootHash,_ := hex.DecodeString("62c171bb022e47d1f745608ff749e676dbd25f78")
	multiStoreCommitInfo=append(multiStoreCommitInfo,iavl.SubstoreCommitID{
		Name:"gov",
		Version:963,
		CommitHash:govRootHash,
	})

	multiStoreCommitInfo=append(multiStoreCommitInfo,iavl.SubstoreCommitID{
		Name:"main",
		Version:963,
		CommitHash:nil,
	})

	accRootHash,_ := hex.DecodeString("ea5d468431015c2cd6295e9a0bb1fc0e49033828")
	multiStoreCommitInfo=append(multiStoreCommitInfo,iavl.SubstoreCommitID{
		Name:"acc",
		Version:963,
		CommitHash:accRootHash,
	})

	err :=  VerifyProofForMultiStore(storeName, substoreRootHash, multiStoreCommitInfo, appHash)
	assert.Nil(t, err)

	appHash,_ = hex.DecodeString("88c14d4f0f4ddbd2bf03f51fc81aa385f51f0fd3")

	err =  VerifyProofForMultiStore(storeName, substoreRootHash, multiStoreCommitInfo, appHash)
	assert.Error(t,err,"appHash doesn't match to the merkle root of multiStoreCommitInfo")
}