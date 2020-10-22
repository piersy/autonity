package tendermint

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
)

type oracle struct {
	lastHeader   *types.Header
	store        *messageStore
	committeeSet committee
	c            *bridge
	ba           *blockAwaiter
}

func newOracle(lh *types.Header, s *messageStore, cs committee, ba *blockAwaiter) *oracle {
	return &oracle{
		lastHeader:   lh,
		store:        s,
		committeeSet: cs,
		ba:           ba,
	}
}

func (o *oracle) FThresh(round int64) bool {
	return o.store.fail(round, o.lastHeader)
}

func (o *oracle) MatchingProposal(cm *algorithm.ConsensusMessage) *algorithm.ConsensusMessage {
	mp := o.store.matchingProposal(cm)
	if mp != nil {
		return mp.consensusMessage
	}
	return nil
}

func (o *oracle) PrecommitQThresh(round int64, value *algorithm.ValueID) bool {
	return o.store.precommitQuorum((*common.Hash)(value), round, o.lastHeader)
}

func (o *oracle) PrevoteQThresh(round int64, value *algorithm.ValueID) bool {
	return o.store.prevoteQuorum((*common.Hash)(value), round, o.lastHeader)
}

func (o *oracle) Proposer(round int64, nodeID algorithm.NodeID) bool {
	return o.committeeSet.GetProposer(round).Address == common.Address(nodeID)
}

func (o *oracle) Valid(value algorithm.ValueID) bool {
	return o.store.isValid(common.Hash(value))
}

func (o *oracle) Height() uint64 {
	return o.lastHeader.Number.Uint64()
}

func (o *oracle) Value() (algorithm.ValueID, error) {
	v, err := o.ba.value(o.lastHeader.Number)
	if err != nil {
		return [32]byte{}, err
	}
	// The tendermint is making a proposal, we need to ensure that we add the proposal block to the msg store, so that
	// it can be picked up in buildMessage.
	o.store.addValue(v.Hash(), v)
	return algorithm.ValueID(v.Hash()), nil
}
