package main

import "testing"

func TestBoundaryPolicyRejectsForbiddenEdges(t *testing.T) {
	packages := []goPackage{
		{ImportPath: module + "/internal/application/messages", Imports: []string{module + "/internal/adapters/postgres"}},
		{ImportPath: module + "/internal/adapters/postgres", Imports: []string{module + "/internal/community"}},
		{ImportPath: module + "/internal/transport/http", Imports: []string{module + "/internal/adapters/postgres"}},
		{ImportPath: module + "/internal/transport/http_test", XTestImports: []string{module + "/internal/community"}},
		{ImportPath: module + "/internal/platform/ids", Imports: []string{module + "/internal/community"}},
	}
	bad := violations(packages)
	if len(bad) != 5 {
		t.Fatalf("violations = %v, want forbidden core, adapter, transport, platform and external-test edges", bad)
	}
}

func TestBoundaryPolicyChecksInternalTestImports(t *testing.T) {
	packages := []goPackage{{
		ImportPath:  module + "/internal/platform/testfixture",
		TestImports: []string{module + "/internal/transport/http"},
	}}
	if bad := violations(packages); len(bad) != 1 {
		t.Fatalf("test import violations = %v, want application boundary enforced for tests", bad)
	}
}

func TestBoundaryPolicyAllowsRecordedLegacyEdges(t *testing.T) {
	packages := []goPackage{{
		ImportPath: module + "/internal/community",
		Imports:    []string{module + "/internal/httpapi", module + "/internal/realtime", module + "/internal/sqlitedb"},
	}}
	if bad := violations(packages); len(bad) != 0 {
		t.Fatalf("legacy service edges rejected: %v", bad)
	}
}

func TestValidationRequiresRealLayerPackages(t *testing.T) {
	if err := validate([]goPackage{{ImportPath: module + "/internal/community"}}); err == nil {
		t.Fatal("checker accepted a graph with no new platform/adapter packages")
	}
}
