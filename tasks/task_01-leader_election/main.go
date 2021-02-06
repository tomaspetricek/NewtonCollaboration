package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"time"
)

func tick(interval time.Duration, tick chan bool){
	for{
		time.Sleep(interval)
		tick <- true
	}
}

func wait(duration time.Duration, done chan bool){
	time.Sleep(duration)
	done <- true
}

func electLeader(e *concurrency.Election, ctx context.Context, done chan bool){
	if err := e.Campaign(ctx, "e"); err != nil {
		log.Fatal(err)
	}
	done <- true
}

func doTillWorkDone(){
	for {
		select{
		case <-printTick:
			fmt.Println("I'm a leader.")
		case <-workDone:
			fmt.Println("Work is done.")
			return
		}
	}
}

func doTillElected(){
	fmt.Println("Waiting to be elected.")
	for {
		select {
		case <- isLeader:
			return
		case <-printTick:
			fmt.Println("I'm a follower.")
		}
	}
}

const (
	workTime = 15 * time.Second
	printInterval = 3 * time.Second
)

var (
	printTick = make(chan bool, 1)
	workDone = make(chan bool, 1)
	isLeader = make(chan bool, 1)
)

func main() {
	// create a etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// create a sessions to elect a Leader
	s, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	e := concurrency.NewElection(s, "/leader-election")
	ctx := context.Background()

	go tick(printInterval, printTick)
	for{
		// get elected
		go electLeader(e, ctx, isLeader)
		doTillElected()

		// do work
		go wait(workTime, workDone)
		doTillWorkDone()

		// resign
		if err := e.Resign(ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Println("I've resigned.")
	}
}
