package engine

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/FactomProject/factomd/common/constants"

	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/factom"

	"crypto/sha256"

	"github.com/FactomProject/factomd/common/entryBlock"
	"github.com/FactomProject/factomd/common/entryCreditBlock"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/common/primitives/random"
	"github.com/FactomProject/factomd/state"
	"github.com/FactomProject/factomd/util"
	"github.com/FactomProject/factomd/util/atomic"
)

type LoadGenerator struct {
	ECKey       *primitives.PrivateKey // Entry Credit private key
	ToSend      int                    // How much to send
	PerSecond   atomic.AtomicInt       // How much per second
	stop        chan bool              // Stop the go routine
	running     atomic.AtomicBool      // We are running
	tight       atomic.AtomicBool      // Only allocate ECs as needed (more EC purchases)
	chainExists atomic.AtomicBool      // We have the FER chain already
	fee         atomic.AtomicBool      // If fee == true, then change the fees as we go
	rate        atomic.AtomicInt64     // Current Rate
	txoffset    int64                  // Offset to be added to the timestamp of created tx to test time limits.
	state       *state.State           // Access to logging
}

// NewLoadGenerator makes a new load generator. The state is used for funding the transaction
func NewLoadGenerator(s *state.State) *LoadGenerator {
	lg := new(LoadGenerator)
	lg.ECKey, _ = primitives.NewPrivateKeyFromHex(ecSec)
	lg.stop = make(chan bool, 5)
	lg.state = s

	return lg
}

//func (lg *LoadGenerator) Fees() {
//
//	for lg.fee.Load() {
//		time.Sleep(1 * time.Second)
//		if lg.rate.Load() > 10000 {
//			lg.rate.Store(900)
//		}
//		lg.rate.Add(100)
//
//		nextRate := lg.rate.Load()
//		for s
//	}
//}

func (lg *LoadGenerator) Run() {
	if lg.running.Load() {
		return
	}

	lg.running.Store(true)
	//FundWallet(fnodes[wsapiNode].State, 15000e8)

	// Every second add the per second amount
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		select {
		case <-lg.stop:
			lg.running.Store(false)
			return
		default:

		}
		addSend := lg.PerSecond.Load()
		lg.ToSend += addSend
		top := lg.ToSend / 10      // ToSend is in tenths so get the integer part
		lg.ToSend = lg.ToSend % 10 // save an fractional part for next iteration
		if addSend == 0 {
			lg.running.Store(false)
			return
		}
		var chain interfaces.IHash = nil

		for i := 0; i < top; i++ {
			var c interfaces.IMsg
			e := RandomEntry(nil)
			if chain == nil {
				c = lg.NewCommitChain(e)
				chain = e.ChainID
			} else {
				e.ChainID = chain
				c = lg.NewCommitEntry(e)
			}
			r := lg.NewRevealEntry(e)
			s := fnodes[wsapiNode].State
			s.APIQueue().Enqueue(c)
			s.APIQueue().Enqueue(r)
			time.Sleep(time.Duration(800/top) * time.Millisecond) // spread the load out over 800ms + overhead
		}
	}
}

func (lg *LoadGenerator) Stop() {
	lg.stop <- true
}

func RandomEntry(entry *entryBlock.Entry) *entryBlock.Entry {
	if entry == nil {
		entry.Content = primitives.ByteSlice{random.RandByteSliceOfLen(rand.Intn(128) + 128)}
		entry.ExtIDs = make([]primitives.ByteSlice, rand.Intn(4)+1)
		raw := make([][]byte, len(entry.ExtIDs))
		for i := range entry.ExtIDs {
			entry.ExtIDs[i] = primitives.ByteSlice{random.RandByteSliceOfLen(rand.Intn(32) + 32)}
			raw[i] = entry.ExtIDs[i].Bytes
		}
	}
	sum := sha256.New()
	for i, v := range entry.ExtIDs {
		fmt.Printf("%d %x\n", i, v)
		x := sha256.Sum256(v.Bytes)
		sum.Write(x[:])
	}
	originalHash := sum.Sum(nil)
	checkHash := primitives.Shad(originalHash)

	entry.ChainID = checkHash
	return entry
}

func (lg *LoadGenerator) NewRevealEntry(entry *entryBlock.Entry) *messages.RevealEntryMsg {
	msg := messages.NewRevealEntryMsg()
	msg.Entry = entry
	msg.Timestamp = primitives.NewTimestampNow()

	return msg
}

var cnt int
var goingUp bool
var limitBuys = true // We limit buys only after the first only attempted purchase, so people can fund identities in testing

