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

type ProfileMappingGenerator struct {
	OktaService
}

func (g *ProfileMappingGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	mappingsList, _, err := client.ProfileMappingAPI.ListProfileMappings(ctx).Execute()
	if err != nil {
		return fmt.Errorf("error listing profile mappings: %w", err)
	}

	var resources []terraformutils.Resource
	for _, mappingSummary := range mappingsList {
		// Fetch full details for properties
		mapping, _, err := client.ProfileMappingAPI.GetProfileMapping(ctx, mappingSummary.GetId()).Execute()
		if err != nil {
			// If fetching fails, maybe log error but continue?
			// For now, let's propagate error as it's critical for functionality
			return fmt.Errorf("error getting profile mapping %s: %w", mappingSummary.GetId(), err)
		}

		name := mapping.GetId()
		// Try to construct a better name
		source := mapping.GetSource()
		target := mapping.GetTarget()
		sourceName := source.GetName()
		targetName := target.GetName()
		if sourceName != "" && targetName != "" {
			name = sourceName + "_to_" + targetName
		}

		attributes := map[string]interface{}{}
		if mapping.Properties != nil {
			var mappingList []interface{}
			for key, prop := range *mapping.Properties {
				m := map[string]interface{}{
					"id":         key,
					"expression": prop.GetExpression(),
					"push_status": prop.GetPushStatus(),
				}
				mappingList = append(mappingList, m)
			}
			attributes["mappings"] = mappingList
		}

		resources = append(resources, terraformutils.NewResource(
			mapping.GetId(),
			normalizeResourceName(name),
			"okta_profile_mapping",
			"okta",
			map[string]string{},
			[]string{},
			attributes,
		))
	}
	g.Resources = resources
	return nil
}
