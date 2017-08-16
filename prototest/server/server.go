// server.go implements a GRPC Report server.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"os"
	"time"

	googletime "github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "hkjn.me/junk/prototest/report"
)

const defaultPort = ":50051"

var (
	debugging  = os.Getenv("REPORT_DEBUGGING") == "true"
	slackToken = os.Getenv("REPORT_SLACK_TOKEN")
)

type (
	// clientInfo represents info about a client that has
	// reported to us.
	//
	clientInfo struct {
		// lastSeen is the last time we heard from the client.
		lastSeen time.Time
		// info is extra info reported by the client.
		info *pb.ClientInfo
	}
	// reportServer is used to implement report.GreeterServer.
	reportServer struct {
		// clients is the known clients.
		clients map[string]clientInfo
	}
)

func debug(format string, a ...interface{}) {
	if !debugging {
		return
	}
	log.Printf(format, a...)
}

// newRpcServer returns the GRPC server.
func newRpcServer() *grpc.Server {
	rpcServer := grpc.NewServer()
	s := &reportServer{map[string]clientInfo{}}
	pb.RegisterReportServer(rpcServer, s)
	log.Printf("Registering GreeterServer to tcp listener on %q..\n", defaultPort)
	reflection.Register(rpcServer)
	return rpcServer
}

// sendSlack sends msg to Slack.
func sendSlack(msg string) error {
	slackUrl := "https://hooks.slack.com/services/" + slackToken
	data := struct {
		Text      string `json:"text"`
		LinkNames uint   `json:"link_names"`
		// TODO: Find reason icon_emoji seems to be ignored.
		// IconEmoji string `json:"icon_emoji"`
	}{
		Text:      fmt.Sprintf("`[report_server]`", msg),
		LinkNames: 1,
		// IconEmoji: ":heavy_exclamation_mark:",
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode as json: %v\n", err)
		return err
	}
	debug("Sending request to %q: %s\n", slackUrl, b)
	resp, err := http.Post(slackUrl, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Failed to send to Slack: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	debug("Slack replied: %v\n", resp)
	return nil
}

// getTime returns the time.Time equivalent of a timestamp proto message.
func getTime(t *googletime.Timestamp) time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

func (s *reportServer) Info(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {
	log.Printf("Received info request\n")
	return &pb.InfoResponse{
		Info: map[string]*pb.ClientInfo{
			"notimplementedyet": &pb.ClientInfo{
				CpuArch: "gelatinous",
				Hostname: "notimplementedyet-inforesponse",
			},
		},
	}, nil
}

// getInfo describes the client info as a string.
func getInfo(info *pb.ClientInfo) string {
	extra := []string{}
	if info.CpuArch != "" {
		extra = append(extra, fmt.Sprintf("`%s`", info.CpuArch))
	}
	if info.KernelName != "" {
		extra = append(extra, fmt.Sprintf("`%s %s`", info.KernelName, info.KernelVersion))
	}
	if info.Platform != "" {
		extra = append(extra, fmt.Sprintf("`%s`", info.Platform))
	}
	for i := range info.Tags {
		extra = append(extra, fmt.Sprintf("`%s`", info.Tags[i]))
	}
	return fmt.Sprintf("`%s` (%s)", info.Hostname, strings.Join(extra, ", "))
}

// Send implements report.ReportServer.
func (s *reportServer) Send(ctx context.Context, req *pb.ReportRequest) (*pb.ReportResponse, error) {
	c, existed := s.clients[req.Name]
	title := "Node"
	if !existed {
		title = "New node"
	}
	msg := fmt.Sprintf("%s reported to us: %s", title, getInfo(req.Info))
	// TODO: Validate data; seems like it can become corrupt:
	// Full info: allowed_ssh_keys:"memory_total_mb" cpu_arch:"7867"

	log.Println(msg)
	log.Printf("Full info: %+v\n", req.Info)
	if existed {
		log.Printf("Heard from known client for the first time in %v: %s\n", time.Since(c.lastSeen), msg)
	} else {
		log.Printf("Heard from new client: %s\n", msg)
		sendSlack(msg)
	}
	c = clientInfo{
		lastSeen: getTime(req.Ts),
		info:     req.Info,
	}
	s.clients[req.Name] = c
	resp := fmt.Sprintf(
		"Hello %q, thanks for writing me at %v, it is now %v.",
		req.Name,
		getTime(req.Ts),
		time.Now().Unix(),
	)
	log.Printf("Responding to client %q: %q..\n", req.Name, resp)
	return &pb.ReportResponse{Message: resp}, nil
}

func main() {
	log.Printf("The report_server %s starting..\n", Version)
	if slackToken == "" {
		log.Println("No REPORT_SLACK_TOKEN specified, can't report to Slack.")
	}
	rpcServer := newRpcServer()
	lis, err := net.Listen("tcp", defaultPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	sendSlack(fmt.Sprintf("%s `report_server` starting..", Version))
	if err := rpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
