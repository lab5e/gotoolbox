package metrics

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
import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server is the monitoring endpoint. The monitoring endpoint provides
// counters for the service and a resource to do live traces of a running
// service. Overall performance is affected by the trace so use with caution
// on running systems under load.
type Server struct {
	Listener     net.Listener
	mux          *http.ServeMux
	srv          *http.Server
	healthStatus *int32
}

// NewMonitoringServer creates a new monitoring endpoint
func NewMonitoringServer(endpoint string) (*Server, error) {
	ret := &Server{
		healthStatus: new(int32),
	}
	ret.SetStatus(http.StatusServiceUnavailable)
	var err error
	ret.Listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	ret.mux = http.NewServeMux()
	ret.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<script language="JavaScript">
					function startTrace() {
						var xhttp = new XMLHttpRequest();
						xhttp.open('POST', '/trace', true);
						xhttp.send('2');
						alert('Trace started');
					}
				</script>
				<ul>
					<li><a href="/pprof/">Profiling</a></li>
					<li><a href="/metrics">Metrics</a></li>
					<li><button onClick="startTrace()">Trace for 2s</button></li>
				</ul>
			</html>
		`))
	})
	ret.mux.Handle("/metrics", promhttp.Handler())
	ret.mux.HandleFunc("/pprof/", pprof.Index)
	ret.mux.HandleFunc("/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	ret.mux.HandleFunc("/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	ret.mux.HandleFunc("/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	ret.mux.HandleFunc("/pprof/block", pprof.Handler("block").ServeHTTP)
	ret.mux.HandleFunc("/pprof/profile", pprof.Handler("profile").ServeHTTP)
	ret.mux.HandleFunc("/pprof/heap", pprof.Handler("heap").ServeHTTP)
	enableTracingRoutine()
	ret.mux.HandleFunc("/trace", traceHandler())
	ret.mux.HandleFunc("/healthz", ret.healthzHandler)
	ret.srv = &http.Server{}
	return ret, nil
}

// Start launches the server
func (s *Server) Start() error {
	go func() {
		if err := http.Serve(s.Listener, s.mux); err != http.ErrServerClosed {
			log.Printf("Unable to listen and serve: %v", err)
		}
	}()
	return nil
}

// ServerURL is the URL for the server
func (s *Server) ServerURL() string {
	return fmt.Sprintf("http://%s", s.Listener.Addr().String())
}

// Shutdown stops the server. There is a 2 second timeout.
func (s *Server) Shutdown() error {
	s.Listener.Close()
	return nil
}

// ListenAddress returns the external reachable address for the monitoring service.
// If it is listening on the loopback adapter the loopback address is returned.
func (s *Server) ListenAddress() net.Addr {
	if s.Listener == nil {
		log.Printf("ListenAddress is nil")
	}
	return s.Listener.Addr()
}

// SetStatus sets the health status that is reported, ie http.StatusOK or
// http.StatusServiceUnavailable
func (s *Server) SetStatus(httpStatus int) {
	atomic.StoreInt32(s.healthStatus, int32(httpStatus))
}

// healthzHandler responds to health requests. When the node is available it
// returns 200, 503 otherwise.
func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(atomic.LoadInt32(s.healthStatus)))
}
