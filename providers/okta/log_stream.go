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
		if stream.LogStreamAws != nil {
			id = stream.LogStreamAws.GetId()
			name = stream.LogStreamAws.GetName()
		} else if stream.LogStreamSplunk != nil {
			id = stream.LogStreamSplunk.GetId()
			name = stream.LogStreamSplunk.GetName()
		} else {
			continue
		}

		resources = append(resources, terraformutils.NewSimpleResource(
			id,
			normalizeResourceName(id+"_"+name),
			"okta_log_stream",
			"okta",
			[]string{},
		))
	}
	g.Resources = resources
	return nil
}
