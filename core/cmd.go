package core

type BulkString string

type RedisCmd struct {
	Cmd string
	Args []string
}