package main

type ValueID [32]byte

var nilValue ValueID

type NodeID [20]byte

type Step uint8

const (
	Propose Step = iota
	Prevote
	Precommit
)

func (s Step) in(steps ...Step) bool {
	for _, step := range steps {
		if s == step {
			return true
		}
	}
	return false
}

type StateTransition uint8

const (
	NewHeight StateTransition = iota
	NewRound
	NewStep
)

type Result struct {
	Transition StateTransition
	Height     uint64
	Round      int64
	Message    *ConsensusMessage
}

type ConsensusMessage struct {
	MsgType    Step
	Height     uint64
	Round      int64
	Value      ValueID
	ValidRound int64
}

type Sender interface {
	Send(cm *ConsensusMessage)
}

type Algorithm struct {
	nodeId         NodeID
	height         uint64
	round          int64
	step           Step
	lockedRound    int64
	lockedValue    ValueID
	validRound     int64
	validValue     ValueID
	sender         Sender
	line34Executed bool
	line36Executed bool
	line47Executed bool
}

type Oracle interface {
	Valid(ValueID) bool
	MatchingProposal(*ConsensusMessage) *ConsensusMessage
	PrevoteQThresh(round int64, value *ValueID) bool
	PrecommitQThresh(round int64, value *ValueID) bool
	FThresh(round int64, value *ValueID) bool
	Proposer(NodeID) bool
	Value() ValueID
}

func (a *Algorithm) msg(msgType Step, value ValueID) *ConsensusMessage {
	cm := &ConsensusMessage{
		MsgType: msgType,
		Height:  a.height,
		Round:   a.round,
		Value:   value,
	}
	if msgType == Propose {
		cm.ValidRound = a.validRound
	}
	return msg
	//a.sender.Send(cm)
}

// Message sent + stepchange to propose (not sure we really care about the step change)
// Schedule timout propose (send a prevote for nil after some time)
func (a *Algorithm) StartRound(round int64, o Oracle) Result {
	a.round = round
	a.step = Propose
	var m *ConsensusMessage
	if o.Proposer(a.nodeId) {
		var v ValueID
		if a.validValue != nilValue {
			v = a.validValue
		} else {
			v = o.Value()
		}
		m = a.msg(Propose, v)
	} else {
		// Schedule on timout propose
	}
}

// MessageSent prevote + stepchange propose to prevote
// schedule TimeoutPrevote (send a precommit for nil after some time)
// MessageSent precommit + stepchange prevote to precommit
// schedule TimeoutPrecommit (start round, round +1)
// Move to next height

func (a *Algorithm) ReceiveMessage(cm *ConsensusMessage, o Oracle) {

	r := a.round
	s := a.step
	t := cm.MsgType

	// look up matching proposal, in the case of a message with msgType
	// proposal the matching proposal is the message.
	p := o.MatchingProposal(cm)

	// Some of the checks in these upon conditions are omitted because they have alrady been checked.
	//
	// - We do not check height because we only execute this code when the
	// message height matches the current height.
	//
	// - We do not check whether the message comes from a proposer since this
	// is checkded before calling this method and we do not process proposals
	// from non proposers.

	// Line 22
	if t.in(Propose) && cm.Round == r && cm.ValidRound == -1 && s == Propose {
		if o.Valid(cm.Value) && a.lockedRound == -1 || a.lockedValue == cm.Value {
			a.msg(Prevote, cm.Value)
		} else {
			a.msg(Prevote, nilValue)
		}
		a.step = Prevote
		s = Prevote
	}

	// Line 28
	if t.in(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQThresh(p.ValidRound, &p.Value) && s == Propose && (p.ValidRound >= 0 && p.ValidRound < r) {
		if o.Valid(p.Value) && (a.lockedRound <= p.ValidRound || a.lockedValue == p.Value) {
			a.msg(Prevote, p.Value)
		} else {
			a.msg(Prevote, nilValue)
		}
		a.step = Prevote
		s = Prevote
	}

	// Line 34
	if t.in(Prevote) && cm.Round == r && o.PrevoteQThresh(r, nil) && s == Prevote && !a.line34Executed {
		//c.prevoteTimeout.scheduleTimeout(c.timeoutPrevote(r), r, h, c.onTimeoutPrecommit)
	}

	// Line 36
	if t.in(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQThresh(r, &p.Value) && o.Valid(p.Value) && s >= Prevote && !a.line36Executed {
		if s == Prevote {
			a.lockedValue = p.Value
			a.lockedRound = r
			a.msg(Precommit, p.Value)
			s = Precommit // TODO set steps in all situations where we set the steps
			a.step = Precommit
		}
		a.validValue = p.Value
		a.validRound = r
	}

	// Line 44
	if t.in(Prevote) && cm.Round == r && o.PrevoteQThresh(r, &nilValue) && s == Prevote {
		a.msg(Precommit, nilValue)
		s = Precommit
		a.step = Precommit
	}

	// Line 47
	if t.in(Precommit) && cm.Round == r && o.PrecommitQThresh(r, nil) && !a.line47Executed {
		//c.precommitTimeout.scheduleTimeout(c.timeoutPrecommit(r), r, h, c.onTimeoutPrecommit) // TODO handle the timers
	}

	// Line 49
	if t.in(Propose, Precommit) && p != nil && o.PrecommitQThresh(p.Round, &p.Value) {
		if o.Valid(p.Value) {
			// TODO commit here commit(p.Value)
			a.height++
			a.lockedRound = -1
			a.lockedValue = nilValue
			a.validRound = -1
			a.validValue = nilValue
		}
		a.StartRound(0)

		// Not quite sure how to start the round nicely
		// need to ensure that we don't stack overflow in the case that the
		// next height messages are sufficient for consensus when we
		// process them and so on and so on. So I need to set the start
		// round states and then queue the messages for processing. And I
		// need to ensure that I get a list of messages to process in an
		// atomic step from the msg cache so that I don't end up trying to
		// process the same message twice.
	}

	// Line 55
	if cm.Round > r && o.FThresh(cm.Round, nil) {
		// TODO account for the fact that many rounds can be skipped here.  so
		// what happens to the old round messages? We don't process them, but
		// we can't remove them from the cache because they may be used in this
		// round. in the conditon at line 28. This means that we only should
		// clean the message cache when there is a height change, clearing out
		// all messages for the height.
		a.StartRound(cm.Round)
	}
}

func (a *Algorithm) SendMessage() {
}