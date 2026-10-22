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
	"google.golang.org/protobuf/types/known/timestamppb"
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
	Id         int64          `json:"id"`
	Builder    Builder        `json:"builder"`
	CreatedBy  string         `json:"created_by"`
	CreateTime *Time          `json:"create_time"`
	StartTime  *Time          `json:"start_time"`
	EndTime    *Time          `json:"end_time"`
	UpdateTime *Time          `json:"update_time"`
	Status     string         `json:"status"`
	Input      map[string]any `json:"input,omitempty"`
}

func toTime(ts *timestamppb.Timestamp) *Time {
	if ts == nil {
		return nil
	}
	return &Time{
		Seconds:    ts.Seconds,
		Nanos:      ts.Nanos,
		TaipeiTime: helper.UnixToTaipeiTime(ts.Seconds, ts.Nanos),
	}
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
		Mask: &buildbucketpb.BuildMask{
			AllFields: true,
		},
	}

	b, err := buildBucketClient.GetBuild(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to get build: %w", err)
	}

	var builder Builder
	if b.Builder != nil {
		builder = Builder{
			Project: b.Builder.Project,
			Bucket:  b.Builder.Bucket,
			Builder: b.Builder.Builder,
		}
	}

	var inputMap map[string]any
	if b.GetInput().GetProperties() != nil {
		inputMap = b.GetInput().GetProperties().AsMap()
	}

	res := GetBuildResponse{
		Id:         opts.ID,
		Builder:    builder,
		CreatedBy:  b.CreatedBy,
		CreateTime: toTime(b.CreateTime),
		StartTime:  toTime(b.StartTime),
		EndTime:    toTime(b.EndTime),
		UpdateTime: toTime(b.UpdateTime),
		Status:     helper.ToBuildStatus(b.Status),
		Input:      inputMap,
	}

	// For example, `res.Input["$chromeos/bad_builds"]` to check if the build is a bad build.

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	return enc.Encode(res)
}

type BuildCancelOptions struct {
	ID     int64
	Reason string
	Key    string
}

func parseBuildCancelFlags(args []string) (*BuildCancelOptions, error) {
	fs := flag.NewFlagSet("build cancel", flag.ContinueOnError)
	opts := &BuildCancelOptions{}

	fs.Int64Var(&opts.ID, "id", 0, "Build ID (required)")
	fs.StringVar(&opts.Reason, "reason", "", "Reason for cancelling the build (required)")
	fs.StringVar(&opts.Key, "key", "", "Path to the service account JSON file (optional)")

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

	if opts.Reason == "" {
		fs.Usage()
		return nil, fmt.Errorf("flag --reason is required")
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

	request := &buildbucketpb.CancelBuildRequest{
		Id:              opts.ID,
		SummaryMarkdown: opts.Reason,
	}

	b, err := buildBucketClient.CancelBuild(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to cancel build: %w", err)
	}

	fmt.Printf("Build %d cancelled successfully: %v\n", opts.ID, b)
	return nil
}
