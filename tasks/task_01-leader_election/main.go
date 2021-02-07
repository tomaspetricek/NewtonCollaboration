package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"time"
)

func electLeader(e *concurrency.Election, ctx context.Context, done chan bool){
	if err := e.Campaign(ctx, "e"); err != nil {
		log.Fatal(err)
	}
	done <- true
}

func doTillWorkDone(){
	for {
		select{
		case <-printTicker.C:
			fmt.Println("I'm a leader.")
		case <-workTimer.C:
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
		case <-printTicker.C:
			fmt.Println("I'm a follower.")
		}
	}
}

const (
	workTime = 15 * time.Second
	printInterval = 3 * time.Second
)

var (
	printTicker   *time.Ticker
	workTimer     *time.Timer
	isLeader    = make(chan bool, 1)
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

	printTicker = time.NewTicker(printInterval)
	for{
		// get elected
		go electLeader(e, ctx, isLeader)
		doTillElected()

		// do work
		workTimer = time.NewTimer(workTime)
		doTillWorkDone()

		// resign
		if err := e.Resign(ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Println("I've resigned.")
	}
}