func (lg *LoadGenerator) KeepUsFunded() {

	s := fnodes[wsapiNode].State

	var level int64

	buys := 0
	totalBought := 0
	for i := 0; ; i++ {

		ts := "false"
		if lg.tight.Load() {
			ts = "true"
		}

		//EC3Eh7yQKShgjkUSFrPbnQpboykCzf4kw9QHxi47GGz5P2k3dbab is EC address

		if !limitBuys {
			level = 200 // Only do this once, after that look for requests for load to drive EC buys.
		} else if lg.tight.Load() {
			level = 10
		} else {
			level = 10000
		}

		outEC, _ := primitives.HexToHash("c23ae8eec2beb181a0da926bd2344e988149fbe839fbc7489f2096e7d6110243")
		outAdd := factoid.NewAddress(outEC.Bytes())
		ecBal := s.GetE(true, outAdd.Fixed())
		ecPrice := s.GetFactoshisPerEC()

		if ecBal < level {
			buys++
			need := level - ecBal + level*2
			totalBought += int(need)
			FundWalletTOFF(s, lg.txoffset, uint64(need)*ecPrice)
		} else {
			limitBuys = true
		}

		if i%5 == 0 {
			lg.state.LogPrintf("loadgenerator", "Tight %7s Total TX %6d for a total of %8d entry credits balance %d.",
				ts, buys, totalBought, ecBal)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (lg *LoadGenerator) NewCommitChain(entry *entryBlock.Entry) *messages.CommitChainMsg {

	msg := new(messages.CommitChainMsg)

	// form commit
	commit := entryCreditBlock.NewCommitChain()
	data, _ := entry.MarshalBinary()
	commit.Credits, _ = util.EntryCost(data)
	commit.Credits += 10

	commit.EntryHash = entry.GetHash()
	var b6 primitives.ByteSlice6
	copy(b6[:], milliTime(lg.txoffset)[:])
	commit.MilliTime = &b6
	var b32 primitives.ByteSlice32
	copy(b32[:], lg.ECKey.Pub[:])
	commit.ECPubKey = &b32
	commit.Weld = entry.GetWeldHash()
	commit.ChainIDHash = entry.ChainID

	commit.Sign(lg.ECKey.Key[:])

	// form msg
	msg.CommitChain = commit
	msg.Sign(lg.ECKey)

	return msg
}

func (lg *LoadGenerator) NewCommitEntry(entry *entryBlock.Entry) *messages.CommitEntryMsg {
	msg := messages.NewCommitEntryMsg()

	// form commit
	commit := entryCreditBlock.NewCommitEntry()
	data, _ := entry.MarshalBinary()

	commit.Credits, _ = util.EntryCost(data)

	commit.EntryHash = entry.GetHash()
	var b6 primitives.ByteSlice6
	copy(b6[:], milliTime(lg.txoffset)[:])
	commit.MilliTime = &b6
	var b32 primitives.ByteSlice32
	copy(b32[:], lg.ECKey.Pub[:])
	commit.ECPubKey = &b32
	commit.Sign(lg.ECKey.Key[:])

	// form msg
	msg.CommitEntry = commit
	msg.Sign(lg.ECKey)

	return msg
}

// milliTime returns a 6 byte slice representing the unix time in milliseconds
func milliTime(offset int64) (r []byte) {
	buf := new(bytes.Buffer)
	t := time.Now().UnixNano()
	m := t/1e6 + offset
	binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes()[2:]
}

// FEREntry
// This is a representation of the FER data.  Basically the json of this will be the factom entry content
type FEREntry struct {
	ExpirationHeight       uint32 `json:"expiration_height"`        // Height Entry is invalid if not activated. (6 blks from activation at most)
	TargetActivationHeight uint32 `json:"target_activation_height"` // Height to activate the entry
	Priority               uint32 `json:"priority"`                 // Some priority > 0.  Lets multiple FEREntries per block, with the highest priority winning
	TargetPrice            uint64 `json:"target_price"`             // Factoshis per EC so 1/10 cent * TargetPrice / 1000 = FCT price
	Version                string `json:"version"`                  // Should always be 1.0 until we update the FER implementation
}

// CreateFERChain
// When running a simulation, we have to create the FERChain to change the exchange rate.
// This function does this for us.
func (lg *LoadGenerator) CreateFERChain() {

	if lg.chainExists.Load() {
		return
	}

	chainID, _ := primitives.HexToHash(constants.FERChainID)
	eb, err := lg.state.DB.FetchEBlockHead(chainID)
	// Only create the FER chain if it doesn't exist.
	if err != nil || eb == nil {
		// create the FER chain by creating an Entry
		fer := entryBlock.NewEntry()
		// Add the Extended IDs (ExtIDs)
		extIDs1 := primitives.ByteSlice{Bytes: []byte("FCT EC Conversion Rate Chain")}
		extIDs2 := primitives.ByteSlice{Bytes: []byte("1950454129")}
		content := primitives.ByteSlice{Bytes: []byte("This chain contains messages which coordinate " +
			"the FCT to EC conversion rate amongst factomd nodes.")}

		fer.ExtIDs = append(fer.ExtIDs, extIDs1)
		fer.ExtIDs = append(fer.ExtIDs, extIDs2)
		// Content can be anything, but we will use the same content from the Main Net
		fer.Content = content
		// Compute the ChainID
		fer = RandomEntry(fer)

		c := lg.NewCommitChain(fer)
		e := lg.NewRevealEntry(fer)
		lg.state.APIQueue().Enqueue(c)
		lg.state.APIQueue().Enqueue(e)
	}
}

// SetExchangeRate -- Write an entry to change the EC Exchange rate into the FER chain.
// This is only going to work if the unit test or the configuration file has set the ExchangeRateAuthorityPublicKey
// to the public key for the null (all zeros) private key.
func (lg *LoadGenerator) SetExchangeRate(DBHeight uint32, priority uint32, factoshis uint64) {

	fer := new(FEREntry)
	fer.ExpirationHeight = DBHeight + 2
	fer.Priority = priority
	fer.TargetPrice = factoshis
	fer.Version = "1.0"

	entryJson, err := json.Marshal(fer)

	if err != nil {
		os.Stderr.WriteString("Could not marshal the data into an FEREntry\n")
		return
	}

	e := new(factom.Entry)
	e.ExtIDs = append(e.ExtIDs, []byte("FCT EC Conversion Rate Chain"))
	e.ChainID = constants.FERChainID

	var signingPrivateKey [64]byte
	ed.GetPublicKey(&signingPrivateKey)
	signingSignature := ed.Sign(&signingPrivateKey, entryJson)

	e.ExtIDs = append(e.ExtIDs, signingSignature[:])
	e.Content = entryJson

}
