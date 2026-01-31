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

type LogStreamGenerator struct {
	OktaService
}

func (g *LogStreamGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	logStreams, _, err := client.LogStreamAPI.ListLogStreams(ctx).Execute()
	if err != nil {
		return fmt.Errorf("error listing log streams: %w", err)
	}

	var resources []terraformutils.Resource
	for _, stream := range logStreams {
		var id, name string
		var attributes map[string]interface{} = make(map[string]interface{})

		if stream.LogStreamAws != nil {
			id = stream.LogStreamAws.GetId()
			name = stream.LogStreamAws.GetName()
			// LogStreamSettingsAws is a struct, not a pointer, so we can't check != nil directly if it's not a pointer field.
			// However, if LogStreamAws is present, Settings should be valid struct (might be empty).
			// Let's check if it has fields populated.

			settings := map[string]interface{}{
				"account_id":        stream.LogStreamAws.Settings.GetAccountId(),
				"region":            stream.LogStreamAws.Settings.GetRegion(),
				"event_source_name": stream.LogStreamAws.Settings.GetEventSourceName(),
			}
			attributes["settings"] = []interface{}{settings}

		} else if stream.LogStreamSplunk != nil {
			id = stream.LogStreamSplunk.GetId()
			name = stream.LogStreamSplunk.GetName()

			settings := map[string]interface{}{
				"host":  stream.LogStreamSplunk.Settings.GetHost(),
				"token": stream.LogStreamSplunk.Settings.GetToken(),
			}
			attributes["settings"] = []interface{}{settings}
		} else {
			continue
		}

		resources = append(resources, terraformutils.NewResource(
			id,
			normalizeResourceName(id+"_"+name),
			"okta_log_stream",
			"okta",
			map[string]string{},
			[]string{},
			attributes,
		))
	}
	g.Resources = resources
	return nil
}
