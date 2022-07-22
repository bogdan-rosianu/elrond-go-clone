package mock

import (
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/testscommon"
)

// VmMachinesContainerFactoryMock -
type VmMachinesContainerFactoryMock struct {
}

// Create -
func (v *VmMachinesContainerFactoryMock) Create() (process.VirtualMachinesContainer, error) {
	return &VMContainerMock{}, nil
}

// Close -
func (v *VmMachinesContainerFactoryMock) Close() error {
	return nil
}

// BlockChainHookImpl -
func (v *VmMachinesContainerFactoryMock) BlockChainHookImpl() process.BlockChainHookHandler {
	return &testscommon.BlockChainHookStub{}
}

// IsInterfaceNil -
func (v *VmMachinesContainerFactoryMock) IsInterfaceNil() bool {
	return v == nil
}
