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

// GRPCClientParam contains gRPC client parameters. These paramters are
// the same for every gRPC client across the system.
type GRPCClientParam struct {
	ServerEndpoint     string `kong:"help='Server endpoint',default='localhost:10000'"` // Host:port address of the server
	TLS                bool   `kong:"help='Enable TLS',default='false'"`                // TLS enabled
	CAFile             string `kong:"help='CA certificate file',type='existingfile'"`   // CA cert file
	ServerHostOverride string `kong:"help='Host name override for certificate'"`        // Server name returned from the TLS handshake (for debugging)

}
