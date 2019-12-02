package main

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
