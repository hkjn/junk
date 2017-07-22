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
	"strings"
	"time"

	googletime "github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "hkjn.me/prototest/report"
)

const defaultPort = ":50051"

var (
	debugging  = os.Getenv("REPORT_DEBUGGING") == "true"
	slackToken = os.Getenv("REPORT_SLACK_TOKEN")
)

type (
	clientInfo struct {
		lastSeen time.Time
		info     map[string]string
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
	go s.checkClients()
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
	debug("Slack replied: %v\n", resp)
	return nil
}

// getTime returns the time.Time equivalent of a timestamp proto message.
func getTime(t *googletime.Timestamp) time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

// Send implements report.ReportServer.
func (s *reportServer) Send(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	c, existed := s.clients[req.Name]
	info := []string{}
	for k, v := range req.Info {
		info = append(info, fmt.Sprintf("  `%s`: `%s`", k, v))
	}
	title := "Node"
	if !existed {
		title = "New node"
	}
	msg := fmt.Sprintf("%s `%s` reported to us (%s)", title, req.Name, strings.Join(info, ","))
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
	return &pb.Response{Message: resp}, nil
}

// checkClients runs forever, and periodically checks if any clients expired.
func (s *reportServer) checkClients() {
	log.Printf("Checking if any clients fell out of touch..")
	maxTime := time.Minute * 10
	for {
		time.Sleep(time.Second * 60)
		for name, v := range s.clients {
			if time.Since(v.lastSeen) > maxTime {
				msg := fmt.Sprintf(
					"Haven't heard from %q in %v, going to forget about them",
					name,
					time.Since(v.lastSeen),
				)
				log.Println(msg)
				delete(s.clients, name)
				sendSlack(msg)
			}
		}
	}
}

func main() {
	if slackToken == "" {
		log.Fatalf("No SLACK_TOKEN specified.\n")
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
