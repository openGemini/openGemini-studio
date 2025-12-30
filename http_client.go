// Copyright 2025 openGemini Authors
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
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/openGemini/opengemini-client-go/opengemini"
)

var _ HttpClient = (*HttpClientCreator)(nil)

type HttpClient interface {
	SetDebug(debug bool)
	SetAuth(username, password string)
	Ping() error
	Query(context.Context, *opengemini.Query) (*opengemini.QueryResult, error)
	Write(ctx context.Context, database, retentionPolicy, raw, precision string) error
	Databases(ctx context.Context) ([]string, error)
	RetentionPolicies(ctx context.Context, database string) ([]*RetentionPolicy, error)
	Measurements(ctx context.Context, database string) ([]string, error)
}

type HttpClientCreator struct {
	HostPort string
	client   *http.Client
	basic    string
	debug    bool
}

func (h *HttpClientCreator) RetentionPolicies(ctx context.Context, database string) ([]*RetentionPolicy, error) {
	response, err := h.Query(ctx, &opengemini.Query{
		Database: database,
		Command:  "SHOW RETENTION POLICIES",
	})
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, fmt.Errorf("show retention policies failed: %s", response.Error)
	}

	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return nil, nil
	}
	var (
		seriesValues    = response.Results[0].Series[0].Values
		column          = response.Results[0].Series[0].Columns
		retentionPolicy = make([]*RetentionPolicy, 0, len(seriesValues))
	)

	for idx, v := range seriesValues {
		columnName := column[idx]
		var name, duration string
		if columnName == "name" {
			name = v[0].(string)
		}
		if columnName == "duration" {
			duration = v[0].(string)
		}
		retentionPolicy = append(retentionPolicy, &RetentionPolicy{
			Name:     name,
			Duration: duration,
		})
	}
	return retentionPolicy, nil
}

func (h *HttpClientCreator) Measurements(ctx context.Context, database string) ([]string, error) {
	response, err := h.Query(ctx, &opengemini.Query{
		Database: database,
		Command:  "SHOW MEASUREMENTS",
	})
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, fmt.Errorf("show measurements failed: %s", response.Error)
	}

	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return nil, nil
	}
	var (
		seriesValues = response.Results[0].Series[0].Values
		measurements = make([]string, 0, len(seriesValues))
	)

	for _, v := range seriesValues {
		measurements = append(measurements, v[0].(string))
	}
	return measurements, nil
}

func (h *HttpClientCreator) Databases(ctx context.Context) ([]string, error) {
	response, err := h.Query(ctx, &opengemini.Query{
		Command: "SHOW DATABASES",
	})
	if err != nil {
		return nil, err
	}
	if len(response.Error) > 0 {
		return nil, fmt.Errorf("show datababse failed: %s", response.Error)
	}
	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return []string{}, nil
	}
	var (
		values   = response.Results[0].Series[0].Values
		dbResult = make([]string, 0, len(values))
	)

	for _, v := range values {
		if len(v) == 0 {
			continue
		}
		val, ok := v[0].(string)
		if !ok {
			continue
		}
		dbResult = append(dbResult, val)
	}
	return dbResult, nil
}

func (h *HttpClientCreator) SetAuth(username, password string) {
	h.basic = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

func (h *HttpClientCreator) SetDebug(debug bool) {
	h.debug = debug
}

func NewHttpClient(cfg *ConnectConfig) (HttpClient, error) {
	var client = &HttpClientCreator{client: &http.Client{
		Timeout: 600 * time.Second,
	}}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	var schema = strings.ToLower(cfg.HTTPSchema)

	if schema == "https" {
		certificateManager, err := NewCertificateManager(cfg.CACertificate, cfg.ClientCertificate, cfg.ClientKey)
		if err != nil {
			return nil, errors.New("cannot load certificate: " + err.Error())
		}
		transport.TLSClientConfig = &tls.Config{
			RootCAs:      certificateManager.CAPool,
			Certificates: []tls.Certificate{certificateManager.Certificate},
		}

		if cfg.InsecureTls {
			transport.TLSClientConfig.InsecureSkipVerify = true
		}

		if cfg.InsecureHostname {
			transport.TLSClientConfig.InsecureSkipVerify = true
			transport.TLSClientConfig.VerifyPeerCertificate = func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
				if len(rawCerts) == 0 {
					return errors.New("no certificates provided by server")
				}
				cert, err := x509.ParseCertificate(rawCerts[0])
				if err != nil {
					return fmt.Errorf("failed to parse server certificate: %w", err)
				}
				now := time.Now()
				if now.After(cert.NotAfter) {
					return errors.New("server certificate has expired")
				}
				opts := x509.VerifyOptions{
					DNSName:       "",
					Roots:         certificateManager.CAPool,
					Intermediates: x509.NewCertPool(),
				}
				for _, rawCert := range rawCerts[1:] {
					intermediateCert, err := x509.ParseCertificate(rawCert)
					if err != nil {
						return fmt.Errorf("failed to parse intermediate certificate: %w", err)
					}
					opts.Intermediates.AddCert(intermediateCert)
				}
				if _, err := cert.Verify(opts); err != nil {
					return fmt.Errorf("server certificate chain validation failed: %w", err)
				}
				return nil
			}
		}
	}

	if cfg.EnableAuth && cfg.Username != "" && cfg.Password != "" {
		client.SetAuth(cfg.Username, cfg.Password)
	}

	client.HostPort = schema + "://" + cfg.Address

	client.client.Transport = transport
	return client, nil
}

