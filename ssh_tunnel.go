// Copyright 2026 openGemini Authors
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

package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

// SSHTunnel manages SSH tunnel connections
type SSHTunnel struct {
	sshClient  *ssh.Client
	listener   net.Listener
	localAddr  string
	remoteAddr string
	sshAddr    string
	config     *ssh.ClientConfig
	wg         sync.WaitGroup
	stopChan   chan struct{}
	started    bool
	mu         sync.Mutex
}

// NewSSHTunnel creates a new SSH tunnel
func NewSSHTunnel(cfg *ConnectConfig) (*SSHTunnel, error) {
	if !cfg.EnableSSH {
		return nil, errors.New("SSH is not enabled")
	}

	if cfg.SSHHost == "" || cfg.SSHPort == 0 {
		return nil, errors.New("SSH host and port are required")
	}

	if cfg.SSHUsername == "" {
		return nil, errors.New("SSH username is required")
	}

	// Create SSH client config
	sshConfig := &ssh.ClientConfig{
		User:            cfg.SSHUsername,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, you should use proper host key verification
		Timeout:         0,
	}

	// Add authentication method
	if cfg.SSHKeyPath != "" {
		// Use private key authentication
		key, err := os.ReadFile(cfg.SSHKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read SSH private key: %w", err)
		}

		var signer ssh.Signer
		if cfg.SSHKeyPassphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(cfg.SSHKeyPassphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(key)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse SSH private key: %w", err)
		}

		sshConfig.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	} else if cfg.SSHPassword != "" {
		// Use password authentication
		sshConfig.Auth = []ssh.AuthMethod{
			ssh.Password(cfg.SSHPassword),
		}
	} else {
		return nil, errors.New("SSH authentication method not configured (password or key required)")
	}

	tunnel := &SSHTunnel{
		sshAddr:    fmt.Sprintf("%s:%d", cfg.SSHHost, cfg.SSHPort),
		remoteAddr: cfg.Address,
		config:     sshConfig,
		stopChan:   make(chan struct{}),
	}

	return tunnel, nil
}

// Start starts the SSH tunnel
func (t *SSHTunnel) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.started {
		return errors.New("tunnel already started")
	}

	// Connect to SSH server
	sshClient, err := ssh.Dial("tcp", t.sshAddr, t.config)
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	t.sshClient = sshClient

	// Create local listener on random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		sshClient.Close()
		return fmt.Errorf("failed to create local listener: %w", err)
	}
	t.listener = listener
	t.localAddr = listener.Addr().String()

	t.started = true

	// Start accepting connections
	t.wg.Add(1)
	go t.acceptConnections()

	return nil
}

// acceptConnections accepts incoming connections and forwards them through SSH tunnel
func (t *SSHTunnel) acceptConnections() {
	defer t.wg.Done()

	for {
		select {
		case <-t.stopChan:
			return
		default:
		}

		localConn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.stopChan:
				return
			default:
				continue
			}
		}

		t.wg.Add(1)
		go t.handleConnection(localConn)
	}
}

// handleConnection handles a single connection through the SSH tunnel
func (t *SSHTunnel) handleConnection(localConn net.Conn) {
	defer t.wg.Done()
	defer localConn.Close()

	// Connect to remote server through SSH tunnel
	remoteConn, err := t.sshClient.Dial("tcp", t.remoteAddr)
	if err != nil {
		return
	}
	defer remoteConn.Close()

	// Forward data bidirectionally
	done := make(chan struct{}, 2)

	// Copy from local to remote
	go func() {
		io.Copy(remoteConn, localConn)
		done <- struct{}{}
	}()

	// Copy from remote to local
	go func() {
		io.Copy(localConn, remoteConn)
		done <- struct{}{}
	}()

	// Wait for one direction to finish
	<-done
}

// Stop stops the SSH tunnel
func (t *SSHTunnel) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.started {
		return nil
	}

	// Signal to stop accepting new connections
	close(t.stopChan)

	// Close the listener
	if t.listener != nil {
		t.listener.Close()
	}

	// Wait for all connections to finish
	t.wg.Wait()

	// Close SSH client
	if t.sshClient != nil {
		t.sshClient.Close()
	}

	t.started = false
	return nil
}

// LocalAddr returns the local address of the tunnel
func (t *SSHTunnel) LocalAddr() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.localAddr
}

// IsStarted returns whether the tunnel is started
func (t *SSHTunnel) IsStarted() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.started
}
