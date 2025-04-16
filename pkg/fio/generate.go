package fio

//go:generate go run go.uber.org/mock/mockgen -package mocks -destination mocks/mock_fio.go github.com/kdwils/splinter/pkg/fio FileIO,WriteCloser,DirEntry
