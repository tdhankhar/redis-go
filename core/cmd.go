package core

type bulkString string

type RedisCmd struct {
	Cmd string
	Args []string
}
