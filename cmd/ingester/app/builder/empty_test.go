// Copyright (c) 2023 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"testing"

	"github.com/jaegertracing/jaeger/pkg/testutils"
)

func TestMain(m *testing.M) {
	testutils.VerifyGoLeaks(m)
}
