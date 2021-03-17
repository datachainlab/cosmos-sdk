package types

import (
	context "context"
	"log"

	ics23 "github.com/confio/ics23/go"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	connectiontypes "github.com/cosmos/cosmos-sdk/x/ibc/core/03-connection/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/23-commitment/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ exported.ClientState = (*ClientState)(nil)

type lightClient struct {
	conn *grpc.ClientConn
	LightClientClient
}

func (lc lightClient) mustClose() {
	if err := lc.conn.Close(); err != nil {
		log.Fatalf("failed to close gRPC connection: %v", err)
	}
}

func connectLightclientd() lightClient {
	// connect
	conn, err := grpc.Dial("localhost:60000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to dial to Corda-IBC lightclientd: %v", err)
	}

	// create gRPC client wrapping the connection
	cli := NewLightClientClient(conn)

	return lightClient{
		conn:              conn,
		LightClientClient: cli,
	}
}

func makeState(clientState *ClientState, cdc codec.BinaryMarshaler, store sdk.KVStore) *State {
	it := store.Iterator(nil, nil)
	defer it.Close()
	if !it.Valid() {
		log.Print("no consensus state")
	}
	consensusState := clienttypes.MustUnmarshalConsensusState(cdc, it.Value()).(*ConsensusState)
	it.Next()
	if it.Valid() {
		log.Print("too many consensus states")
	}
	return &State{
		ClientState:    clientState,
		ConsensusState: consensusState,
	}
}
func makeStateWithoutConsensusState(clientState *ClientState) *State {
	return &State{
		ClientState:    clientState,
		ConsensusState: nil,
	}
}

func (cs *ClientState) ClientType() string {
	lc := connectLightclientd()
	defer lc.mustClose()

	res, err := lc.ClientType(context.TODO(), &ClientTypeRequest{
		makeStateWithoutConsensusState(cs),
	})
	if err != nil {
		log.Fatalf("failed to call gRPC function ClientType: %v", err)
	}
	return res.ClientType
}

func (cs *ClientState) GetLatestHeight() exported.Height {
	lc := connectLightclientd()
	defer lc.mustClose()

	res, err := lc.GetLatestHeight(context.TODO(), &GetLatestHeightRequest{
		makeStateWithoutConsensusState(cs),
	})
	if err != nil {
		log.Fatalf("failed to call gRPC function GetLatestHeight: %v", err)
	}
	return res.Height
}

func (cs *ClientState) IsFrozen() bool {
	lc := connectLightclientd()
	defer lc.mustClose()

	res, err := lc.IsFrozen(context.TODO(), &IsFrozenRequest{
		makeStateWithoutConsensusState(cs),
	})
	if err != nil {
		log.Fatalf("failed to call gRPC function IsFrozen: %v", err)
	}
	return res.IsFrozen
}

func (cs *ClientState) GetFrozenHeight() exported.Height {
	lc := connectLightclientd()
	defer lc.mustClose()

	res, err := lc.GetFrozenHeight(context.TODO(), &GetFrozenHeightRequest{
		makeStateWithoutConsensusState(cs),
	})
	if err != nil {
		log.Fatalf("failed to call gRPC function GetFrozenHeight: %v", err)
	}
	return res.Height
}

