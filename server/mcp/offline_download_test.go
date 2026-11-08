package mcp

import (
	"encoding/json"
	"testing"
)

func TestParseFSAddOfflineDownloadArgsRequiresURLsAndPath(t *testing.T) {
	_, err := parseFSAddOfflineDownloadArgs(json.RawMessage(`{"path":"/dl"}`))
	if err == nil || err.Message != "urls is required" {
		t.Fatalf("unexpected error: %+v", err)
	}
	_, err = parseFSAddOfflineDownloadArgs(json.RawMessage(`{"urls":["https://example.com/a.bin"]}`))
	if err == nil || err.Message != "path is required" {
		t.Fatalf("unexpected error: %+v", err)
	}
	_, err = parseFSAddOfflineDownloadArgs(json.RawMessage(`{"urls":["  ","\n"],"path":"/dl"}`))
	if err == nil || err.Message != "urls is required" {
		t.Fatalf("unexpected error: %+v", err)
	}
}

func TestParseFSAddOfflineDownloadArgsRejectsInvalidJSON(t *testing.T) {
	for _, raw := range []json.RawMessage{nil, json.RawMessage(`null`), json.RawMessage(`"bad"`)} {
		_, err := parseFSAddOfflineDownloadArgs(raw)
		if err == nil || err.Code != -32602 {
			t.Fatalf("raw=%s err=%+v", string(raw), err)
		}
	}
}

func TestParseFSAddOfflineDownloadArgs(t *testing.T) {
	args, err := parseFSAddOfflineDownloadArgs(json.RawMessage(`{
		"urls":[" https://example.com/a.bin ","","https://example.com/b.bin"],
		"path":"/downloads",
		"tool":"SimpleHttp",
		"delete_policy":"delete_on_upload_succeed"
	}`))
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if len(args.URLs) != 2 || args.URLs[0] != "https://example.com/a.bin" || args.URLs[1] != "https://example.com/b.bin" {
		t.Fatalf("urls=%v", args.URLs)
	}
	if args.Path != "/downloads" || args.Tool != "SimpleHttp" || args.DeletePolicy != "delete_on_upload_succeed" {
		t.Fatalf("args=%+v", args)
	}
}
