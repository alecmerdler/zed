//go:build wasm
// +build wasm

// To Run:
// 1) Install wasmbrowsertest: `go install github.com/agnivade/wasmbrowsertest@latest`
// 2) Run: `GOOS=js GOARCH=wasm go test -exec wasmbrowsertest`

package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	v1 "github.com/authzed/spicedb/pkg/proto/core/v1"
	devinterface "github.com/authzed/spicedb/pkg/proto/developer/v1"
	"github.com/authzed/spicedb/pkg/tuple"
)

func TestZedCommand(t *testing.T) {
	requestCtx := &devinterface.RequestContext{
		Schema: `definition user {}
		
		caveat somecaveat(somecondition int) {
			somecondition == 42
		}

		definition document {
			relation viewer: user | user with somecaveat
			permission view = viewer
		}`,

		Relationships: []*v1.RelationTuple{
			tuple.MustParse(`document:first#viewer@user:fred[somecaveat:{"somecondition": 42}]`).ToCoreTuple(),
			tuple.MustParse("document:first#viewer@user:tom").ToCoreTuple(),
		},
	}

	encodedContext, err := protojson.Marshal(requestCtx)
	require.NoError(t, err)

	rootCmd := buildRootCmd()

	// Run with --help
	result := runZedCommand(rootCmd, string(encodedContext), []string{"permission", "check", "--help"})
	require.Contains(t, result.Output, "Usage:", "failed to run 'permission check --help' command: %v", result)

	// Run the actual command.
	result = runZedCommand(rootCmd, string(encodedContext), []string{"permission", "check", "document:first", "view", "user:fred"})
	require.True(t, strings.HasSuffix(strings.TrimSpace(result.Output), "true"), "expected true at end of: %s", result.Output)

	updatedContext := &devinterface.RequestContext{}
	err = protojson.Unmarshal([]byte(result.UpdatedContext), updatedContext)
	require.NoError(t, err)

	require.Contains(t, updatedContext.Schema, "definition document")
	require.Equal(t, `document:first#viewer@user:fred[somecaveat:{"somecondition":42}]`, tuple.MustCoreRelationToString(updatedContext.Relationships[0]))
	require.Equal(t, "document:first#viewer@user:tom", tuple.MustCoreRelationToString(updatedContext.Relationships[1]))
	require.Len(t, updatedContext.Relationships, 2)

	// Run the actual command.
	result = runZedCommand(rootCmd, string(encodedContext), []string{"relationship", "create", "document:1", "viewer", "user:1"})
	require.Empty(t, result.Error, "failed to run relationship create: %s", result.Error)

	updatedContext = &devinterface.RequestContext{}
	err = protojson.Unmarshal([]byte(result.UpdatedContext), updatedContext)
	require.NoError(t, err)
	require.Len(t, updatedContext.Relationships, 3)
}
