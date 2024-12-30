package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/alecthomas/kong"
)

var CLI struct {
	List struct{} `cmd:"" help:"List hosts."`
	Mute struct {
		Duration int `help:"How many minutes to mute the host for." default:"60"`

		Hosts []string `arg:"" name:"host" help:"Hosts to mute. Note that hosts that don't exist do not generate an error from Datadog."`
	} `cmd:"" help:"Mute hosts."`

	Unmute struct {
		Hosts []string `arg:"" name:"host" help:"Hosts to mute."`
	} `cmd:"" help:"Unmute hosts."`
}

func newDatadogApi() (context.Context, *datadogV1.HostsApi) {
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV1.NewHostsApi(apiClient)
	return ctx, api
}

func main() {
	ctx := kong.Parse(&CLI, kong.Description("Quickly mute and unmute hosts in Datadog."))
	switch ctx.Command() {
	case "list":
		ctx, api := newDatadogApi()
		resp, r, err := api.ListHosts(ctx, *datadogV1.NewListHostsOptionalParameters())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `HostsApi.ListHosts`: %v\n", err)
			if r != nil {
				fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
			}
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Hostname\tMute Status\t")
		for _, host := range resp.HostList {
			muteStatus := "unmuted"
			if host.MuteTimeout.IsSet() && host.MuteTimeout.Get() != nil {
				muteTime := time.Unix(*host.MuteTimeout.Get(), 0)
				muteStatus = "until " + muteTime.Format(time.RFC3339)
			}

			fmt.Fprintf(w, "%s\t%s\t\n", *host.HostName, muteStatus)
		}
		w.Flush()

	case "mute <host>":
		ctx, api := newDatadogApi()

		body := datadogV1.HostMuteSettings{
			End:      datadog.PtrInt64(time.Now().Add(time.Duration(CLI.Mute.Duration) * time.Minute).Unix()),
			Message:  datadog.PtrString("Host muted via dd-mute-host"),
			Override: datadog.PtrBool(false),
		}

		for _, host_name := range CLI.Mute.Hosts {
			resp, r, err := api.MuteHost(ctx, host_name, body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to mute host %v: %v %v\n", host_name, err, r.Body)
				continue
			}
			fmt.Fprintf(os.Stdout, "resp %v \n", resp)
			fmt.Fprintf(os.Stdout, "Host %v status: %v\n", host_name, resp.GetAction())
		}

	case "unmute <host>":
		ctx, api := newDatadogApi()

		for _, host_name := range CLI.Unmute.Hosts {
			resp, r, err := api.UnmuteHost(ctx, host_name)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to unmute host %v: %v %v\n", host_name, err, r.Body)
				continue
			}

			fmt.Fprintf(os.Stdout, "Host %v status: %v\n", host_name, resp.GetAction())
		}
	default:
		panic(ctx.Command())
	}
}
