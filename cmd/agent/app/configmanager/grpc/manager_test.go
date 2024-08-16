// Copyright (c) 2018 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jaegertracing/jaeger/pkg/testutils"
	"github.com/jaegertracing/jaeger/proto-gen/api_v2"
)

func TestSamplingManager_GetSamplingStrategy(t *testing.T) {
	s, addr := initializeGRPCTestServer(t, func(s *grpc.Server) {
		api_v2.RegisterSamplingManagerServer(s, &mockSamplingHandler{})
	})
	conn, err := grpc.NewClient(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() { require.NoError(t, conn.Close()) })
	require.NoError(t, err)
	defer s.GracefulStop()
	manager := NewConfigManager(conn)
	resp, err := manager.GetSamplingStrategy(context.Background(), "any")
	require.NoError(t, err)
	assert.Equal(t, &api_v2.SamplingStrategyResponse{StrategyType: api_v2.SamplingStrategyType_PROBABILISTIC}, resp)
}

func TestSamplingManager_GetSamplingStrategy_error(t *testing.T) {
	conn, err := grpc.NewClient("foo", grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() { require.NoError(t, conn.Close()) })
	require.NoError(t, err)
	manager := NewConfigManager(conn)
	resp, err := manager.GetSamplingStrategy(context.Background(), "any")
	require.Nil(t, resp)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get sampling strategy")
}

func TestSamplingManager_GetBaggageRestrictions(t *testing.T) {
	manager := NewConfigManager(nil)
	rest, err := manager.GetBaggageRestrictions(context.Background(), "foo")
	require.Nil(t, rest)
	require.EqualError(t, err, "baggage not implemented")
}

type mockSamplingHandler struct{}

func (*mockSamplingHandler) GetSamplingStrategy(context.Context, *api_v2.SamplingStrategyParameters) (*api_v2.SamplingStrategyResponse, error) {
	return &api_v2.SamplingStrategyResponse{StrategyType: api_v2.SamplingStrategyType_PROBABILISTIC}, nil
}

func initializeGRPCTestServer(t *testing.T, beforeServe func(server *grpc.Server)) (*grpc.Server, net.Addr) {
	server := grpc.NewServer()
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	beforeServe(server)
	go func() {
		err := server.Serve(lis)
		require.NoError(t, err)
	}()
	return server, lis.Addr()
}

func TestMain(m *testing.M) {
	testutils.VerifyGoLeaks(m)
}
