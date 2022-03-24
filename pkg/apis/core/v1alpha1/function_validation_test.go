// Copyright 2022 The hyperfunction Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/ptr"
)

func TestFunctionValidation(t *testing.T) {
	tests := []struct {
		name    string
		in      *Function
		wantErr *apis.FieldError
	}{
		{
			name: "code with no fields",
			in: &Function{
				Spec: FunctionSpec{
					Code: FunctionCode{},
				},
			},
			wantErr: apis.ErrMissingOneOf("spec.code.inline"),
		},
		{
			name: "simple ok function",
			in: &Function{
				Spec: FunctionSpec{
					Public: true,
					Code: FunctionCode{
						Inline: ptr.String("def handle(*args):\n  pass"),
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Validate(context.Background())
			assert.Equal(t, tt.wantErr.Error(), got.Error())
		})
	}
}
