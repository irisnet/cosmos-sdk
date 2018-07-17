package store

import (
	"github.com/tendermint/tendermint/crypto/merkle"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/iavl"
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func VerifyMultiStoreCommitInfo(storeName string, multiStoreCommitInfo []iavl.SubstoreCommitID, appHash []byte) ([]byte, error) {
	var substoreCommitHash []byte
	var kvPairs cmn.KVPairs
	for _,multiStoreCommitID := range multiStoreCommitInfo {

		if multiStoreCommitID.Name == storeName {
			substoreCommitHash = multiStoreCommitID.CommitHash;
		}

		kHash := []byte(multiStoreCommitID.Name)
		storeInfo := storeInfo{
			Core:storeCore{
				CommitID:sdk.CommitID{
					Version: multiStoreCommitID.Version,
					Hash: multiStoreCommitID.CommitHash,
				},
			},
		}

		kvPairs = append(kvPairs, cmn.KVPair{
			Key:   kHash,
			Value: storeInfo.Hash(),
		})
	}
	if len(substoreCommitHash) == 0 {
		return nil, cmn.NewError("failed to get substore root commit hash by store name")
	}
	if kvPairs == nil {
		return nil, cmn.NewError("failed to extract information from multiStoreCommitInfo")
	}
	//sort the kvPair list
	kvPairs.Sort()

	//Rebuild simple merkle hash tree
	var hashList [][]byte
	for _, kvPair := range kvPairs {
		hashResult := merkle.SimpleHashFromTwoHashes(kvPair.Key,kvPair.Value)
		hashList=append(hashList,hashResult)
	}

	if !bytes.Equal(appHash,simpleHashFromHashes(hashList)){
		return nil, cmn.NewError("The merkle root of multiStoreCommitInfo doesn't equal to appHash")
	}
	return substoreCommitHash, nil
}

func VerifyRangeProof(key, value []byte, substoreCommitHash []byte, rangeProof *iavl.RangeProof) (error){

	// Validate the proof to ensure data integrity.
	err := rangeProof.Verify(substoreCommitHash)
	if err != nil {
		return  errors.Wrap(err, "proof root hash doesn't equal to substore commit root hash")
	}

	if len(value) != 0 {
		// Validate existence proof
		err = rangeProof.VerifyItem(key, value)
		if err != nil {
			return  errors.Wrap(err, "failed in existence verification")
		}
	} else {
		// Validate absence proof
		err = rangeProof.VerifyAbsence(key)
		if err != nil {
			return  errors.Wrap(err, "failed in absence verification")
		}
	}

	return nil
}

func simpleHashFromHashes(hashes [][]byte) []byte {
	// Recursive impl.
	switch len(hashes) {
	case 0:
		return nil
	case 1:
		return hashes[0]
	default:
		left := simpleHashFromHashes(hashes[:(len(hashes)+1)/2])
		right := simpleHashFromHashes(hashes[(len(hashes)+1)/2:])
		return merkle.SimpleHashFromTwoHashes(left,right)
	}
}
