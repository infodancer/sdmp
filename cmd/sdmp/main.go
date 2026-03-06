package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/infodancer/sdmp/internal/sdmp"
	pb "github.com/infodancer/sdmp/proto/sdmp/v1"
	"google.golang.org/grpc"
)

func main() {
	grpcAddr := envOrDefault("SDMP_GRPC_ADDR", ":9400")
	httpAddr := envOrDefault("SDMP_HTTP_ADDR", ":9401")
	udpAddr := envOrDefault("SDMP_UDP_ADDR", ":9402")

	// gRPC server for authenticated domain-to-domain operations.
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterFetchServiceServer(grpcServer, &sdmp.FetchServer{})
	pb.RegisterKeyServiceServer(grpcServer, &sdmp.KeyServer{})
	pb.RegisterTransferServiceServer(grpcServer, &sdmp.TransferServer{})

	go func() {
		log.Printf("sdmp: gRPC server listening on %s", grpcAddr)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	// HTTP server for CDN-cacheable blob fetch.
	mux := http.NewServeMux()
	mux.Handle("GET /message/{id}", &sdmp.BlobHandler{})
	httpServer := &http.Server{Addr: httpAddr, Handler: mux}

	go func() {
		log.Printf("sdmp: HTTP blob server listening on %s", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve: %v", err)
		}
	}()

	// UDP listener for fire-and-forget notifications.
	notifListener, err := sdmp.NewNotificationListener(udpAddr)
	if err != nil {
		log.Fatalf("failed to start notification listener on %s: %v", udpAddr, err)
	}

	go func() {
		log.Printf("sdmp: UDP notification listener on %s", udpAddr)
		if err := notifListener.Listen(func(n *pb.EncryptedNotification) {
			log.Printf("sdmp: notification received for domain %s", n.GetDestinationDomain())
		}); err != nil {
			log.Printf("notification listener stopped: %v", err)
		}
	}()

	fmt.Println("sdmp: Secure Domain Messaging Protocol server running")
	fmt.Printf("  gRPC: %s  HTTP: %s  UDP: %s\n", grpcAddr, httpAddr, udpAddr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("sdmp: shutting down")
	grpcServer.GracefulStop()
	notifListener.Close()
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
