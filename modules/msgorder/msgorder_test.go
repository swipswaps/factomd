package msgorder_test

import (
	"github.com/FactomProject/factomd/modules/msgorder"
	"testing"

	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/modules/leader"
	"github.com/stretchr/testify/assert"

	"github.com/FactomProject/factom"
	. "github.com/FactomProject/factomd/testHelper"
)

func TestMsgOrderList(t *testing.T) {
	state := CreateEmptyTestState()
	ml := msgorder.NewOrderedMessageList()

	extIDs := [][]byte{[]byte("foo"), []byte("bar")}

	e := factom.Entry{
		ChainID: factom.ChainIDFromFields(extIDs),
		ExtIDs:  extIDs,
		Content: []byte("Hello World!"),
	}
	chain := factom.NewChain(&e)

	b := AccountFromFctSecret("Fs2BNvoDgSoGJpWg4PvRUxqvLE28CQexp5FZM9X5qU6QvzFBUn6D")
	commit, _ := ComposeChainCommit(b.Priv, chain)
	reveal, _ := ComposeRevealEntryMsg(b.Priv, chain.FirstEntry)

	// generate some Acks
	leader := leader.New(state.StateConfig.LeaderConfig).Leader
	commitAck := leader.NewAck(commit, nil).(*messages.Ack)
	revealAck := leader.NewAck(reveal, nil).(*messages.Ack)

	// should be true by design
	assert.Equal(t, commit.GetMsgHash(), commitAck.MessageHash)
	assert.Equal(t, reveal.GetMsgHash(), revealAck.MessageHash)
	assert.Equal(t, commit.GetHash(), reveal.GetHash())

	// load up 2 matched messages & 1 missing Ack
	ml.Add(commit)
	ml.Add(commitAck)

	ml.Add(reveal)
	ml.Add(revealAck)

	assert.Equal(t, 2, len(ml.PairList))
	assert.Equal(t, 0, len(ml.MsgList))
	assert.Equal(t, 0, len(ml.AckList))
}
