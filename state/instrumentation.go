package state

import "github.com/FactomProject/factomd/telemetry"

var (
	RegisterPrometheus = telemetry.RegisterPrometheus

	// TODO: refactor to create metrics during initialization

	// Entry Syncing Controller
	ESMissingQueue   = telemetry.ESMissingQueue
	ESMissing        = telemetry.ESMissing
	ESFound          = telemetry.ESFound
	ESAsking         = telemetry.ESAsking
	ESHighestAsking  = telemetry.ESHighestAsking
	ESHighestMissing = telemetry.ESHighestMissing
	ESFirstMissing   = telemetry.ESFirstMissing
	ESDBHTComplete   = telemetry.ESDBHTComplete
	ESAvgRequests    = telemetry.ESAvgRequests
	HighestAck       = telemetry.HighestAck
	HighestKnown     = telemetry.HighestKnown
	HighestSaved     = telemetry.HighestSaved
	HighestCompleted = telemetry.HighestCompleted

	// TPS
	TotalTransactionPerSecond   = telemetry.TotalTransactionPerSecond
	InstantTransactionPerSecond = telemetry.InstantTransactionPerSecond

	// Torrent
	stateTorrentSyncingLower = telemetry.StateTorrentSyncingLower
	stateTorrentSyncingUpper = telemetry.StateTorrentSyncingUpper

	// Queues
	CurrentMessageQueueInMsgGeneralVec   = telemetry.CurrentMessageQueueInMsgGeneralVec
	TotalMessageQueueInMsgGeneralVec     = telemetry.TotalMessageQueueInMsgGeneralVec
	CurrentMessageQueueApiGeneralVec     = telemetry.CurrentMessageQueueApiGeneralVec
	TotalMessageQueueApiGeneralVec       = telemetry.TotalMessageQueueApiGeneralVec
	TotalMessageQueueNetOutMsgGeneralVec = telemetry.TotalMessageQueueNetOutMsgGeneralVec

	// MsgQueue chan
	TotalMsgQueueInputs  = telemetry.TotalMsgQueueInputs
	TotalMsgQueueOutputs = telemetry.TotalMsgQueueOutputs

	// Holding Queue
	TotalHoldingQueueInputs        = telemetry.TotalHoldingQueueInputs
	TotalHoldingQueueOutputs       = telemetry.TotalHoldingQueueOutputs
	TotalHoldingQueueRecycles      = telemetry.TotalHoldingQueueRecycles
	HoldingQueueDBSigInputs        = telemetry.HoldingQueueDBSigInputs
	HoldingQueueDBSigOutputs       = telemetry.HoldingQueueDBSigOutputs
	HoldingQueueCommitEntryInputs  = telemetry.HoldingQueueCommitEntryInputs
	HoldingQueueCommitEntryOutputs = telemetry.HoldingQueueCommitEntryOutputs
	HoldingQueueCommitChainInputs  = telemetry.HoldingQueueCommitChainInputs
	HoldingQueueCommitChainOutputs = telemetry.HoldingQueueCommitChainOutputs
	HoldingQueueRevealEntryInputs  = telemetry.HoldingQueueRevealEntryInputs
	HoldingQueueRevealEntryOutputs = telemetry.HoldingQueueRevealEntryOutputs

	// Acks Queue                          // Acks Queue
	TotalAcksInputs  = telemetry.TotalAcksInputs
	TotalAcksOutputs = telemetry.TotalAcksOutputs

	// Commits map                         // Commits map
	TotalCommitsInputs  = telemetry.TotalCommitsInputs
	TotalCommitsOutputs = telemetry.TotalCommitsOutputs

	// XReview Queue                       // XReview Queue
	TotalXReviewQueueInputs  = telemetry.TotalXReviewQueueInputs
	TotalXReviewQueueOutputs = telemetry.TotalXReviewQueueOutputs

	// Executions                          // Executions
	LeaderExecutions             = telemetry.LeaderExecutions
	FollowerExecutions           = telemetry.FollowerExecutions
	LeaderEOMExecutions          = telemetry.LeaderEOMExecutions
	FollowerEOMExecutions        = telemetry.FollowerEOMExecutions
	FollowerMissingMsgExecutions = telemetry.FollowerMissingMsgExecutions

	// ProcessList                         // ProcessList
	TotalProcessListInputs    = telemetry.TotalProcessListInputs
	TotalProcessListProcesses = telemetry.TotalProcessListProcesses
	TotalProcessEOMs          = telemetry.TotalProcessEOMs

	// Durations                           // Durations
	TotalReviewHoldingTime   = telemetry.TotalReviewHoldingTime
	TotalProcessXReviewTime  = telemetry.TotalProcessXReviewTime
	TotalProcessProcChanTime = telemetry.TotalProcessProcChanTime
	TotalEmptyLoopTime       = telemetry.TotalEmptyLoopTime
	TotalAckLoopTime         = telemetry.TotalAckLoopTime
	TotalExecuteMsgTime      = telemetry.TotalExecuteMsgTime
)

