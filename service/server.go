package service

import (
	"net/http"
	"strings"

	"github.com/GXK666/eosTransfer/service/general"
	"github.com/GXK666/eosTransfer/transfer"

	"crypto/tls"
	"net"

	"context"
	"time"

	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kazegusuri/grpc-panic-handler"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Serve() {
	setupLog()
	panichandler.InstallPanicHandler(LogPanicHandler)

	// Protect GRPC from aborting by panic
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(panichandler.UnaryPanicHandler),
		grpc.StreamInterceptor(panichandler.StreamPanicHandler),
	}

	grpcServer := grpc.NewServer(opts...)
	general.RegisterServiceServer(grpcServer, transfer.Server)

	mux := http.NewServeMux()

	addr := viper.GetString("rpc.addr")
	endpoint := addr
	parts := strings.Split(addr, ":")
	if parts[0] == "0.0.0.0" {
		endpoint = "127.0.0.1:" + parts[1]
	}
	certFile := viper.GetString("rpc.tls.certFile")
	keyFile := viper.GetString("rpc.tls.keyFile")

	dialCreds, err := credentials.NewClientTLSFromFile(certFile, endpoint)
	if err != nil {
		panic(err)
	}
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(dialCreds)}
	jsonpb := &gateway.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		// This is necessary to get error details properly
		// marshalled in unary requests.
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)
	go func() {
		time.Sleep(time.Second) // Avoid immediate connection failure
		err = general.RegisterServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dialOpts)
		if err != nil {
			panic(err)
		}
	}()
	mux.Handle("/", gwmux)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	})

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"h2"},
		},
	}

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
	if err != nil {
		panic(err)
	}
	return
}
