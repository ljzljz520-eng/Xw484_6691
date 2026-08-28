package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"genealogy-story-organizer/internal/application"
	"genealogy-story-organizer/internal/store"
	"genealogy-story-organizer/internal/web"
)

func main() {
	dbPath := flag.String("db", "genealogy.db", "path to bbolt database")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	database, err := store.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	service, err := application.NewService(database, application.NewSequenceID(), store.StaticClock{})
	if err != nil {
		log.Fatal(err)
	}
	server := web.NewServer(service)
	fmt.Printf("家谱故事整理系统 listening on %s\n", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
