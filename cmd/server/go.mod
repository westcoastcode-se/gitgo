module gitgo/server

go 1.17

replace gitgo/api v1.0.0 => ../../api

require (
	gitgo/api v1.0.0
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b
)

require golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