// Ping will check to see if the server is up
func (h *HttpClientCreator) Ping() error {
	var urlPath = h.HostPort + "/ping"
	response, err := h.innerRequest(context.Background(), http.MethodGet, urlPath, nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 204 {
		return fmt.Errorf("ping failed: %s", response.Status)
	}
	return nil
}

func (h *HttpClientCreator) Query(ctx context.Context, query *opengemini.Query) (*opengemini.QueryResult, error) {
	urlPath := h.HostPort + "/query"

	var queryValues = make(url.Values)
	queryValues.Add("db", query.Database)
	queryValues.Add("rp", query.RetentionPolicy)
	queryValues.Add("q", query.Command)
	queryValues.Add("epoch", query.Precision.Epoch())

	response, err := h.innerRequest(ctx, http.MethodPost, urlPath, strings.NewReader(queryValues.Encode()))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("response status_code: " + response.Status + ", body: " + string(data))
	}
	var qr = new(opengemini.QueryResult)
	err = json.Unmarshal(data, qr)
	if err != nil {
		return nil, err
	}
	return qr, nil
}

func (h *HttpClientCreator) Write(ctx context.Context, database, retentionPolicy, raw, precision string) error {
	urlPath := h.HostPort + "/write"
	u, err := url.Parse(urlPath)
	if err != nil {
		return err
	}

	var writeValues = make(url.Values)
	writeValues.Add("db", database)
	writeValues.Add("rp", retentionPolicy)
	writeValues.Add("precision", precision)
	u.RawQuery = writeValues.Encode()

	response, err := h.innerRequest(ctx, http.MethodPost, u.String(), strings.NewReader(raw))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return errors.New("write failed: " + response.Status)
	}
	return nil
}

func (h *HttpClientCreator) innerRequest(ctx context.Context, method, urlPath string, reader io.Reader) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, urlPath, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "opengemini-cli")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if h.basic != "" {
		request.Header.Set("Authorization", "Basic "+h.basic)
	}

	if h.debug {
		dumpRequest, _ := httputil.DumpRequest(request, true)
		fmt.Printf("---------- REQUEST DEBUG ----------\n%s\n---------- REQUEST DEBUG ----------\n", string(dumpRequest))
	}

	response, err := h.client.Do(request)

	if h.debug {
		dumpResponse, _ := httputil.DumpResponse(response, true)
		fmt.Printf("---------- RESPONSE DEBUG ----------\n%s\n---------- RESPONSE DEBUG ----------\n", string(dumpResponse))
	}

	return response, err
}

type CertificateManager struct {
	CAContent   []byte
	CAPool      *x509.CertPool
	Certificate tls.Certificate
}

func NewCertificateManager(ca, certificate, certificateKey string) (*CertificateManager, error) {
	var cm = new(CertificateManager)
	if ca != "" {
		content, err := os.ReadFile(ca)
		if err != nil {
			return nil, err
		}
		cm.CAPool = x509.NewCertPool()
		if !cm.CAPool.AppendCertsFromPEM(content) {
			return nil, errors.New("failed to parse ca certificate")
		}
	}
	if certificate != "" && certificateKey != "" {
		keyPair, err := tls.LoadX509KeyPair(certificate, certificateKey)
		if err != nil {
			return nil, err
		}
		cm.Certificate = keyPair
	}

	return cm, nil
}

func (cm *CertificateManager) CreateTls(insecureTls, insecureHostname bool) *tls.Config {
	var cfg = &tls.Config{
		RootCAs:      cm.CAPool,
		Certificates: []tls.Certificate{cm.Certificate},
	}
	if insecureTls {
		cfg.InsecureSkipVerify = true
	}

	if insecureHostname {
		cfg.InsecureSkipVerify = true
		cfg.VerifyPeerCertificate = func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
			if len(rawCerts) == 0 {
				return errors.New("no certificates provided by server")
			}
			cert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return fmt.Errorf("failed to parse server certificate: %w", err)
			}
			now := time.Now()
			if now.After(cert.NotAfter) {
				return errors.New("server certificate has expired")
			}
			opts := x509.VerifyOptions{
				DNSName:       "",
				Roots:         cm.CAPool,
				Intermediates: x509.NewCertPool(),
			}
			for _, rawCert := range rawCerts[1:] {
				intermediateCert, err := x509.ParseCertificate(rawCert)
				if err != nil {
					return fmt.Errorf("failed to parse intermediate certificate: %w", err)
				}
				opts.Intermediates.AddCert(intermediateCert)
			}
			if _, err := cert.Verify(opts); err != nil {
				return fmt.Errorf("server certificate chain validation failed: %w", err)
			}
			return nil
		}
	}
	return cfg
}
