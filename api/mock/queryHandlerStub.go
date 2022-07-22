package mock

// QueryHandlerStub -
type QueryHandlerStub struct {
	QueryCalled func(search string) []string
}

// Query -
func (qhs *QueryHandlerStub) Query(search string) []string {
	if qhs.QueryCalled != nil {
		return qhs.QueryCalled(search)
	}

	return make([]string, 0)
}

// Close -
func (qhs *QueryHandlerStub) Close() error {
	return nil
}

// IsInterfaceNil -
func (qhs *QueryHandlerStub) IsInterfaceNil() bool {
	return qhs == nil
}
