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
	"github.com/okta/okta-sdk-golang/v5/okta"
)

type BrandGenerator struct {
	OktaService
}

func (g *BrandGenerator) createResources(brands []okta.BrandWithEmbedded) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, brand := range brands {
		resources = append(resources, terraformutils.NewSimpleResource(
			brand.GetId(),
			normalizeResourceName(brand.GetId()+"_"+brand.GetName()),
			"okta_brand",
			"okta",
			[]string{},
		))
	}
	return resources
}

func (g *BrandGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	brands, _, err := client.BrandsAPI.ListBrands(ctx).Execute()
	if err != nil {
		return fmt.Errorf("error listing brands: %w", err)
	}

	g.Resources = g.createResources(brands)
	return nil
}

type ThemeGenerator struct {
	OktaService
}

func (g *ThemeGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	brands, _, err := client.BrandsAPI.ListBrands(ctx).Execute()
	if err != nil {
		return fmt.Errorf("error listing brands for themes: %w", err)
	}

	var resources []terraformutils.Resource
	for _, brand := range brands {
		themes, _, err := client.ThemesAPI.ListBrandThemes(ctx, brand.GetId()).Execute()
		if err != nil {
			return fmt.Errorf("error listing themes for brand %s: %w", brand.GetId(), err)
		}
		for _, theme := range themes {
			attributes := map[string]string{
				"brand_id": brand.GetId(),
			}

			if theme.HasPrimaryColorHex() {
				attributes["primary_color_hex"] = theme.GetPrimaryColorHex()
			}
			if theme.HasSecondaryColorHex() {
				attributes["secondary_color_hex"] = theme.GetSecondaryColorHex()
			}
			if theme.HasSignInPageTouchPointVariant() {
				attributes["sign_in_page_touch_point_variant"] = string(theme.GetSignInPageTouchPointVariant())
			}
			if theme.HasEndUserDashboardTouchPointVariant() {
				attributes["end_user_dashboard_touch_point_variant"] = string(theme.GetEndUserDashboardTouchPointVariant())
			}
			if theme.HasErrorPageTouchPointVariant() {
				attributes["error_page_touch_point_variant"] = string(theme.GetErrorPageTouchPointVariant())
			}
			if theme.HasEmailTemplateTouchPointVariant() {
				attributes["email_template_touch_point_variant"] = string(theme.GetEmailTemplateTouchPointVariant())
			}

			resources = append(resources, terraformutils.NewResource(
				brand.GetId()+"/"+theme.GetId(),
				normalizeResourceName(brand.GetId()+"_"+theme.GetId()),
				"okta_theme",
				"okta",
				attributes,
				[]string{},
				map[string]interface{}{},
			))
		}
	}
	g.Resources = resources
	return nil
}
