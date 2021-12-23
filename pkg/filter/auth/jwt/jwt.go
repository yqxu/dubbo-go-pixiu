/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jwt

import (
	"encoding/json"
	"fmt"
	stdHttp "net/http"
	"strings"
	"time"
)

import (
	"github.com/MicahParks/keyfunc"
	jwt4 "github.com/golang-jwt/jwt/v4"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind = constant.HTTPAuthJwtFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}

	// Filter is http filter instance
	Filter struct {
		cfg          *Config
		errMsg       []byte
		providerJwks map[string]Provider
	}

	// Config describe the config of Filter
	Config struct {
		ErrMsg    string      `yaml:"err_msg" json:"err_msg" mapstructure:"err_msg"`
		Rules     []Rules     `yaml:"rules" json:"rules" mapstructure:"rules"`
		Providers []Providers `yaml:"providers" json:"providers" mapstructure:"providers"`
	}
)

func (p Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &Config{}, providerJwks: map[string]Provider{}}, nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(ctx *http.HttpContext) {

	path := ctx.Request.RequestURI

	router := false

	for _, rule := range f.cfg.Rules {
		if strings.HasPrefix(path, rule.Match.Prefix) {
			router = true
			if f.validAny(rule, ctx) || f.validAll(rule, ctx) {
				router = false
				break
			}
		}
	}

	if router {
		ctx.WriteJSONWithStatus(stdHttp.StatusUnauthorized, f.errMsg)
		ctx.Abort()
		return
	}

	ctx.Next()

}

func valuePrefix(value, prefix string) string {
	if prefix == "" {
		return value
	}
	return strings.TrimPrefix(value, prefix)
}

// validAny route single provider verification
func (f *Filter) validAny(rule Rules, ctx *http.HttpContext) bool {

	providerName := rule.Requires.RequiresAny.ProviderName

	if provider, ok := f.providerJwks[providerName]; ok {
		ctx.Request.Header.Set(provider.forwardPayloadHeader, provider.issuer)
		if key := ctx.Request.Header.Get(provider.headers.Name); key != "" {
			token, err := jwt4.Parse(valuePrefix(key, provider.headers.ValuePrefix), provider.jwk.Keyfunc)
			if err != nil {
				logger.Warnf("failed to parse JWKs from JSON. provider：%s Error: %s", providerName, err.Error())
				return false
			}
			return token.Valid
		}
	}

	return false
}

// validAll route multiple provider verification
func (f *Filter) validAll(rule Rules, ctx *http.HttpContext) bool {

	for _, requirement := range rule.Requires.RequiresAll {
		if provider, ok := f.providerJwks[requirement.ProviderName]; ok {
			ctx.Request.Header.Set(provider.forwardPayloadHeader, provider.issuer)
			if key := ctx.Request.Header.Get(provider.headers.Name); key != "" {
				token, err := jwt4.Parse(valuePrefix(key, provider.headers.ValuePrefix), provider.jwk.Keyfunc)
				if err != nil {
					logger.Warnf("failed to parse JWKs from JSON. provider：%s Error: %s", requirement.ProviderName, err.Error())
					continue
				}

				if token.Valid {
					return true
				}
			}
		}
	}

	return false
}

func (f *Filter) Apply() error {

	if f.cfg.ErrMsg == "" {
		f.cfg.ErrMsg = "token invalid"
	}

	errMsg, _ := json.Marshal(http.ErrResponse{Message: f.cfg.ErrMsg})
	f.errMsg = errMsg

	if len(f.cfg.Providers) == 0 {
		return fmt.Errorf("providers is null")
	}

	for _, provider := range f.cfg.Providers {

		if provider.Local != nil {
			jwksJSON := json.RawMessage(provider.Local.InlineString)
			jwks, err := keyfunc.NewJSON(jwksJSON)
			if err != nil {
				logger.Warnf("failed to create JWKs from JSON. provider：%s Error: %s", provider.Name, err.Error())
			} else {
				provider.FromHeaders.setDefault()
				f.providerJwks[provider.Name] = Provider{jwk: jwks, headers: provider.FromHeaders,
					issuer: provider.Issuer, forwardPayloadHeader: provider.ForwardPayloadHeader}
				continue
			}
		}

		if provider.Remote != nil {
			uri := provider.Remote.HttpURI
			timeout, err := time.ParseDuration(uri.TimeOut)
			if err != nil {
				logger.Warnf("jwt provides timeout parse fail: %s", err.Error())
				continue
			}

			options := keyfunc.Options{RefreshTimeout: timeout}
			jwks, err := keyfunc.Get(uri.Uri, options)
			if err != nil {
				logger.Warnf("failed to create JWKs from resource at the given URL. provider：%s Error: %s", provider.Name, err.Error())
			} else {
				provider.FromHeaders.setDefault()
				f.providerJwks[provider.Name] = Provider{jwk: jwks, headers: provider.FromHeaders,
					issuer: provider.Issuer, forwardPayloadHeader: provider.ForwardPayloadHeader}
			}
		}
	}

	if len(f.providerJwks) == 0 {
		return fmt.Errorf("providers is null")
	}

	return nil
}

func (h *FromHeaders) setDefault() {
	if h.Name == "" {
		h.Name = "Authorization"
	}

	if h.ValuePrefix == "" {
		h.ValuePrefix = "Bearer "
	}
}

func (f *Filter) Config() interface{} {
	return f.cfg
}