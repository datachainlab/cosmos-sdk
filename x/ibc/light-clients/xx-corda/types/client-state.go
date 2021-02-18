package types

import (
	ics23 "github.com/confio/ics23/go"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
)

var _ exported.ClientState = (*ClientState)(nil)

const VersionNumber = 0
const VersionHeight = 1

func (*ClientState) ClientType() string {
	return "corda"
}

func (*ClientState) GetLatestHeight() exported.Height {
	return types.Height{
		VersionNumber: VersionNumber,
		VersionHeight: VersionHeight,
	}
}

func (*ClientState) IsFrozen() bool {
	panic("not implemented")
}

func (*ClientState) GetFrozenHeight() exported.Height {
	panic("not implemented")
}

func (*ClientState) Validate() error {
	return nil
}

func (*ClientState) GetProofSpecs() []*ics23.ProofSpec {
	panic("not implemented")
}

func (*ClientState) CheckHeaderAndUpdateState(sdk.Context, codec.BinaryMarshaler, sdk.KVStore, exported.Header) (exported.ClientState, exported.ConsensusState, error) {
	panic("not implemented")
}

func (*ClientState) CheckMisbehaviourAndUpdateState(sdk.Context, codec.BinaryMarshaler, sdk.KVStore, exported.Misbehaviour) (exported.ClientState, error) {
	panic("not implemented")
}

func (*ClientState) CheckProposedHeaderAndUpdateState(sdk.Context, codec.BinaryMarshaler, sdk.KVStore, exported.Header) (exported.ClientState, exported.ConsensusState, error) {
	panic("not implemented")
}

func (*ClientState) VerifyUpgrade(
	ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	newClient exported.ClientState,
	upgradeHeight exported.Height,
	proofUpgrade []byte,
) error {
	panic("not implemented")
}

func (*ClientState) ZeroCustomFields() exported.ClientState {
	panic("not implemented")
}

func (*ClientState) VerifyClientState(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	counterpartyClientIdentifier string,
	proof []byte,
	clientState exported.ClientState,
) error {
	return nil
}

func (*ClientState) VerifyClientConsensusState(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	counterpartyClientIdentifier string,
	consensusHeight exported.Height,
	prefix exported.Prefix,
	proof []byte,
	consensusState exported.ConsensusState,
) error {
	return nil
}

func (*ClientState) VerifyConnectionState(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	connectionID string,
	connectionEnd exported.ConnectionI,
) error {
	return nil
}

func (*ClientState) VerifyChannelState(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	channel exported.ChannelI,
) error {
	return nil
}

func (*ClientState) VerifyPacketCommitment(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	commitmentBytes []byte,
) error {
	return nil
}

func (*ClientState) VerifyPacketAcknowledgement(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	acknowledgement []byte,
) error {
	return nil
}

func (*ClientState) VerifyPacketReceiptAbsence(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
) error {
	return nil
}

func (*ClientState) VerifyNextSequenceRecv(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	nextSequenceRecv uint64,
) error {
	return nil
}
