/*
Copyright © 2020 Romber Li <romber2001@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package server

import (
	"context"
	"fmt"
	"github.com/pingcap/errors"
	"net/http"
	"os"
	"time"

	"github.com/romberli/go-template-web/router"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
)

const defaultGracefulShutdownTimeout = 5 * time.Second

type Server interface {
	// Addr returns listen address
	Addr() string
	// PidFile returns pid file path
	PidFile() string
	// Router returns router
	Router() router.Router
	// Run runs server
	Run()
	// Stop stops server
	Stop() error
}

var _ Server = (*server)(nil)

type server struct {
	*http.Server
	addr    string
	pidFile string
	router  router.Router
}

// NewServer returns new *server
func NewServer(addr string, pidFile string, readTimeout, writeTimeout int, router router.Router) Server {
	return &server{
		Server: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
		},
		addr:    addr,
		pidFile: pidFile,
		router:  router,
	}
}

// Addr returns listen address
func (s *server) Addr() string {
	return s.addr
}

// PidFile returns pid file path
func (s *server) PidFile() string {
	return s.pidFile
}

// Router returns router
func (s *server) Router() router.Router {
	return s.router
}

// Run runs server
func (s *server) Run() {
	err := s.router.Run(s.addr)
	if err != nil {
		fmt.Println(fmt.Sprintf("server run failed. addr: %s, pid file: %s\n%+v", s.Addr(), s.PidFile(), err))
		log.Errorf("server run failed. addr: %s, pid file: %s\n%+v", s.Addr(), s.PidFile(), err)
		os.Exit(constant.DefaultAbnormalExitCode)
	}
}

// Stop stops server
func (s *server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultGracefulShutdownTimeout)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
