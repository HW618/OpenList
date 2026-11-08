package mcp

import (
	"encoding/json"
	"strings"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	offlinetool "github.com/OpenListTeam/OpenList/v4/internal/offline_download/tool"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/internal/task"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/OpenListTeam/OpenList/v4/server/handles"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type fsAddOfflineDownloadArgs struct {
	URLs         []string `json:"urls"`
	Path         string   `json:"path"`
	Tool         string   `json:"tool"`
	DeletePolicy string   `json:"delete_policy"`
}

func (s *Server) callFSAddOfflineDownload(c *gin.Context, raw json.RawMessage) (any, *rpcError) {
	args, mcpErr := parseFSAddOfflineDownloadArgs(raw)
	if mcpErr != nil {
		return nil, mcpErr
	}

	user, ok := c.Request.Context().Value(conf.UserKey).(*model.User)
	if !ok || user == nil {
		return nil, &rpcError{Code: -32603, Message: "missing user context"}
	}
	if !user.CanAddOfflineDownloadTasks() {
		return nil, &rpcError{Code: -32003, Message: "permission denied"}
	}

	reqPath, err := user.JoinPath(args.Path)
	if err != nil {
		return nil, &rpcError{Code: -32003, Message: err.Error()}
	}
	meta, err := op.GetNearestMeta(reqPath)
	if err != nil && !errors.Is(errors.Cause(err), errs.MetaNotFound) {
		return nil, &rpcError{Code: -32603, Message: err.Error()}
	}
	if !common.CanWrite(user, meta, reqPath) {
		return nil, &rpcError{Code: -32003, Message: errs.PermissionDenied.Error()}
	}

	var tasks []task.TaskExtensionInfo
	for _, u := range args.URLs {
		t, addErr := offlinetool.AddURL(c.Request.Context(), &offlinetool.AddURLArgs{
			URL:          u,
			DstDirPath:   reqPath,
			Tool:         args.Tool,
			DeletePolicy: offlinetool.DeletePolicy(args.DeletePolicy),
		})
		if addErr != nil {
			return nil, &rpcError{Code: -32603, Message: addErr.Error()}
		}
		if t != nil {
			tasks = append(tasks, t)
		}
	}

	return gin.H{
		"tasks": handles.TaskInfos(tasks),
	}, nil
}

func parseFSAddOfflineDownloadArgs(raw json.RawMessage) (*fsAddOfflineDownloadArgs, *rpcError) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, &rpcError{Code: -32602, Message: "invalid openlist.fs.add_offline_download arguments"}
	}
	var args fsAddOfflineDownloadArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, &rpcError{Code: -32602, Message: "invalid openlist.fs.add_offline_download arguments"}
	}
	cleaned := make([]string, 0, len(args.URLs))
	for _, u := range args.URLs {
		trimmed := strings.TrimSpace(u)
		if trimmed == "" {
			continue
		}
		cleaned = append(cleaned, trimmed)
	}
	args.URLs = cleaned
	if len(args.URLs) == 0 {
		return nil, &rpcError{Code: -32602, Message: "urls is required"}
	}
	if args.Path == "" {
		return nil, &rpcError{Code: -32602, Message: "path is required"}
	}
	return &args, nil
}
