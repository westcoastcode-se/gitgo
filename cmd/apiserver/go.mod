module gitgo/apiserver

go 1.17

replace gitgo/api v1.0.0 => ../../api

require (
	gitgo/api v1.0.0
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
)

require golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
