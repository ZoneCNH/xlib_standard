// SPDX-License-Identifier: Apache-2.0
package main

import (
	"reflect"
	"testing"
)

func TestEnvBoolParsesKnownValuesAndFallback(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		fallback bool
		want     bool
	}{
		{name: "true digit", value: "1", fallback: false, want: true},
		{name: "true word", value: " true ", fallback: false, want: true},
		{name: "true yes", value: "YES", fallback: false, want: true},
		{name: "true on", value: "on", fallback: false, want: true},
		{name: "false digit", value: "0", fallback: true, want: false},
		{name: "false word", value: " false ", fallback: true, want: false},
		{name: "false no", value: "NO", fallback: true, want: false},
		{name: "false off", value: "off", fallback: true, want: false},
		{name: "empty uses true fallback", value: "", fallback: true, want: true},
		{name: "empty uses false fallback", value: "", fallback: false, want: false},
		{name: "unknown uses fallback", value: "maybe", fallback: true, want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			key := "RELEASEMANIFEST_TEST_ENV_BOOL"
			t.Setenv(key, tc.value)

			got := envBool(key, tc.fallback)

			if got != tc.want {
				t.Fatalf("envBool(%q, %v) = %v, want %v", tc.value, tc.fallback, got, tc.want)
			}
		})
	}
}

func TestEnvCSVDefaultTrimsFiltersAndCopiesFallback(t *testing.T) {
	fallback := []string{"fallback-a", "fallback-b"}

	t.Run("blank returns fallback copy", func(t *testing.T) {
		key := "RELEASEMANIFEST_TEST_ENV_CSV_BLANK"
		t.Setenv(key, " ")

		got := envCSVDefault(key, fallback)

		if !reflect.DeepEqual(got, fallback) {
			t.Fatalf("envCSVDefault blank = %v, want %v", got, fallback)
		}
		got[0] = "mutated"
		if fallback[0] != "fallback-a" {
			t.Fatalf("fallback was mutated: %v", fallback)
		}
	})

	t.Run("trims and drops empty fields", func(t *testing.T) {
		key := "RELEASEMANIFEST_TEST_ENV_CSV_VALUES"
		t.Setenv(key, " alpha, , beta ,gamma ")

		got := envCSVDefault(key, fallback)
		want := []string{"alpha", "beta", "gamma"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("envCSVDefault values = %v, want %v", got, want)
		}
	})

	t.Run("only empty fields returns fallback copy", func(t *testing.T) {
		key := "RELEASEMANIFEST_TEST_ENV_CSV_EMPTY_FIELDS"
		t.Setenv(key, " , , ")

		got := envCSVDefault(key, fallback)

		if !reflect.DeepEqual(got, fallback) {
			t.Fatalf("envCSVDefault empty fields = %v, want %v", got, fallback)
		}
		got[1] = "mutated"
		if fallback[1] != "fallback-b" {
			t.Fatalf("fallback was mutated: %v", fallback)
		}
	})
}
