package grpcutil

//
//Copyright 2019 Telenor Digital AS
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//

// GRPCServerParam holds parameters for a GRPC server
type GRPCServerParam struct {
	Endpoint string `kong:"help='Service endpoint',default='localhost:0'"`
	TLS      bool   `kong:"help='Enable TLS',default='false'"`
	CertFile string `kong:"help='Certificate file',type='existingfile'"`
	KeyFile  string `kong:"help='Certificate key file',type='existingfile'"`
	Metrics  bool   `kong:"help='Add Prometheus interceptors for server',default='true'"`
}
