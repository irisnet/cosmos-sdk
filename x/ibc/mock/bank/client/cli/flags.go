package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagSrcPort     = "src-port"
	FlagSrcChannel  = "src-channel"
	FlagDenom       = "denom"
	FlagAmount      = "amount"
	FlagReceiver    = "receiver"
	FlagSource      = "source"
	FlagProofPath   = "proof-path"
	FlagProofHeight = "proof-height"
)

var (
	FsTransfer = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsTransfer.String(FlagSrcPort, "", "the source port ID")
	FsTransfer.String(FlagSrcChannel, "", "the source channel ID")
	FsTransfer.String(FlagDenom, "", "the denomination to be transferred")
	FsTransfer.String(FlagAmount, "", "the amount to be transferred")
	FsTransfer.String(FlagReceiver, "", "the recipient")
	FsTransfer.Bool(FlagSource, true, "indicate if the sender is the source chain of the token")
	FsTransfer.String(FlagProofPath, "", "the path of the proof file")
	FsTransfer.Uint64(FlagProofHeight, 0, "the block height at which the proof is generated")
}
