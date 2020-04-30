/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package value

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/go-cmp/cmp"
	"k8s.io/component-base/metrics/legacyregistry"
	"k8s.io/component-base/metrics/testutil"
)

type testTransformer struct {
	from, to                 []byte
	err                      error
	stale                    bool
	receivedFrom, receivedTo []byte
}

func (t *testTransformer) TransformFromStorage(from []byte, context Context) (data []byte, stale bool, err error) {
	t.receivedFrom = from
	return t.from, t.stale, t.err
}

func (t *testTransformer) TransformToStorage(to []byte, context Context) (data []byte, err error) {
	t.receivedTo = to
	return t.to, t.err
}

func TestPrefixFrom(t *testing.T) {
	testErr := fmt.Errorf("test error")
	transformErr := status.Error(codes.Internal, "transform error")
	transformers := []PrefixTransformer{
		{Prefix: []byte("first:"), Transformer: &testTransformer{from: []byte("value1")}},
		{Prefix: []byte("second:"), Transformer: &testTransformer{from: []byte("value2")}},
		{Prefix: []byte("fails:"), Transformer: &testTransformer{err: transformErr}},
		{Prefix: []byte("stale:"), Transformer: &testTransformer{from: []byte("value3"), stale: true}},
	}
	p := NewPrefixTransformers(testErr, transformers...)

	testCases := []struct {
		desc     string
		input    []byte
		expect   []byte
		stale    bool
		err      error
		matchIdx int
	}{
		{
			desc:   "matches the first transformer",
			input:  []byte("first:value"),
			expect: []byte("value1"),
		},
		{
			desc:     "matches the second transformer",
			input:    []byte("second:value"),
			expect:   []byte("value2"),
			stale:    true,
			matchIdx: 1,
		},
		{
			desc:     "no match",
			input:    []byte("third:value"),
			err:      status.Error(codes.Internal, testErr.Error()),
			matchIdx: -1,
		},
		{
			desc:     "matches but fails",
			input:    []byte("fails:value"),
			err:      transformErr,
			matchIdx: 2,
		},
		{
			desc:     "stale",
			input:    []byte("stale:value"),
			expect:   []byte("value3"),
			stale:    true,
			matchIdx: 3,
		},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			got, stale, err := p.TransformFromStorage(test.input, nil)

			if d := cmp.Diff(test.err, err); d != "" {
				t.Errorf("error mismatch (-want +got):\n%s", d)
			}

			if d := cmp.Diff(test.expect, got); d != "" {
				t.Errorf("transformation mismatch (-want +got):\n%s", d)
			}

			if stale != test.stale {
				t.Errorf("got %t want %t for stale", stale, test.stale)
			}

			if test.matchIdx != -1 && !bytes.Equal([]byte("value"), transformers[test.matchIdx].Transformer.(*testTransformer).receivedFrom) {
				t.Errorf("unexpected value received by transformer: %s", transformers[test.matchIdx].Transformer.(*testTransformer).receivedFrom)
			}
		})
	}
}

func TestPrefixTo(t *testing.T) {
	testErr := fmt.Errorf("test error")
	transformErr := fmt.Errorf("test error")
	testCases := []struct {
		transformers []PrefixTransformer
		expect       []byte
		err          error
	}{
		{[]PrefixTransformer{{Prefix: []byte("first:"), Transformer: &testTransformer{to: []byte("value1")}}}, []byte("first:value1"), nil},
		{[]PrefixTransformer{{Prefix: []byte("second:"), Transformer: &testTransformer{to: []byte("value2")}}}, []byte("second:value2"), nil},
		{[]PrefixTransformer{{Prefix: []byte("fails:"), Transformer: &testTransformer{err: transformErr}}}, nil, transformErr},
	}
	for i, test := range testCases {
		p := NewPrefixTransformers(testErr, test.transformers...)
		got, err := p.TransformToStorage([]byte("value"), nil)
		if err != test.err || !bytes.Equal(got, test.expect) {
			t.Errorf("%d: unexpected out: %q %#v", i, string(got), err)
			continue
		}
		if !bytes.Equal([]byte("value"), test.transformers[0].Transformer.(*testTransformer).receivedTo) {
			t.Errorf("%d: unexpected value received by transformer: %s", i, test.transformers[0].Transformer.(*testTransformer).receivedTo)
		}
	}
}

