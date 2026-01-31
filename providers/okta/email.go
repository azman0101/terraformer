// Copyright 2024 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package okta

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

type EmailSenderGenerator struct {
	OktaService
}

func (g *EmailSenderGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	servers, _, err := client.EmailServerAPI.ListEmailServers(ctx).Execute()
	if err != nil {
		return fmt.Errorf("error listing email servers: %w", err)
	}

	var resources []terraformutils.Resource
	for _, server := range servers.GetEmailServers() {
		name := server.GetAlias()
		if name == "" {
			name = server.GetHost()
		}
		resources = append(resources, terraformutils.NewSimpleResource(
			server.GetId(),
			normalizeResourceName(name),
			"okta_email_sender",
			"okta",
			[]string{},
		))
	}
	g.Resources = resources
	return nil
}
