package coordinator

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/atomic"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/process"
)

const printReportHeader = "double transactions found (this is not critical, thus)\nshowing the whole block body:\n"
const nilBlockBodyMessage = "nil block body in printDoubleTransactionsDetector.ProcessBlockBody"
const noDoubledTransactionsFoundMessage = "no double transactions found"
const doubledTransactionsFoundButFlagActive = "double transactions found but this is expected until the AddFailedRelayedTxToInvalidMBsDisableEpoch is deactivated"

// ArgsPrintDoubleTransactionsDetector is the argument DTO structure used in the NewPrintDoubleTransactionsDetector function
type ArgsPrintDoubleTransactionsDetector struct {
	Marshaller    marshal.Marshalizer
	Hasher        hashing.Hasher
	EpochNotifier process.EpochNotifier

	AddFailedRelayedTxToInvalidMBsDisableEpoch uint32
}

type printDoubleTransactionsDetector struct {
	marshaller marshal.Marshalizer
	hasher     hashing.Hasher
	logger     logger.Logger

	addFailedRelayedTxToInvalidMBsDisableEpoch uint32
	flagAddFailedRelayedTxToInvalidMBs         atomic.Flag
}

// NewPrintDoubleTransactionsDetector creates a new instance of printDoubleTransactionsDetector
func NewPrintDoubleTransactionsDetector(args ArgsPrintDoubleTransactionsDetector) (*printDoubleTransactionsDetector, error) {
	err := checkArgsPrintDoubleTransactionsDetector(args)
	if err != nil {
		return nil, err
	}

	detector := &printDoubleTransactionsDetector{
		marshaller: args.Marshaller,
		hasher:     args.Hasher,
		addFailedRelayedTxToInvalidMBsDisableEpoch: args.AddFailedRelayedTxToInvalidMBsDisableEpoch,
		logger: log,
	}

	args.EpochNotifier.RegisterNotifyHandler(detector)

	return detector, nil
}

func checkArgsPrintDoubleTransactionsDetector(args ArgsPrintDoubleTransactionsDetector) error {
	if check.IfNil(args.Marshaller) {
		return process.ErrNilMarshalizer
	}
	if check.IfNil(args.Hasher) {
		return process.ErrNilHasher
	}
	if check.IfNil(args.EpochNotifier) {
		return process.ErrNilEpochNotifier
	}

	return nil
}

// ProcessBlockBody processes the block body provided in search of doubled transactions. If there are doubled transactions,
// this method will log as error the event providing as much information as possible
func (detector *printDoubleTransactionsDetector) ProcessBlockBody(body *block.Body) {
	if body == nil {
		detector.logger.Error(nilBlockBodyMessage)
		return
	}

	transactions := make(map[string]int)
	doubleTransactionsExist := false
	printReport := strings.Builder{}

	for _, miniBlock := range body.MiniBlocks {
		mbHash, _ := core.CalculateHash(detector.marshaller, detector.hasher, miniBlock)
		log.Debug("checking for double transactions: miniblock",
			"sender shard", miniBlock.SenderShardID,
			"receiver shard", miniBlock.ReceiverShardID,
			"type", miniBlock.Type,
			"num txs", len(miniBlock.TxHashes),
			"hash", mbHash)
		printReport.WriteString(fmt.Sprintf(" miniblock hash %s, type %s, %d -> %d\n",
			hex.EncodeToString(mbHash), miniBlock.Type.String(), miniBlock.SenderShardID, miniBlock.ReceiverShardID))

		for _, txHash := range miniBlock.TxHashes {
			transactions[string(txHash)]++
			printReport.WriteString(fmt.Sprintf("  tx hash %s\n", hex.EncodeToString(txHash)))

			doubleTransactionsExist = doubleTransactionsExist || transactions[string(txHash)] > 1
		}
	}

	if !doubleTransactionsExist {
		detector.logger.Debug(noDoubledTransactionsFoundMessage)
		return
	}
	if detector.flagAddFailedRelayedTxToInvalidMBs.IsSet() {
		detector.logger.Debug(doubledTransactionsFoundButFlagActive)
		return
	}

	detector.logger.Error(printReportHeader + printReport.String())
}

// EpochConfirmed is called whenever a new epoch is confirmed
func (detector *printDoubleTransactionsDetector) EpochConfirmed(epoch uint32, _ uint64) {
	detector.flagAddFailedRelayedTxToInvalidMBs.SetValue(epoch < detector.addFailedRelayedTxToInvalidMBsDisableEpoch)
	log.Debug("printDoubleTransactionsDetector: add failed relayed tx to invalid miniblocks",
		"enabled", detector.flagAddFailedRelayedTxToInvalidMBs.IsSet())
}

// IsInterfaceNil returns true if there is no value under the interface
func (detector *printDoubleTransactionsDetector) IsInterfaceNil() bool {
	return detector == nil
}
