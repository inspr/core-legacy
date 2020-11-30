/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	pb "gitlab.inspr.com/ptcar/core/pkg/operator/comunication"
	"google.golang.org/grpc"
)

const (
	port = ":49051"
)

var even int

// server is used to implement operator.OperatorServer.
type server struct {
	pb.UnimplementedOperatorServer
}

// ApplyOperation implements operator.OperatorServer
func (s *server) ApplyOperation(ctx context.Context, in *pb.OperationRequest) (*pb.OperationReply, error) {
	log.Printf("Received: %v -> %v", in.GetKind(), in.GetValue())
	even = even + 1
	if even%2 == 1 {
		return &pb.OperationReply{
			Err:    "OK! OP (" + in.GetKind() + ") - " + in.GetValue() + "\nGo ahead!",
			Status: true,
		}, nil
	} else {
		return &pb.OperationReply{
			Err:    "Stop! Error to parser OP (" + in.GetKind() + ") \nPermission denied!",
			Status: false,
		}, nil
	}
}

func main() {
	even = 0
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOperatorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
