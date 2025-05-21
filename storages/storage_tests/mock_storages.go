package storages_tests

func (m *MockRedisClient) Set(key string, value interface{}) error {
	return nil
}

func (m *MockRedisClient) Get(key string, typeInfo interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockRedisClient) Del(key string) error {
	return nil
}

func (m *MockRedisClient) FlushAll() error {
	return nil
}
