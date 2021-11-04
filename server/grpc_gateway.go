package server

import (
	"app/pkg/ui/swagger"
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/antsgo/antsgo/conf"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewGrpcGateway(c conf.Conf, logger *logrus.Logger, registerHttpHandler func(*runtime.ServeMux, *grpc.ClientConn)) {
	grpcAddr := fmt.Sprintf(":%d", c.ServerGrpc.Port)
	grpcGwAddr := fmt.Sprintf(":%d", c.ServerGrpcGw.Port)
	conn, err := grpc.DialContext(context.Background(), grpcAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("failed to dial server: %v", err)
	}

	gwmux := runtime.NewServeMux()

	registerHttpHandler(gwmux, conn)

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.HandleFunc("/swagger/", serverSwaggerFile)
	serverSwaggerUI(mux)

	fmt.Printf("Serving Grpc-Gateway on http://0.0.0.0%s\n", grpcGwAddr)
	logger.Fatal(http.ListenAndServe(grpcGwAddr, mux))
}

func serverSwaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		http.NotFound(w, r)
		return
	}

	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join("./api/", p)
	http.ServeFile(w, r, p)
}

func serverSwaggerUI(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "assets/swagger-ui/dist",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
