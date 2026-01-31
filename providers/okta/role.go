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

type RoleGenerator struct {
	OktaService
}

func (g *RoleGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	roles, _, err := client.RoleAPI.ListRoles(ctx).Execute()
	if err != nil {
		return fmt.Errorf("error listing roles: %w", err)
	}

	var resources []terraformutils.Resource
	for _, role := range roles.GetRoles() {
		// We import all roles as okta_admin_role_custom.
		// Users might need to filter manually if they only want custom roles.
		resources = append(resources, terraformutils.NewSimpleResource(
			role.GetId(),
			normalizeResourceName(role.GetId()+"_"+role.GetLabel()),
			"okta_admin_role_custom",
			"okta",
			[]string{},
		))
	}
	g.Resources = resources
	return nil
}
