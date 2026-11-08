package mcp

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestParseFSUploadArgsRequiresPath(t *testing.T) {
	_, err := parseFSUploadArgs(json.RawMessage(`{"as_task":true}`))
	if err == nil || err.Message != "path is required" {
		t.Fatalf("unexpected error: %+v", err)
	}
}

func TestParseFSUploadArgsDefaults(t *testing.T) {
	args, err := parseFSUploadArgs(json.RawMessage(`{"path":"/movies/a.bin"}`))
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if args.Path != "/movies/a.bin" || args.AsTask || !args.overwriteOrDefault() {
		t.Fatalf("args=%+v", args)
	}
}

func TestParseFSUploadArgsExplicitFlags(t *testing.T) {
	args, err := parseFSUploadArgs(json.RawMessage(`{"path":"/movies/a.bin","as_task":true,"overwrite":false,"password":"x"}`))
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if !args.AsTask || args.overwriteOrDefault() || args.Password != "x" {
		t.Fatalf("args=%+v", args)
	}
}

func TestBuildFSUploadGuide(t *testing.T) {
	guide := buildFSUploadGuide("/movies/demo file.bin", true, false)
	if guide.Method != "PUT" || guide.URL != "/api/fs/put" {
		t.Fatalf("guide=%+v", guide)
	}
	if guide.Headers["File-Path"] != url.PathEscape("/movies/demo file.bin") {
		t.Fatalf("File-Path=%q", guide.Headers["File-Path"])
	}
	if guide.Headers["As-Task"] != "true" || guide.Headers["Overwrite"] != "false" {
		t.Fatalf("headers=%v", guide.Headers)
	}
	if guide.Headers["Content-Type"] != "application/octet-stream" {
		t.Fatalf("content-type=%q", guide.Headers["Content-Type"])
	}
	if guide.Notes == "" {
		t.Fatal("expected notes")
	}
}

