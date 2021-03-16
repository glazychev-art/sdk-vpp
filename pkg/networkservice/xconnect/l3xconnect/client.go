// Copyright (c) 2021 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package l3xconnect

import (
	"context"

	"git.fd.io/govpp.git/api"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"
)

type l3XConnectClient struct {
	vppConn api.Connection
}

// NewClient returns a Client chain element that will l3 cross connect a client and server vpp interface (if present)
func NewClient(vppConn api.Connection) networkservice.NetworkServiceClient {
	return &l3XConnectClient{
		vppConn: vppConn,
	}
}

func (v *l3XConnectClient) Request(ctx context.Context, request *networkservice.NetworkServiceRequest, opts ...grpc.CallOption) (*networkservice.Connection, error) {
	conn, err := next.Client(ctx).Request(ctx, request, opts...)
	if err != nil {
		return nil, err
	}

	if err := create(ctx, conn, v.vppConn); err != nil {
		_, _ = v.Close(ctx, conn, opts...)
		return nil, err
	}
	return conn, nil
}

func (v *l3XConnectClient) Close(ctx context.Context, conn *networkservice.Connection, opts ...grpc.CallOption) (*empty.Empty, error) {
	rv, err := next.Client(ctx).Close(ctx, conn, opts...)
	if err != nil {
		return nil, err
	}
	if err := del(ctx, conn, v.vppConn); err != nil {
		return nil, err
	}

	return rv, nil
}
