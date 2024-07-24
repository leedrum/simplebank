package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/leedrum/simplebank/api"
	db "github.com/leedrum/simplebank/db/sqlc"
	"github.com/leedrum/simplebank/gapi"
	"github.com/leedrum/simplebank/pb"
	"github.com/leedrum/simplebank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	// runGinServer(config, store)
	runRPCServer(config, store)
}

func runRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listen, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	log.Printf("start gRPC server on %s", listen.Addr().String())
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal("cannot start gRPC server: ", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
