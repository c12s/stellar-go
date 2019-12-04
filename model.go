package main

type Spanner interface {
	Child(name string) *Span
	AddLog(kv *KV)
	AddTag(kv *KV)
	AddBaggage(kv *KV)
	StartTime()
	EndTime()
	Finish() // send data to collecto and maybe serialize to ctx ot request
	Serialize() *Values
	Marshall() ([]byte, error)
}

type KV struct {
	key   string
	value string
}

type Values struct {
	md map[string][]string
}

func (v Values) Get(key string) []string {
	return v.md[key]
}
