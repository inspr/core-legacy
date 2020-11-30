// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

syntax = "proto3";

option go_package = "gitlab.inspr.com/ptcar/core/pkg/meta";
option java_multiple_files = true;
option java_package = "io.grpc.insprd.meta";
option java_outer_classname = "InsprDMetaProto";

package meta;

//import "k8s.io/kubernetes/pkg/api/v1/generated.proto";
//import "k8s.io/api/core/v1/generated.proto";
option go_package = "k8s.io/api/core/v1";

// A stub message to search data
message stub {}

// Annotation is a method to describe annotations
message Annotation {
    string name = 1;
    string value = 2;
}

// Metadata describes a metadata that should be used by inspr
message Metadata {
    string name = 1;
    string reference = 2;
    repeated Annotation annotations = 3;
    string parent = 4;
    string sha256 = 5;
}

// NodeSpec represents a configuration for a node. The image represents the Docker image for the main container of the Node.
// If the node has an specific Kubernetes configuration, the configuration can be injected via the Kubernetes field. When
// Kubernetes is set, the Image field gets igored.
message NodeSpec {
    string image = 1;
    PodSpec kubernetes = 2;
}

// Node represents an inspr component that is a node.
message Node {
    Metadata metadata = 1;
    PodSpec spec = 2;
}

// AppBoundary represents the connections this app can make to other apps. These are the fields that can be overriten
// by the ChannelAliases when instantiating the app.
message AppBoundary {
    repeated string input = 1;
    repeated string ouput = 2;
}

// AppSpec represents the configuration of an App.
//
// The app contains a list of apps and a list of nodes. The apps and nodes can be dereferenced by it's metadata
// reference, at CLI time.
//
// The boundary represent the possible connections to other apps, and the fields that can be overriten when instantiating the app.
message AppSpec {
    Node node = 1;
    repeated App apps = 2;
    repeated Channel channel = 3;
    AppBoundary boundary = 4;
}

// App is an inspr component that represents an App. An App can contain other apps, channels and other components.
message App {
    Metadata metadata = 1;
    AppSpec spec = 2;
}

// ChannelType is the type of the channel. It can be a reference to an outsourced type or can be a local type. This local
// type will be defined via the workspace and instantiated as a []byte on the cluster
message ChannelType {
    string reference = 1;
    string fileReference = 2;
    bytes schema = 3;
}

// ChannelSpec is the specification of a channel. (the external variable is just an idea)
message ChannelSpec {
    ChannelType type = 1;
    bool external = 2;
}

// Channel is an Inspr component that represents a Channel. The channel can be instantiated by using a reference to either
// a local file or an URL to an uploaded file.
message Channel {
    Metadata metadata = 1;
    ChannelSpec spec = 2;
}
