package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"tztask/api"
	"tztask/conf"
	"tztask/domain/dispatch"
	"tztask/domain/service"
	"tztask/repo_impl"
	"tztask/utils/di"
	"tztask/utils/sequence"
)

func init() {
	conf.Init()
	sequence.Init()
	repo_impl.DIRegister()
	service.DIRegister()
	di.MustBindALL()
}

func main() {
	defer Recover()

	if err := dispatch.Start(); err != nil {
		panic(err)
	}
	signalExit(func() {
		dispatch.Close()
	})
	log.Println("tztask start ...")
	if err := api.HTTPServe(); err != nil {
		panic(err)
	}

}

func Recover() {
	if err := recover(); err != nil {
		log.Printf("tztask panic: %v", err)
	}
}

func signalExit(f func()) {
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		s := <-c
		log.Printf("service get a signal: %v", s)
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGHUP:
			f()
			log.Println("service closed")
			os.Exit(0)
			return
		default:
			return
		}
	}
}
