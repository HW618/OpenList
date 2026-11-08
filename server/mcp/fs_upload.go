package mcp

import (
	"encoding/json"
	"net/url"
	stdpath "path"
	"strconv"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type fsUploadArgs struct {
	Path      string `json:"path"`
	AsTask    bool   `json:"as_task"`
	Overwrite *bool  `json:"overwrite"`
	Password  string `json:"password"`
}

type fsUploadGuide struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Path    string            `json:"path"`
	Notes   string            `json:"notes"`
}

func (a *fsUploadArgs) overwriteOrDefault() bool {
	if a.Overwrite == nil {
		return true
	}
	return *a.Overwrite
}

func (s *Server) callFSUpload(c *gin.Context, raw json.RawMessage) (any, *rpcError) {
	args, mcpErr := parseFSUploadArgs(raw)
	if mcpErr != nil {
		return nil, mcpErr
	}

	user, ok := c.Request.Context().Value(conf.UserKey).(*model.User)
	if !ok || user == nil {
		return nil, &rpcError{Code: -32603, Message: "missing user context"}
	}
	if user.IsGuest() && user.Disabled {
		return nil, &rpcError{Code: -32001, Message: "guest user is disabled"}
	}

	reqPath, err := user.JoinPath(args.Path)
	if err != nil {
		return nil, &rpcError{Code: -32003, Message: err.Error()}
	}
	parentPath := stdpath.Dir(reqPath)
	parentMeta, err := op.GetNearestMeta(parentPath)
	if err != nil && !errors.Is(errors.Cause(err), errs.MetaNotFound) {
		return nil, &rpcError{Code: -32603, Message: err.Error()}
	}
	if !common.CanAccess(user, parentMeta, parentPath, args.Password) {
		return nil, &rpcError{Code: -32003, Message: "password is incorrect or you have no permission"}
	}
	if !user.CanWriteContent() && !common.CanWriteContentBypassUserPerms(parentMeta, parentPath) {
		return nil, &rpcError{Code: -32003, Message: errs.PermissionDenied.Error()}
	}
	if !common.CanWrite(user, parentMeta, parentPath) {
		return nil, &rpcError{Code: -32003, Message: errs.PermissionDenied.Error()}
	}

	return buildFSUploadGuide(args.Path, args.AsTask, args.overwriteOrDefault()), nil
}

func buildFSUploadGuide(path string, asTask, overwrite bool) fsUploadGuide {
	return fsUploadGuide{
		Method: "PUT",
		URL:    "/api/fs/put",
		Path:   path,
		Headers: map[string]string{
			"File-Path":    url.PathEscape(path),
			"As-Task":      strconv.FormatBool(asTask),
			"Overwrite":    strconv.FormatBool(overwrite),
			"Content-Type": "application/octet-stream",
		},
		Notes: "Stream the file body to this endpoint with the same Authorization (or token) used for MCP. Do not send file bytes via tools/call.",
	}
}

func parseFSUploadArgs(raw json.RawMessage) (*fsUploadArgs, *rpcError) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, &rpcError{Code: -32602, Message: "invalid openlist.fs.upload arguments"}
	}
	var args fsUploadArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, &rpcError{Code: -32602, Message: "invalid openlist.fs.upload arguments"}
	}
	if args.Path == "" {
		return nil, &rpcError{Code: -32602, Message: "path is required"}
	}
	return &args, nil
}
