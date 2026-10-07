package community

import (
	"math"
	"testing"
)

func permissionSnapshot() *permSnapshot {
	return &permSnapshot{
		owners: map[string]bool{},
		roles: map[int64]*role{
			everyoneRoleID: {ID: everyoneRoleID, perms: permViewChannel | permSendMessages},
		},
		memberRoles: map[string][]int64{},
		channels:    map[int64]*channel{},
		restricted:  map[string]string{},
	}
}

func TestPermissionSnapshotRoleOverrideAggregation(t *testing.T) {
	ps := permissionSnapshot()
	ps.roles[2] = &role{ID: 2, perms: 0}
	ps.roles[3] = &role{ID: 3, perms: 0}
	ps.memberRoles["member"] = []int64{2, 3}
	c := &channel{Overrides: []override{
		{Type: targetRole, TargetID: "2", allow: permSendMessages},
		{Type: targetRole, TargetID: "3", deny: permSendMessages},
	}}
	if got := ps.applyOverrides(permViewChannel, c, "member"); got&permSendMessages == 0 {
		t.Fatal("role allow should win over a deny from another member role")
	}
}

func TestPermissionSnapshotOverridePrecedence(t *testing.T) {
	ps := permissionSnapshot()
	ps.roles[2] = &role{ID: 2}
	ps.memberRoles["member"] = []int64{2}
	c := &channel{Overrides: []override{
		{Type: targetRole, TargetID: "1", deny: permSendMessages}, // everyone
		{Type: targetRole, TargetID: "2", allow: permSendMessages},
		{Type: targetMember, TargetID: "member", deny: permSendMessages},
	}}
	if got := ps.applyOverrides(permViewChannel|permSendMessages, c, "member"); got&permSendMessages != 0 {
		t.Fatal("member deny should override everyone and role allow")
	}
	c.Overrides[2].allow = permSendMessages
	if got := ps.applyOverrides(permViewChannel, c, "member"); got&permSendMessages == 0 {
		t.Fatal("member allow should override earlier everyone and role rules")
	}
}

func TestPermissionSnapshotOwnerAdministratorAndRestrictions(t *testing.T) {
	t.Run("owner bypasses restrictions", func(t *testing.T) {
		ps := permissionSnapshot()
		ps.owners["owner"] = true
		ps.restricted["owner"] = restrictTimeout
		if got := ps.base("owner"); got != permAll {
			t.Fatalf("owner base = %b, want all permissions", got)
		}
	})

	t.Run("administrator bypasses channel overrides", func(t *testing.T) {
		ps := permissionSnapshot()
		ps.roles[2] = &role{ID: 2, perms: permAdministrator}
		ps.memberRoles["admin"] = []int64{2}
		ps.restricted["admin"] = restrictTimeout
		ps.channels[10] = &channel{ID: 10, Type: chanText, Overrides: []override{{Type: targetMember, TargetID: "admin", deny: permViewChannel}}}
		if got := ps.inChannel("admin", 10); got != permAll {
			t.Fatalf("administrator channel permissions = %b, want all permissions", got)
		}
	})

	t.Run("restriction is reapplied after channel allows", func(t *testing.T) {
		ps := permissionSnapshot()
		ps.restricted["limited"] = restrictRules
		ps.channels[10] = &channel{ID: 10, Type: chanText, Overrides: []override{{Type: targetMember, TargetID: "limited", allow: permSendMessages}}}
		if got := ps.inChannel("limited", 10); got != permViewChannel {
			t.Fatalf("restricted channel permissions = %b, want read-only", got)
		}
	})
}

func TestPermissionSnapshotViewCategoryAndThreadInheritance(t *testing.T) {
	ps := permissionSnapshot()
	ps.channels[10] = &channel{ID: 10, Type: chanCategory, Overrides: []override{{Type: targetMember, TargetID: "hidden", deny: permViewChannel}}}
	parent := int64(10)
	ps.channels[11] = &channel{ID: 11, Type: chanText, ParentID: &parent}
	ps.channels[12] = &channel{ID: 12, Type: chanThread, ParentID: ptrInt64(11), Overrides: []override{{Type: targetMember, TargetID: "visible", allow: permSendMessages}}}

	if got := ps.inChannel("hidden", 11); got != 0 {
		t.Fatalf("channel without view permission returned %b, want no permissions", got)
	}
	if got := ps.inChannel("visible", 12); got != permViewChannel|permSendMessages {
		t.Fatalf("thread permissions = %b, want parent-channel permissions", got)
	}
}

func TestPermissionSnapshotRoleHierarchy(t *testing.T) {
	ps := permissionSnapshot()
	ps.owners["owner"] = true
	ps.owners["other-owner"] = true
	ps.roles[2] = &role{ID: 2, Position: 2}
	ps.roles[3] = &role{ID: 3, Position: 7}
	ps.memberRoles["actor"] = []int64{2, 3}
	ps.memberRoles["target"] = []int64{2}
	ps.memberRoles["same"] = []int64{3}

	if !ps.outranks("actor", "target") || ps.outranks("actor", "same") || ps.outranks("actor", "other-owner") {
		t.Fatal("role hierarchy should require a strictly higher role and never outrank an owner")
	}
	if ps.top("owner") != math.MaxInt64 {
		t.Fatal("owner should rank above every role")
	}
}

func TestPermissionNamesRejectUnknownAndRoundTripKnownFlags(t *testing.T) {
	allNames := permAll.names()
	got, err := parsePerms(allNames)
	if err != nil || got != permAll {
		t.Fatalf("known permission round trip = (%b, %v), want (%b, nil)", got, err, permAll)
	}
	if _, err := parsePerms([]string{"not_a_permission"}); err == nil {
		t.Fatal("unknown permission name was accepted")
	}
	if got := perm(1 << 29).names(); len(got) != 0 {
		t.Fatalf("unknown stored permission flag converted to names %v", got)
	}
}

func ptrInt64(v int64) *int64 { return &v }
