package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"lab_cli/cmd/helper"
	"lab_cli/service"
	"lab_cli/site"
	"os"

	buildbucketpb "go.chromium.org/luci/buildbucket/proto"
)

type Builder struct {
	Project string `json:"project"`
	Bucket  string `json:"bucket"`
	Builder string `json:"builder"`
}

type Time struct {
	Seconds    int64  `json:"seconds"`
	Nanos      int32  `json:"nanos"`
	TaipeiTime string `json:"taipei_time"`
}

type GetBuildResponse struct {
	Id         int64   `json:"id"`
	Builder    Builder `json:"builder"`
	CreatedBy  string  `json:"created_by"`
	CreateTime Time    `json:"create_time"`
	StartTime  Time    `json:"start_time"`
	EndTime    Time    `json:"end_time"`
	UpdateTime Time    `json:"update_time"`
	Status     string  `json:"status"`
	Input      string  `json:"input"`
}

func RunBuild(args []string) error {
	if len(args) == 0 {
		PrintBuildUsage()
		return fmt.Errorf("no build subcommand specified")
	}
	subcommand := args[0]
	subArgs := args[1:]
	switch subcommand {
	case "get":
		return RunBuildGet(subArgs)
	case "cancel":
		return RunBuildCancel(subArgs)
	case "-h", "--help", "help":
		PrintBuildUsage()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "Unknown build command: %q\n\n", subcommand)
		PrintBuildUsage()
		return fmt.Errorf("unknown build command: %s", subcommand)
	}
}

func PrintBuildUsage() {
	NewHelpBuilder().
		SetUsage("lab_cli build <command> [options]").
		AddCommand("get", "Get detailed build information").
		AddCommand("cancel", "Cancel a running build").
		AddOption("-h, --help", "Show help message").
		SetHintCommand("lab_cli build <command>").
		Print()
}

type BuildGetOptions struct {
	ID  int64
	Key string
}

func parseBuildGetFlags(args []string) (*BuildGetOptions, error) {
	fs := flag.NewFlagSet("build get", flag.ContinueOnError)
	opts := &BuildGetOptions{}

	fs.Int64Var(&opts.ID, "id", 0, "Build ID (required)")
	fs.StringVar(&opts.Key, "key", "", "Path to the service account JSON file (optional)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s lab_cli build get --id <ID> [options]\n\n%s\n", Header("Usage:"), Header("Options:"))
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if opts.ID == 0 {
		fs.Usage()
		return nil, fmt.Errorf("flag --id is required")
	}

	return opts, nil
}

func RunBuildGet(args []string) error {
	opts, err := parseBuildGetFlags(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	path := site.GetServiceAccountPath(opts.Key)
	if path == "" {
		return fmt.Errorf("cannot get service account")
	}

	buildBucketClient, err := service.NewBuildBucketClient(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to create buildbucket client: %w", err)
	}

	request := &buildbucketpb.GetBuildRequest{
		Id: opts.ID,
	}
	b, err := buildBucketClient.GetBuild(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to get build: %w", err)
	}

	res := GetBuildResponse{
		Id: opts.ID,
		Builder: Builder{
			Project: b.Builder.Project,
			Bucket:  b.Builder.Bucket,
			Builder: b.Builder.Builder,
		},
		CreatedBy: b.CreatedBy,
		CreateTime: Time{
			Seconds:    b.CreateTime.Seconds,
			Nanos:      b.CreateTime.Nanos,
			TaipeiTime: helper.UnixToTaipeiTime(b.CreateTime.Seconds, b.CreateTime.Nanos),
		},
		StartTime: Time{
			Seconds:    b.StartTime.Seconds,
			Nanos:      b.StartTime.Nanos,
			TaipeiTime: helper.UnixToTaipeiTime(b.StartTime.Seconds, b.StartTime.Nanos),
		},
		EndTime: Time{
			Seconds:    b.EndTime.Seconds,
			Nanos:      b.EndTime.Nanos,
			TaipeiTime: helper.UnixToTaipeiTime(b.EndTime.Seconds, b.EndTime.Nanos),
		},
		UpdateTime: Time{
			Seconds:    b.UpdateTime.Seconds,
			Nanos:      b.UpdateTime.Nanos,
			TaipeiTime: helper.UnixToTaipeiTime(b.UpdateTime.Seconds, b.UpdateTime.Nanos),
		},
		Status: helper.ToBuildStatus(b.Status),
		Input:  b.Input.String(),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	return enc.Encode(res)
}

type BuildCancelOptions struct {
	ID int64
}

func parseBuildCancelFlags(args []string) (*BuildCancelOptions, error) {
	fs := flag.NewFlagSet("build cancel", flag.ContinueOnError)
	opts := &BuildCancelOptions{}

	fs.Int64Var(&opts.ID, "id", 0, "Build ID (required)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s lab_cli build cancel --id <ID> [options]\n\n%s\n", Header("Usage:"), Header("Options:"))
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if opts.ID == 0 {
		fs.Usage()
		return nil, fmt.Errorf("flag --id is required")
	}

	return opts, nil
}

func RunBuildCancel(args []string) error {
	opts, err := parseBuildCancelFlags(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	// TODO:
	_ = opts
	return nil
}
