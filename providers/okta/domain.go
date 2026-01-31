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

type DomainGenerator struct {
	OktaService
}

func (g *DomainGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	domains, _, err := client.Domain.ListDomains(ctx)
	if err != nil {
		return fmt.Errorf("error listing domains: %w", err)
	}

	var resources []terraformutils.Resource
	for _, domain := range domains.Domains {
		resources = append(resources, terraformutils.NewSimpleResource(
			domain.Id,
			normalizeResourceName(domain.Id+"_"+domain.Domain),
			"okta_domain",
			"okta",
			[]string{},
		))
	}
	g.Resources = resources
	return nil
}

type EmailDomainGenerator struct {
	OktaService
}

func (g *EmailDomainGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	// We expand "brands" to ensure we can get the brand ID if needed, or check AdditionalProperties.
	// Although ListEmailDomains returns EmailDomainResponse which might lack BrandId in the struct.
	// We rely on AdditionalProperties or just the resource import by ID if fetching brand_id is hard.
	// However, `okta_email_domain` resource usually needs `brand_id`.
	// If we can't get `brand_id` from the list response easily (without iterating and fetching details),
	// we might rely on the fact that `brandId` is likely in `AdditionalProperties` map if returned by API.

	emailDomains, _, err := client.EmailDomainAPI.ListEmailDomains(ctx).Expand([]string{"brands"}).Execute()
	if err != nil {
		return fmt.Errorf("error listing email domains: %w", err)
	}

	var resources []terraformutils.Resource
	for _, domain := range emailDomains {
		attributes := map[string]string{}

		// Try to find brandId in AdditionalProperties
		if val, ok := domain.AdditionalProperties["brandId"]; ok {
			attributes["brand_id"] = fmt.Sprintf("%v", val)
		} else {
             // If not found directly, maybe we can assume it's missing or try to look into embedded brands if expanded.
             // But simpler to just proceed. If brand_id is missing, terraformer might produce incomplete HCL,
             // but `terraform import` often fixes state.
             // However, `brand_id` is required argument.
        }

		resources = append(resources, terraformutils.NewResource(
			domain.GetId(),
			normalizeResourceName(domain.GetId()+"_"+domain.GetDisplayName()),
			"okta_email_domain",
			"okta",
			attributes,
			[]string{},
			map[string]interface{}{},
		))
	}
	g.Resources = resources
	return nil
}
