package kv

type Kv interface {
	Add(col, key string, val interface{}) error
	Get(col, key string) (interface{}, error)
	GetAll(col string) (map[string]interface{}, error)
	Delete(col, key string) error
}
