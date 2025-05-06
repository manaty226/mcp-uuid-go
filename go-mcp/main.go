package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"hello-world-uuid",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	t := mcp.NewTool(
		"uuid-tool",
		mcp.WithDescription("A tool that responses a UUID value"),
		mcp.WithNumber("version",
			mcp.Required(),
			mcp.Description("The version of UUID which is restricted to 4 or 7"),
		),
	)

	s.AddTool(t, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		version, ok := request.Params.Arguments["version"].(float64)
		if !ok {
			fmt.Println(reflect.TypeOf(request.Params.Arguments["version"]))
			return nil, fmt.Errorf("version parameter is required")
		}

		switch version {
		case 4:
			return mcp.NewToolResultText(uuid.New().String()), nil
		case 7:
			id, err := uuid.NewV7()
			if err != nil {
				return nil, fmt.Errorf("error generating UUIDv7: %v", err)
			}
			return mcp.NewToolResultText(id.String()), nil
		default:
			return nil, fmt.Errorf("unsupported UUID version: %v", version)
		}
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
