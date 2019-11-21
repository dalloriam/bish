package command

type ShellContext interface {
	GetKey(domain string, key string) (interface{}, bool)
	SetKey(domain string, key string, value interface{})
}
