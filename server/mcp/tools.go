package mcp

import "encoding/json"

type tool struct {
	Name        string          `json:"name"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	InputSchema toolInputSchema `json:"inputSchema"`
}

type toolInputSchema struct {
	Type       string                    `json:"type"`
	Properties map[string]schemaProperty `json:"properties,omitempty"`
	Required   []string                  `json:"required,omitempty"`
}

type schemaProperty struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type toolsListParams struct {
	Cursor string `json:"cursor,omitempty"`
}

var openListTools = []tool{
	{
		Name:        "openlist.fs.list",
		Title:       "OpenList FS List",
		Description: "List files and directories under a mount path that the current user can access.",
		InputSchema: toolInputSchema{
			Type: "object",
			Properties: map[string]schemaProperty{
				"path": {
					Type:        "string",
					Description: "Mount path to list, for example \"/\" or \"/movies\".",
				},
				"refresh": {
					Type:        "boolean",
					Description: "Refresh the directory listing before returning results.",
				},
				"password": {
					Type:        "string",
					Description: "Optional password for protected paths.",
				},
				"page": {
					Type:        "integer",
					Description: "1-based page number.",
				},
				"per_page": {
					Type:        "integer",
					Description: "Page size.",
				},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "openlist.fs.get",
		Title:       "OpenList FS Get",
		Description: "Get file or directory details for a mount path that the current user can access.",
		InputSchema: toolInputSchema{
			Type: "object",
			Properties: map[string]schemaProperty{
				"path": {
					Type:        "string",
					Description: "Mount path to inspect, for example \"/movies/demo.mp4\".",
				},
				"password": {
					Type:        "string",
					Description: "Optional password for protected paths.",
				},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "openlist.fs.link",
		Title:       "OpenList FS Link",
		Description: "Return usable link information for a file path that the current user can access.",
		InputSchema: toolInputSchema{
			Type: "object",
			Properties: map[string]schemaProperty{
				"path": {
					Type:        "string",
					Description: "File mount path, for example \"/movies/demo.mp4\".",
				},
				"password": {
					Type:        "string",
					Description: "Optional password for protected paths.",
				},
				"type": {
					Type:        "string",
					Description: "Optional link type forwarded to storage drivers.",
				},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "openlist.fs.add_offline_download",
		Title:       "OpenList Add Offline Download",
		Description: "Create offline download tasks into a writable mount directory (same semantics as POST /api/fs/add_offline_download).",
		InputSchema: toolInputSchema{
			Type: "object",
			Properties: map[string]schemaProperty{
				"urls": {
					Type:        "array",
					Description: "One or more download URLs (magnet/http/etc. depending on configured tools).",
				},
				"path": {
					Type:        "string",
					Description: "Destination directory mount path.",
				},
				"tool": {
					Type:        "string",
					Description: "Offline download tool name, for example \"SimpleHttp\", \"aria2\", \"qBittorrent\".",
				},
				"delete_policy": {
					Type:        "string",
					Description: "Temp-file delete policy: delete_on_upload_succeed, delete_on_upload_failed, delete_never, delete_always, upload_download_stream.",
				},
			},
			Required: []string{"urls", "path"},
		},
	},
	{
		Name:        "openlist.fs.upload",
		Title:       "OpenList FS Upload Guide",
		Description: "Validate upload permission and return streaming PUT instructions for /api/fs/put. Does not accept file bytes over MCP.",
		InputSchema: toolInputSchema{
			Type: "object",
			Properties: map[string]schemaProperty{
				"path": {
					Type:        "string",
					Description: "Full destination file mount path, for example \"/movies/demo.mp4\".",
				},
				"as_task": {
					Type:        "boolean",
					Description: "When true, upload runs as a background task (As-Task header).",
				},
				"overwrite": {
					Type:        "boolean",
					Description: "When false, refuse if the file already exists. Defaults to true.",
				},
				"password": {
					Type:        "string",
					Description: "Optional password for protected parent directories.",
				},
			},
			Required: []string{"path"},
		},
	},
}

func (s *Server) handleToolsList(req request) response {
	var params toolsListParams
	if len(req.Params) > 0 {
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &rpcError{Code: -32602, Message: "invalid tools/list params"},
			}
		}
	}

	return response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]any{
			"tools": openListTools,
		},
	}
}