func TestPrefixFromMetrics(t *testing.T) {
	testErr := fmt.Errorf("test error")
	transformerErr := fmt.Errorf("test error")
	identityTransformer := PrefixTransformer{Prefix: []byte{}, Transformer: &testTransformer{from: []byte("value1")}}
	identityTransformerErr := PrefixTransformer{Prefix: []byte{}, Transformer: &testTransformer{err: transformerErr}}
	otherTransformer := PrefixTransformer{Prefix: []byte("other:"), Transformer: &testTransformer{from: []byte("value1")}}
	otherTransformerErr := PrefixTransformer{Prefix: []byte("other:"), Transformer: &testTransformer{err: transformerErr}}

	testCases := []struct {
		desc    string
		input   []byte
		prefix  Transformer
		metrics []string
		want    string
		err     error
	}{
		{
			desc:   "identity prefix",
			input:  []byte("value"),
			prefix: NewPrefixTransformers(testErr, identityTransformer, otherTransformer),
			metrics: []string{
				"apiserver_storage_transformation_operations_total",
			},
			want: `
	# HELP apiserver_storage_transformation_operations_total [ALPHA] Total number of transformations.
  # TYPE apiserver_storage_transformation_operations_total counter
  apiserver_storage_transformation_operations_total{status="OK",transformation_type="from_storage",transformer_prefix="identity"} 1
  `,
			err: nil,
		},
		{
			desc:   "other prefix (ok)",
			input:  []byte("other:value"),
			prefix: NewPrefixTransformers(testErr, identityTransformerErr, otherTransformer),
			metrics: []string{
				"apiserver_storage_transformation_operations_total",
			},
			want: `
	# HELP apiserver_storage_transformation_operations_total [ALPHA] Total number of transformations.
  # TYPE apiserver_storage_transformation_operations_total counter
  apiserver_storage_transformation_operations_total{status="OK",transformation_type="from_storage",transformer_prefix="other:"} 1
  `,
			err: nil,
		},
		{
			desc:   "other prefix (error)",
			input:  []byte("other:value"),
			prefix: NewPrefixTransformers(testErr, identityTransformerErr, otherTransformerErr),
			metrics: []string{
				"apiserver_storage_transformation_operations_total",
			},
			want: `
	# HELP apiserver_storage_transformation_operations_total [ALPHA] Total number of transformations.
  # TYPE apiserver_storage_transformation_operations_total counter
  apiserver_storage_transformation_operations_total{status="Unknown",transformation_type="from_storage",transformer_prefix="other:"} 1
  `,
			err: nil,
		},
		{
			desc:   "unknown prefix",
			input:  []byte("foo:value"),
			prefix: NewPrefixTransformers(testErr, identityTransformerErr, otherTransformer),
			metrics: []string{
				"apiserver_storage_transformation_operations_total",
			},
			want: `
	# HELP apiserver_storage_transformation_operations_total [ALPHA] Total number of transformations.
  # TYPE apiserver_storage_transformation_operations_total counter
  apiserver_storage_transformation_operations_total{status="Internal",transformation_type="from_storage",transformer_prefix="unknown"} 1
  `,
			err: nil,
		},
	}

	RegisterMetrics()
	transformerOperationsTotal.Reset()

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.prefix.TransformFromStorage(tc.input, nil)
			defer transformerOperationsTotal.Reset()
			if err := testutil.GatherAndCompare(legacyregistry.DefaultGatherer, strings.NewReader(tc.want), tc.metrics...); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPrefixToMetrics(t *testing.T) {
	testErr := fmt.Errorf("test error")
	transformerErr := fmt.Errorf("test error")
	otherTransformer := PrefixTransformer{Prefix: []byte("other:"), Transformer: &testTransformer{from: []byte("value1")}}
	otherTransformerErr := PrefixTransformer{Prefix: []byte("other:"), Transformer: &testTransformer{err: transformerErr}}

	testCases := []struct {
		desc    string
		input   []byte
		prefix  Transformer
		metrics []string
		want    string
		err     error
	}{
		{
			desc:   "ok",
			input:  []byte("value"),
			prefix: NewPrefixTransformers(testErr, otherTransformer),
			metrics: []string{
				"apiserver_storage_transformation_operations_total",
			},
			want: `
	# HELP apiserver_storage_transformation_operations_total [ALPHA] Total number of transformations.
  # TYPE apiserver_storage_transformation_operations_total counter
  apiserver_storage_transformation_operations_total{status="OK",transformation_type="to_storage",transformer_prefix="other:"} 1
  `,
			err: nil,
		},
		{
			desc:   "error",
			input:  []byte("value"),
			prefix: NewPrefixTransformers(testErr, otherTransformerErr),
			metrics: []string{
				"apiserver_storage_transformation_operations_total",
			},
			want: `
	# HELP apiserver_storage_transformation_operations_total [ALPHA] Total number of transformations.
  # TYPE apiserver_storage_transformation_operations_total counter
  apiserver_storage_transformation_operations_total{status="Unknown",transformation_type="to_storage",transformer_prefix="other:"} 1
  `,
			err: nil,
		},
	}

	RegisterMetrics()
	transformerOperationsTotal.Reset()

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.prefix.TransformToStorage(tc.input, nil)
			defer transformerOperationsTotal.Reset()
			if err := testutil.GatherAndCompare(legacyregistry.DefaultGatherer, strings.NewReader(tc.want), tc.metrics...); err != nil {
				t.Fatal(err)
			}
		})
	}
}
