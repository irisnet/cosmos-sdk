package rootmulti

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/merkle"
)

// RequireProof returns whether proof is required for the subpath.
func RequireProof(subpath string) bool {
	// XXX: create a better convention.
	// Currently, only when query subpath is "/key", will proof be included in
	// response. If there are some changes about proof building in iavlstore.go,
	// we must change code here to keep consistency with iavlStore#Query.
	return subpath == "/key"
}

//-----------------------------------------------------------------------------

// XXX: This should be managed by the rootMultiStore which may want to register
// more proof ops?
func DefaultProofRuntime() (prt *merkle.ProofRuntime) {
	prt = merkle.NewProofRuntime()
	prt.RegisterOpDecoder(merkle.ProofOpSimpleValue, merkle.SimpleValueOpDecoder)
	prt.RegisterOpDecoder(iavl.ProofOpIAVLValue, iavl.ValueOpDecoder)
	prt.RegisterOpDecoder(iavl.ProofOpIAVLAbsence, iavl.AbsenceOpDecoder)
	return
}
