/*
Copyright 2021 The Dapr Authors
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

package config

import (
	"flag"
	"strings"

	"github.com/dapr/kit/logger"

	"github.com/liuxd6825/dapr/pkg/metrics"
	"github.com/liuxd6825/dapr/pkg/placement/raft"
)

const (
	defaultCredentialsPath   = "/var/run/dapr/credentials"
	defaultHealthzPort       = 8080
	defaultPlacementPort     = 50005
	defaultReplicationFactor = 100
)

type config struct {
	// Raft protocol configurations
	RaftID           string
	RaftPeerString   string
	RaftPeers        []raft.PeerInfo
	RaftInMemEnabled bool
	RaftLogStorePath string

	// Placement server configurations
	PlacementPort int
	HealthzPort   int
	CertChainPath string
	TlsEnabled    bool

	ReplicationFactor int

	// Log and metrics configurations
	LoggerOptions   logger.Options
	MetricsExporter metrics.Exporter
}

func NewConfig() *config {
	// Default configuration
	cfg := config{
		RaftID:           "dapr-placement-0",
		RaftPeerString:   "dapr-placement-0=127.0.0.1:8201",
		RaftPeers:        []raft.PeerInfo{},
		RaftInMemEnabled: true,
		RaftLogStorePath: "",

		PlacementPort: defaultPlacementPort,
		HealthzPort:   defaultHealthzPort,
		CertChainPath: defaultCredentialsPath,
		TlsEnabled:    false,
	}

	flag.StringVar(&cfg.RaftID, "id", cfg.RaftID, "Placement server ID.")
	flag.StringVar(&cfg.RaftPeerString, "initial-cluster", cfg.RaftPeerString, "raft cluster peers")
	flag.BoolVar(&cfg.RaftInMemEnabled, "inmem-store-enabled", cfg.RaftInMemEnabled, "Enable in-memory log and snapshot store unless --raft-logstore-path is set")
	flag.StringVar(&cfg.RaftLogStorePath, "raft-logstore-path", cfg.RaftLogStorePath, "raft log store path.")
	flag.IntVar(&cfg.PlacementPort, "port", cfg.PlacementPort, "sets the gRPC port for the placement service")
	flag.IntVar(&cfg.HealthzPort, "healthz-port", cfg.HealthzPort, "sets the HTTP port for the healthz server")
	flag.StringVar(&cfg.CertChainPath, "certchain", cfg.CertChainPath, "Path to the credentials directory holding the cert chain")
	flag.BoolVar(&cfg.TlsEnabled, "tls-enabled", cfg.TlsEnabled, "Should TLS be enabled for the placement gRPC server")
	flag.IntVar(&cfg.ReplicationFactor, "ReplicationFactor", defaultReplicationFactor, "sets the replication factor for actor distribution on vnodes")

	cfg.LoggerOptions = logger.DefaultOptions()
	cfg.LoggerOptions.AttachCmdFlags(flag.StringVar, flag.BoolVar)

	cfg.MetricsExporter = metrics.NewExporter(metrics.DefaultMetricNamespace)
	cfg.MetricsExporter.Options().AttachCmdFlags(flag.StringVar, flag.BoolVar)

	flag.Parse()

	cfg.RaftPeers = parsePeersFromFlag(cfg.RaftPeerString)
	if cfg.RaftLogStorePath != "" {
		cfg.RaftInMemEnabled = false
	}

	return &cfg
}

func parsePeersFromFlag(val string) []raft.PeerInfo {
	peers := []raft.PeerInfo{}

	p := strings.Split(val, ",")
	for _, addr := range p {
		peer := strings.Split(addr, "=")
		if len(peer) != 2 {
			continue
		}

		peers = append(peers, raft.PeerInfo{
			ID:      strings.TrimSpace(peer[0]),
			Address: strings.TrimSpace(peer[1]),
		})
	}

	return peers
}
