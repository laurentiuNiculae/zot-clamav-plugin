package main

import (
	"fmt"
	"log"
	"net"

	. "github.com/laurentiuNiculae/zot-clamav-plugin/server"
	"google.golang.org/grpc"
	"zotregistry.io/zot/pkg/plugins/scan"
)

const (
	port = 9100
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	scan.RegisterScanServer(s, ScanServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
