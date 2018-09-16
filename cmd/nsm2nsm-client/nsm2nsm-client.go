package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workspace/ligato/nsm2nsm/pkg/apis/serverendpoints"
	customresourceclient "workspace/ligato/nsm2nsm/pkg/client/clientset/versioned"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file.")
)

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

func getServerEndpointAddress(coreClient *kubernetes.Clientset,
	customResourceClient *customresourceclient.Clientset) (string, error) {

	selector := labels.SelectorFromSet(labels.Set(map[string]string{"nsm-network-service-name": "pingServer"}))
	options := metav1.ListOptions{LabelSelector: selector.String()}
	endpointsList, err := customResourceClient.SbezverkV1().ServerEndpoints("default").List(options)
	if err != nil {
		return "", err
	}
	if len(endpointsList.Items) == 0 {
		return "", fmt.Errorf("no network service name 'pingServer' is found")
	}
	endpoint, err := customResourceClient.SbezverkV1().ServerEndpoints("default").Get(endpointsList.Items[0].ObjectMeta.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return endpoint.Spec.ServerAddress, nil
}

func dial(ctx context.Context, serverAddress string) (*grpc.ClientConn, error) {
	c, err := grpc.DialContext(ctx, serverAddress, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("tcp", addr, timeout)
		}),
	)

	return c, err
}

// socketOperationCheck checks for liveness of a gRPC server socket.
func socketOperationCheck(serverAddress string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := dial(ctx, serverAddress)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func pingServer(clientPingPong serverendpoints.ServerPingPongClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ping := &serverendpoints.ServerPing{
		Data: "abcdefghijklmnopqrstuvwxyz",
	}
	pong, err := clientPingPong.ServerPingPong(ctx, ping)
	if err != nil {
		return err
	}
	glog.Infof("Client sent this: %s", ping.Data)
	glog.Infof("Server returned this: %s", pong.Data)

	return nil
}

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	glog.Infof("Starting Server Endpoints Service...")

	coreClient, customResourceClient, err := k8sBuildClient()
	if err != nil {
		glog.Fatalf("Shutting down Server Endpoints Service with error: %+v", err)
	}
	serverAddress, err := getServerEndpointAddress(coreClient, customResourceClient)
	if err != nil {
		glog.Fatalf("fail to get Server Endpoint address with error: %+v", err)
		os.Exit(1)
	}

	// Need to check if server address is actual IP address or DNS name
	// gRPC requires to add dns:// prefix
	host, _, err := net.SplitHostPort(serverAddress)
	if err != nil {
		glog.Fatalf("fail to parse server address %s with error: %+v", serverAddress, err)
		os.Exit(1)
	}
	if net.ParseIP(host) == nil {
		// so host is not a valid IP address, which means it could be a dns name
		// working: serverAddress = "dns:///" + host + ".default.svc" + ":" + port
		serverAddress = "dns:///" + serverAddress
	}

	glog.Infof("Server endpoint address: %s", serverAddress)

	// Setting up gRPC client connection to discovered server address
	conn, err := socketOperationCheck(serverAddress)
	if err != nil {
		glog.Fatalf("fail to set up connection to %s with error: %+v", serverAddress, err)
		os.Exit(1)
	}
	pingServerClient := serverendpoints.NewServerPingPongClient(conn)
	// respond to syscalls for termination
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	pingTicker := time.NewTicker(2 * time.Second)
	for {
		select {
		case sig := <-sigChannel:
			glog.Infof("Received signal \"%v\", shutting down.", sig)
			os.Exit(0)
		case <-pingTicker.C:
			if err := pingServer(pingServerClient); err != nil {
				glog.Errorf("ping failed with error: %+v", err)
				continue
			}
		}
	}
}
