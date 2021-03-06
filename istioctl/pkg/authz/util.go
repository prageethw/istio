// Copyright 2019 Istio Authors
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

package authz

import (
	"strings"

	"istio.io/istio/pilot/pkg/model"

	envoy_auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
)

// PolicyTypeToConfigs maps policy type (e.g. service-role) to a list of its config.
type PolicyTypeToConfigs map[string][]model.Config

func getCertificate(ctx *envoy_auth.CommonTlsContext) string {
	cert := "none"
	if ctx == nil {
		return cert
	}

	staticConfigs := make([]string, 0)
	for _, cert := range ctx.GetTlsCertificates() {
		name := cert.GetCertificateChain().GetFilename()
		if name == "" {
			name = "<inline>"
		}
		staticConfigs = append(staticConfigs, name)
	}

	sdsConfigs := make([]string, 0)
	for _, sds := range ctx.GetTlsCertificateSdsSecretConfigs() {
		sdsConfigs = append(sdsConfigs, sds.Name)
	}

	if len(staticConfigs) != 0 {
		cert = strings.Join(staticConfigs, ",")
	}

	if len(sdsConfigs) != 0 {
		if len(staticConfigs) != 0 {
			cert += "; "
		}
		cert += "SDS: " + strings.Join(sdsConfigs, ",")
	}
	return cert
}
