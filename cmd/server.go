package cmd

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/giantswarm/fulfillment/aws"
	"github.com/giantswarm/fulfillment/handlers"
	"github.com/giantswarm/fulfillment/slack"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Run:   runServer,
}

var (
	awsAccessKeyId     string
	awsSecretAccessKey string
	slackToken         string

	mockAws   bool
	mockSlack bool
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().StringVar(&awsAccessKeyId, "aws-access-key-id", "", "AWS access key id (or set AWS_ACCESS_KEY_ID)")
	serverCmd.PersistentFlags().StringVar(&awsSecretAccessKey, "aws-secret-access-key", "", "AWS secret access key (or set AWS_SECRET_ACCESS_KEY)")

	serverCmd.PersistentFlags().StringVar(&slackToken, "slack-token", "", "Slack API token (or set SLACK_TOKEN)")

	serverCmd.PersistentFlags().BoolVar(&mockAws, "mock-aws", false, "Mock calls to AWS")
	serverCmd.PersistentFlags().BoolVar(&mockSlack, "mock-slack", false, "Mock calls to Slack")
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func runServer(cmd *cobra.Command, args []string) {
	if awsAccessKeyId == "" {
		awsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if awsSecretAccessKey == "" {
		awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}
	if slackToken == "" {
		slackToken = os.Getenv("SLACK_TOKEN")
	}

	flag.Parse()

	awsService, err := aws.New(awsAccessKeyId, awsSecretAccessKey, mockAws)
	if err != nil {
		log.Fatalf("failed to create aws service: %s", err)
	}

	slackService, err := slack.New(slackToken, mockSlack)
	if err != nil {
		log.Fatalf("failed to create slack service: %s", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Root(w, r, awsService, slackService)
	})
	mux.HandleFunc("/success", handlers.Success)
	mux.HandleFunc("/webhook", handlers.Webhook)
	mux.Handle("/content/", http.StripPrefix("/content/", http.FileServer(http.Dir("./content"))))

	loggedMux := loggingMiddleware(mux)

	server := &http.Server{
		Addr:              ":8000",
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           loggedMux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}
