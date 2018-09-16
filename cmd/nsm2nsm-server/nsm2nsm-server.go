package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sbezverk/nsm2nsm/pkg/apis/serverendpoints"
	customresourceclient "github.com/sbezverk/nsm2nsm/pkg/client/clientset/versioned"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file.")
)

type serverEndpointsService struct {
	k8s        *kubernetes.Clientset
	crClient   *customresourceclient.Clientset
	grpcServer *grpc.Server
}

func newServerEndpointsService(k8s *kubernetes.Clientset, crClient *customresourceclient.Clientset) *serverEndpointsService {
	return &serverEndpointsService{
		k8s:        k8s,
		crClient:   crClient,
		grpcServer: grpc.NewServer(),
	}
}

func (s serverEndpointsService) ServerPingPong(ctx context.Context, ping *serverendpoints.ServerPing) (*serverendpoints.ServerPong, error) {
	return &serverendpoints.ServerPong{
		Data: strings.ToUpper(ping.Data),
	}, nil
}

func k8sBuildClient() (*kubernetes.Clientset, *customresourceclient.Clientset, error) {
	var config *rest.Config
	var err error

	if *kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create config: %v", err)
	}
	k8s, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client: %v", err)
	}
	customResourceClient, err := customresourceclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create customresource clientset: %v", err)
	}

	return k8s, customResourceClient, nil
}

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	glog.Infof("Starting Server Endpoints Service...")

	coreClient, customResourceClient, err := k8sBuildClient()
	if err != nil {
		glog.Fatalf("Shutting down Server Endpoints Service with error: %+v", err)
	}

	lis, err := net.Listen("tcp", ":14141")
	if err != nil {
		glog.Errorf("Error. Starting SRIOV Network Device Plugin server failed: %v", err)
	}
	server := newServerEndpointsService(coreClient, customResourceClient)

	serverendpoints.RegisterServerPingPongServer(server.grpcServer, server)
	go func() {
		if err := server.grpcServer.Serve(lis); err != nil {
			glog.Fatalf("PingPong Server failed to listen on port with error: %+v", err)
			os.Exit(1)
		}
	}()

	// respond to syscalls for termination
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Catch termination signals
	select {
	case sig := <-sigChannel:
		glog.Infof("Received signal \"%v\", shutting down.", sig)
		return
	}
}
