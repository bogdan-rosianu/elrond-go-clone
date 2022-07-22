package genericMocks

import (
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/storage"
)

// ChainStorerMock -
type ChainStorerMock struct {
	Transactions *StorerMock
	Rewards      *StorerMock
	Unsigned     *StorerMock
	HdrNonce     *StorerMock
}

// NewChainStorerMock -
func NewChainStorerMock(epoch uint32) *ChainStorerMock {
	return &ChainStorerMock{
		Transactions: NewStorerMock("Transactions", epoch),
		Rewards:      NewStorerMock("Rewards", epoch),
		Unsigned:     NewStorerMock("Unsigned", epoch),
		HdrNonce:     NewStorerMock("HeaderNonce", epoch),
	}
}

// CloseAll -
func (sm *ChainStorerMock) CloseAll() error {
	return nil
}

// AddStorer -
func (sm *ChainStorerMock) AddStorer(_ dataRetriever.UnitType, _ storage.Storer) {
	panic("not supported")
}

// GetStorer -
func (sm *ChainStorerMock) GetStorer(unitType dataRetriever.UnitType) storage.Storer {
	if unitType == dataRetriever.TransactionUnit {
		return sm.Transactions
	}
	if unitType == dataRetriever.RewardTransactionUnit {
		return sm.Rewards
	}
	if unitType == dataRetriever.UnsignedTransactionUnit {
		return sm.Unsigned
	}

	return sm.HdrNonce
}

// Has -
func (sm *ChainStorerMock) Has(_ dataRetriever.UnitType, _ []byte) error {
	return nil
}

// Get -
func (sm *ChainStorerMock) Get(unitType dataRetriever.UnitType, key []byte) ([]byte, error) {
	return sm.GetStorer(unitType).Get(key)
}

// Put -
func (sm *ChainStorerMock) Put(unitType dataRetriever.UnitType, key []byte, value []byte) error {
	return sm.GetStorer(unitType).Put(key, value)
}

// GetAll -
func (sm *ChainStorerMock) GetAll(unitType dataRetriever.UnitType, keys [][]byte) (map[string][]byte, error) {
	storer := sm.GetStorer(unitType)

	data := make(map[string][]byte)
	for _, key := range keys {
		buff, err := storer.Get(key)
		if err != nil {
			return nil, err
		}

		data[string(key)] = buff
	}

	return data, nil
}

// SetEpochForPutOperation -
func (sm *ChainStorerMock) SetEpochForPutOperation(_ uint32) {
	panic("not supported")
}

// GetAllStorers -
func (sm *ChainStorerMock) GetAllStorers() map[dataRetriever.UnitType]storage.Storer {
	return map[dataRetriever.UnitType]storage.Storer{
		dataRetriever.TransactionUnit:          sm.Transactions,
		dataRetriever.RewardTransactionUnit:    sm.Rewards,
		dataRetriever.UnsignedTransactionUnit:  sm.Unsigned,
		dataRetriever.MetaHdrNonceHashDataUnit: sm.HdrNonce,
	}
}

// Destroy -
func (sm *ChainStorerMock) Destroy() error {
	return nil
}

// IsInterfaceNil -
func (sm *ChainStorerMock) IsInterfaceNil() bool {
	return sm == nil
}
