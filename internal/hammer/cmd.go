package hammer

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"

	"github.com/jimmykodes/kodes/internal/command"
)

var hammerCmd = &cobra.Command{
	Use:   "hammer",
	Short: "Send a whole bunch of http requests. Hammer the server.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	Run: run,
}

func New() *command.Command {
	return command.New(hammerCmd, flagInit)
}

const (
	bodyString    = "body-string"
	bodyFile      = "body-file"
	contentType   = "content-type"
	url           = "url"
	headers       = "headers"
	method        = "method"
	count         = "count"
	maxConcurrent = "max-concurrent"
)

func flagInit() error {
	// Add flags here
	hammerCmd.Flags().StringP(bodyString, "s", "", "string to use for http request body")
	hammerCmd.Flags().StringP(bodyFile, "f", "", "file to use for http request body")
	hammerCmd.Flags().StringP(contentType, "t", "application/json", "content type header to set on request")
	hammerCmd.Flags().StringP(url, "u", "", "url to send request to")
	hammerCmd.Flags().StringP(method, "m", http.MethodGet, "request method")
	hammerCmd.Flags().StringToStringP(headers, "H", nil, "additional headers to send on request")
	hammerCmd.Flags().IntP(count, "c", 10, "number of times to send request")
	hammerCmd.Flags().Int64P(maxConcurrent, "M", 10, "max number of concurrent requests")

	return nil
}

func run(cmd *cobra.Command, args []string) {
	maxProcs := viper.GetInt64(maxConcurrent)
	sem := semaphore.NewWeighted(maxProcs)
	var body *template.Template
	if bs := viper.GetString(bodyString); bs != "" {
		body = template.Must(template.New("body").Parse(bs))
	} else if bf := viper.GetString(bodyFile); bf != "" {
		f, err := os.Open(bf)
		if err != nil {
			log.Println("error opening body file", err)
			return
		}
		bodyData, err := io.ReadAll(f)
		if err != nil {
			log.Println("error reading file", err)
			return
		}
		body = template.Must(template.New("body").Parse(string(bodyData)))
	}
	for i := 0; i < viper.GetInt(count); i++ {
		go func(index int) {
			err := makeReq(cmd.Context(), sem, body, index)
			if err != nil {
				log.Println("error making request", err)
				return
			}
		}(i)
	}
	time.Sleep(time.Second * 5)
	err := sem.Acquire(cmd.Context(), maxProcs)
	if err != nil {
		log.Println("error gathering semaphore", err)
		return
	}
}

type data struct {
	Index int
}

func makeReq(ctx context.Context, sem *semaphore.Weighted, body *template.Template, index int) error {
	if err := sem.Acquire(ctx, 1); err != nil {
		return err
	}
	defer sem.Release(1)
	var sb strings.Builder
	if body != nil {
		err := body.Execute(&sb, data{Index: index})
		if err != nil {
			return err
		}
	}
	log.Println("sending data", sb.String())

	req, err := http.NewRequest(viper.GetString(method), viper.GetString(url), bytes.NewBuffer([]byte(sb.String())))
	if err != nil {
		return err
	}

	for k, v := range viper.GetStringMapString(headers) {
		req.Header.Set(k, v)
	}
	log.Println("making request to", viper.GetString(url))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(respData))
	return nil
}
