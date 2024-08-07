// Copyright (c) 2023 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adjuster

import (
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/otelsemconv"
)

var otelLibraryKeys = map[string]struct{}{
	string(otelsemconv.TelemetrySDKLanguageKey):   {},
	string(otelsemconv.TelemetrySDKNameKey):       {},
	string(otelsemconv.TelemetrySDKVersionKey):    {},
	string(otelsemconv.TelemetryDistroNameKey):    {},
	string(otelsemconv.TelemetryDistroVersionKey): {},
}

func OTelTagAdjuster() Adjuster {
	adjustSpanTags := func(span *model.Span) {
		newI := 0
		for i, tag := range span.Tags {
			if _, ok := otelLibraryKeys[tag.Key]; ok {
				span.Process.Tags = append(span.Process.Tags, tag)
				continue
			}
			if i != newI {
				span.Tags[newI] = tag
			}
			newI++
		}
		span.Tags = span.Tags[:newI]
	}

	return Func(func(trace *model.Trace) (*model.Trace, error) {
		for _, span := range trace.Spans {
			adjustSpanTags(span)
			model.KeyValues(span.Process.Tags).Sort()
		}
		return trace, nil
	})
}