var registered bool = false

// RegisterPrometheus registers the variables to be exposed. This can only be run once, hence the
// boolean flag to prevent panics if launched more than once. This is called in NetStart
func RegisterPrometheus() {
	if registered {
		return
	}
	registered = true
	// 		Example Cont.
	// prometheus.MustRegister(stateRandomCounter)

	// Entry syncing
	prometheus.MustRegister(ESAsking)
	prometheus.MustRegister(ESHighestAsking)
	prometheus.MustRegister(ESFirstMissing)
	prometheus.MustRegister(ESMissing)
	prometheus.MustRegister(ESFound)
	prometheus.MustRegister(ESDBHTComplete)
	prometheus.MustRegister(ESMissingQueue)
	prometheus.MustRegister(ESHighestMissing)
	prometheus.MustRegister(ESAvgRequests)
	prometheus.MustRegister(HighestAck)
	prometheus.MustRegister(HighestKnown)
	prometheus.MustRegister(HighestSaved)
	prometheus.MustRegister(HighestCompleted)

	// TPS
	prometheus.MustRegister(TotalTransactionPerSecond)
	prometheus.MustRegister(InstantTransactionPerSecond)

	// Torrent
	prometheus.MustRegister(stateTorrentSyncingLower)
	prometheus.MustRegister(stateTorrentSyncingUpper)

	// Queues
	prometheus.MustRegister(CurrentMessageQueueInMsgGeneralVec)
	prometheus.MustRegister(TotalMessageQueueInMsgGeneralVec)
	prometheus.MustRegister(CurrentMessageQueueApiGeneralVec)
	prometheus.MustRegister(TotalMessageQueueApiGeneralVec)
	prometheus.MustRegister(TotalMessageQueueNetOutMsgGeneralVec)

	// MsgQueue chan
	prometheus.MustRegister(TotalMsgQueueInputs)
	prometheus.MustRegister(TotalMsgQueueOutputs)

	// Holding
	prometheus.MustRegister(TotalHoldingQueueInputs)
	prometheus.MustRegister(TotalHoldingQueueOutputs)
	prometheus.MustRegister(HoldingQueueDBSigInputs)
	prometheus.MustRegister(HoldingQueueDBSigOutputs)
	prometheus.MustRegister(HoldingQueueCommitEntryInputs)
	prometheus.MustRegister(HoldingQueueCommitEntryOutputs)
	prometheus.MustRegister(HoldingQueueCommitChainInputs)
	prometheus.MustRegister(HoldingQueueCommitChainOutputs)
	prometheus.MustRegister(HoldingQueueRevealEntryInputs)
	prometheus.MustRegister(HoldingQueueRevealEntryOutputs)

	// Acks
	prometheus.MustRegister(TotalAcksInputs)
	prometheus.MustRegister(TotalAcksOutputs)

	// Execution
	prometheus.MustRegister(LeaderExecutions)
	prometheus.MustRegister(FollowerExecutions)
	prometheus.MustRegister(LeaderEOMExecutions)
	prometheus.MustRegister(FollowerEOMExecutions)
	prometheus.MustRegister(FollowerMissingMsgExecutions)

	// ProcessList
	prometheus.MustRegister(TotalProcessListInputs)
	prometheus.MustRegister(TotalProcessListProcesses)
	prometheus.MustRegister(TotalProcessEOMs)

	// XReview Queue
	prometheus.MustRegister(TotalXReviewQueueInputs)
	prometheus.MustRegister(TotalXReviewQueueOutputs)

	// Commits map
	prometheus.MustRegister(TotalCommitsInputs)
	prometheus.MustRegister(TotalCommitsOutputs)

	// Durations
	prometheus.MustRegister(TotalReviewHoldingTime)
	prometheus.MustRegister(TotalProcessXReviewTime)
	prometheus.MustRegister(TotalProcessProcChanTime)
	prometheus.MustRegister(TotalEmptyLoopTime)
	prometheus.MustRegister(TotalAckLoopTime)
	prometheus.MustRegister(TotalExecuteMsgTime)
}
