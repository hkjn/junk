// server.go implements a GRPC Report server.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
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
		Text:      msg,
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

// Send implements report.ReportServer.
func (s *reportServer) Send(ctx context.Context, req *pb.ReportRequest) (*pb.ReportResponse, error) {
	c, existed := s.clients[req.Name]
	info := fmt.Sprintf("`%s` (`%s`)", c.info.Hostname, c.info.CpuArch)
	title := "Node"
	if !existed {
		title = "New node"
	}
	msg := fmt.Sprintf("%s `%s` reported to us (%s)", title, req.Name, info)
	log.Println(msg)
	if existed {
		log.Printf("Heard from known client for the first time in %v: %s\n", time.Since(c.lastSeen), msg)
	} else {
		log.Printf("Heard from new client: %s\n", msg)
		sendSlack(msg)
	}
	s.clients[req.Name] = clientInfo{
		lastSeen: getTime(req.Ts),
		info:     req.Info,
	}
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
	if slackToken == "" {
		log.Fatalf("No REPORT_SLACK_TOKEN specified.\n")
	}
	rpcServer := newRpcServer()
	lis, err := net.Listen("tcp", defaultPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := rpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
