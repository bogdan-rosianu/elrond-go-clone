package builtInFunctions

import (
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-go/state"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	vmcommonBuiltInFunctions "github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
)

// ArgsCreateBuiltInFunctionContainer defines the argument structure to create new built in function container
type ArgsCreateBuiltInFunctionContainer struct {
	GasSchedule                      core.GasScheduleNotifier
	MapDNSAddresses                  map[string]struct{}
	EnableUserNameChange             bool
	Marshalizer                      marshal.Marshalizer
	Accounts                         state.AccountsAdapter
	ShardCoordinator                 sharding.Coordinator
	EpochNotifier                    vmcommon.EpochNotifier
	ESDTMultiTransferEnableEpoch     uint32
	ESDTTransferRoleEnableEpoch      uint32
	GlobalMintBurnDisableEpoch       uint32
	ESDTTransferMetaEnableEpoch      uint32
	OptimizeNFTStoreEnableEpoch      uint32
	CheckCorrectTokenIDEnableEpoch   uint32
	CheckFunctionArgumentEnableEpoch uint32
}

// CreateBuiltInFuncContainerAndNFTStorageHandler creates a container that will hold all the available built in functions
func CreateBuiltInFuncContainerAndNFTStorageHandler(args ArgsCreateBuiltInFunctionContainer) (vmcommon.BuiltInFunctionContainer, vmcommon.SimpleESDTNFTStorageHandler, vmcommon.ESDTGlobalSettingsHandler, error) {
	if check.IfNil(args.GasSchedule) {
		return nil, nil, nil, process.ErrNilGasSchedule
	}
	if check.IfNil(args.Marshalizer) {
		return nil, nil, nil, process.ErrNilMarshalizer
	}
	if check.IfNil(args.Accounts) {
		return nil, nil, nil, process.ErrNilAccountsAdapter
	}
	if args.MapDNSAddresses == nil {
		return nil, nil, nil, process.ErrNilDnsAddresses
	}
	if check.IfNil(args.ShardCoordinator) {
		return nil, nil, nil, process.ErrNilShardCoordinator
	}
	if check.IfNil(args.EpochNotifier) {
		return nil, nil, nil, process.ErrNilEpochNotifier
	}

	vmcommonAccounts, ok := args.Accounts.(vmcommon.AccountsAdapter)
	if !ok {
		return nil, nil, nil, process.ErrWrongTypeAssertion
	}

	modifiedArgs := vmcommonBuiltInFunctions.ArgsCreateBuiltInFunctionContainer{
		GasMap:                              args.GasSchedule.LatestGasSchedule(),
		MapDNSAddresses:                     args.MapDNSAddresses,
		EnableUserNameChange:                args.EnableUserNameChange,
		Marshalizer:                         args.Marshalizer,
		Accounts:                            vmcommonAccounts,
		ShardCoordinator:                    args.ShardCoordinator,
		EpochNotifier:                       args.EpochNotifier,
		ESDTNFTImprovementV1ActivationEpoch: args.ESDTMultiTransferEnableEpoch,
		ESDTTransferToMetaEnableEpoch:       args.ESDTTransferMetaEnableEpoch,
		ESDTTransferRoleEnableEpoch:         args.ESDTTransferRoleEnableEpoch,
		GlobalMintBurnDisableEpoch:          args.GlobalMintBurnDisableEpoch,
		SaveNFTToSystemAccountEnableEpoch:   args.OptimizeNFTStoreEnableEpoch,
		CheckCorrectTokenIDEnableEpoch:      args.CheckCorrectTokenIDEnableEpoch,
		CheckFunctionArgumentEnableEpoch:    args.CheckFunctionArgumentEnableEpoch,
	}

	bContainerFactory, err := vmcommonBuiltInFunctions.NewBuiltInFunctionsCreator(modifiedArgs)
	if err != nil {
		return nil, nil, nil, err
	}

	container, err := bContainerFactory.CreateBuiltInFunctionContainer()
	if err != nil {
		return nil, nil, nil, err
	}

	args.GasSchedule.RegisterNotifyHandler(bContainerFactory)

	return container, bContainerFactory.NFTStorageHandler(), bContainerFactory.ESDTGlobalSettingsHandler(), nil
}
