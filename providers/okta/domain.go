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

	// Fetch brands to map EmailDomainId -> BrandId
	brands, _, err := client.BrandsAPI.ListBrands(ctx).Execute()
	if err != nil {
		return fmt.Errorf("error listing brands for email domains: %w", err)
	}

	domainToBrand := make(map[string]string)
	for _, brand := range brands {
		if brand.HasEmailDomainId() {
			domainToBrand[brand.GetEmailDomainId()] = brand.GetId()
		}
	}

	// Try to get expanded brands as well
	emailDomains, _, err := client.EmailDomainAPI.ListEmailDomains(ctx).Expand([]string{"brands"}).Execute()
	if err != nil {
		return fmt.Errorf("error listing email domains: %w", err)
	}

	var resources []terraformutils.Resource
	for _, domain := range emailDomains {
		attributes := map[string]string{}

		fmt.Printf("DEBUG: Processing domain %s (%s)\n", domain.GetId(), domain.GetDisplayName())

		brandId := ""
		if bid, ok := domainToBrand[domain.GetId()]; ok {
			brandId = bid
			fmt.Printf("DEBUG: Found brandId %s in domainToBrand map\n", brandId)
		}

		if brandId == "" && domain.HasEmbedded() {
			embedded := domain.GetEmbedded()
			if embedded.Brands != nil && len(embedded.Brands) > 0 {
				brandId = embedded.Brands[0].GetId()
				fmt.Printf("DEBUG: Found brandId %s in Embedded.Brands\n", brandId)
			} else {
				fmt.Println("DEBUG: Embedded.Brands is empty or nil")
			}
		} else if brandId == "" {
			fmt.Println("DEBUG: No Embedded data")
		}

		// Fallback: check AdditionalProperties for "brandId"
		if brandId == "" && domain.AdditionalProperties != nil {
			if val, ok := domain.AdditionalProperties["brandId"]; ok {
				brandId = fmt.Sprintf("%v", val)
				fmt.Printf("DEBUG: Found brandId %s in AdditionalProperties\n", brandId)
			}
		}

		if brandId != "" {
			attributes["brand_id"] = brandId
		} else {
			fmt.Printf("DEBUG: Failed to find brandId for domain %s\n", domain.GetId())
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