func (cs *ClientState) Validate() error {
	lc := connectLightclientd()
	defer lc.mustClose()

	_, err := lc.Validate(context.TODO(), &ValidateRequest{
		makeStateWithoutConsensusState(cs),
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function Validate: %v", err)
	}
	return err
}

func (cs *ClientState) GetProofSpecs() []*ics23.ProofSpec {
	lc := connectLightclientd()
	defer lc.mustClose()

	res, err := lc.GetProofSpecs(context.TODO(), &GetProofSpecsRequest{
		makeStateWithoutConsensusState(cs),
	})
	if err != nil {
		log.Fatalf("failed to call gRPC function GetProofSpecs: %v", err)
	}
	return res.ProofSpecs
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

func (cs *ClientState) VerifyUpgrade(
	ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	newClient exported.ClientState,
	upgradeHeight exported.Height,
	proofUpgrade []byte,
) error {
	lc := connectLightclientd()
	defer lc.mustClose()

	_, err := lc.VerifyUpgrade(context.TODO(), &VerifyUpgradeRequest{
		State:         makeState(cs, cdc, store),
		NewClient:     newClient.(*ClientState),
		UpgradeHeight: upgradeHeight.(*clienttypes.Height),
		ProofUpgrade:  proofUpgrade,
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function VerifyUpgrade: %v", err)
	}
	return err
}

func (cs *ClientState) ZeroCustomFields() exported.ClientState {
	lc := connectLightclientd()
	defer lc.mustClose()

	res, err := lc.ZeroCustomFields(context.TODO(), &ZeroCustomFieldsRequest{
		makeStateWithoutConsensusState(cs),
	})
	if err != nil {
		log.Fatalf("failed to call gRPC function ZeroCustomFields: %v", err)
	}
	return res.ClientState
}

func (cs *ClientState) VerifyClientState(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	counterpartyClientIdentifier string,
	proof []byte,
	clientState exported.ClientState,
) error {
	lc := connectLightclientd()
	defer lc.mustClose()

	any, err := codectypes.NewAnyWithValue(clientState.(proto.Message))
	if err != nil {
		log.Fatalf("failed to make Any from exported.ClientState: %v", err)
	}
	_, err = lc.VerifyClientState(context.TODO(), &VerifyClientStateRequest{
		State:                        makeState(cs, cdc, store),
		Height:                       height.(*clienttypes.Height),
		Prefix:                       prefix.(*commitmenttypes.MerklePrefix),
		CounterpartyClientIdentifier: counterpartyClientIdentifier,
		Proof:                        proof,
		ClientState:                  any,
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function VerifyClient: %v", err)
	}
	return err
}

func (cs *ClientState) VerifyClientConsensusState(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	counterpartyClientIdentifier string,
	consensusHeight exported.Height,
	prefix exported.Prefix,
	proof []byte,
	consensusState exported.ConsensusState,
) error {
	lc := connectLightclientd()
	defer lc.mustClose()

	any, err := codectypes.NewAnyWithValue(consensusState.(proto.Message))
	if err != nil {
		log.Fatalf("failed to make Any from exported.ClientState: %v", err)
	}
	_, err = lc.VerifyClientConsensusState(context.TODO(), &VerifyClientConsensusStateRequest{
		State:                        makeState(cs, cdc, store),
		Height:                       height.(*clienttypes.Height),
		CounterpartyClientIdentifier: counterpartyClientIdentifier,
		ConsensusHeight:              consensusHeight.(*clienttypes.Height),
		Prefix:                       prefix.(*commitmenttypes.MerklePrefix),
		Proof:                        proof,
		ConsensusState:               any,
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function VerifyClientConsensusState: %v", err)
	}
	return err
}

func (cs *ClientState) VerifyConnectionState(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	connectionID string,
	connectionEnd exported.ConnectionI,
) error {
	lc := connectLightclientd()
	defer lc.mustClose()

	_, err := lc.VerifyConnectionState(context.TODO(), &VerifyConnectionStateRequest{
		State:         makeState(cs, cdc, store),
		Height:        height.(*clienttypes.Height),
		Prefix:        prefix.(*commitmenttypes.MerklePrefix),
		Proof:         proof,
		ConnectionId:  connectionID,
		ConnectionEnd: connectionEnd.(*connectiontypes.ConnectionEnd),
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function VerifyConnectionState: %v", err)
	}
	return err
}

func (cs *ClientState) VerifyChannelState(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	channel exported.ChannelI,
) error {
	lc := connectLightclientd()
	defer lc.mustClose()

	_, err := lc.VerifyChannelState(context.TODO(), &VerifyChannelStateRequest{
		State:     makeState(cs, cdc, store),
		Height:    height.(*clienttypes.Height),
		Prefix:    prefix.(*commitmenttypes.MerklePrefix),
		Proof:     proof,
		PortId:    portID,
		ChannelId: channelID,
		Channel:   channel.(*channeltypes.Channel),
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function VerifyChannelState: %v", err)
	}
	return err
}

func (cs *ClientState) VerifyPacketCommitment(
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
	lc := connectLightclientd()
	defer lc.mustClose()

	_, err := lc.VerifyPacketCommitment(context.TODO(), &VerifyPacketCommitmentRequest{
		State:           makeState(cs, cdc, store),
		Height:          height.(*clienttypes.Height),
		Prefix:          prefix.(*commitmenttypes.MerklePrefix),
		Proof:           proof,
		PortId:          portID,
		ChannelId:       channelID,
		Sequence:        sequence,
		CommitmentBytes: commitmentBytes,
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function VerifyPacketCommitment: %v", err)
	}
	return err
}

func (cs *ClientState) VerifyPacketAcknowledgement(
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
	lc := connectLightclientd()
	defer lc.mustClose()

	_, err := lc.VerifyPacketAcknowledgement(context.TODO(), &VerifyPacketAcknowledgementRequest{
		State:           makeState(cs, cdc, store),
		Height:          height.(*clienttypes.Height),
		Prefix:          prefix.(*commitmenttypes.MerklePrefix),
		Proof:           proof,
		PortId:          portID,
		ChannelId:       channelID,
		Sequence:        sequence,
		Acknowledgement: acknowledgement,
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function VerifyPacketAcknowledgement: %v", err)
	}
	return err
}

func (cs *ClientState) VerifyPacketReceiptAbsence(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
) error {
	lc := connectLightclientd()
	defer lc.mustClose()

	_, err := lc.VerifyPacketReceiptAbsence(context.TODO(), &VerifyPacketReceiptAbsenceRequest{
		State:     makeState(cs, cdc, store),
		Height:    height.(*clienttypes.Height),
		Prefix:    prefix.(*commitmenttypes.MerklePrefix),
		Proof:     proof,
		PortId:    portID,
		ChannelId: channelID,
		Sequence:  sequence,
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function VerifyPacketReceiptAbsence: %v", err)
	}
	return err
}

func (cs *ClientState) VerifyNextSequenceRecv(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	nextSequenceRecv uint64,
) error {
	lc := connectLightclientd()
	defer lc.mustClose()

	_, err := lc.VerifyNextSequenceRecv(context.TODO(), &VerifyNextSequenceRecvRequest{
		State:            makeState(cs, cdc, store),
		Height:           height.(*clienttypes.Height),
		Prefix:           prefix.(*commitmenttypes.MerklePrefix),
		Proof:            proof,
		PortId:           portID,
		ChannelId:        channelID,
		NextSequenceRecv: nextSequenceRecv,
	})
	switch status.Convert(err).Code() {
	case codes.OK:
	case codes.Unknown:
	default:
		log.Fatalf("failed to call gRPC function VerifyNextSequenceRecv: %v", err)
	}
	return err
}
