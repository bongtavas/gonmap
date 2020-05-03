package netutil

import "testing"

func TestBuildIPListCIDR(t *testing.T) {
	cidr := "202.92.128.0/24"
	t.Logf("Building IP List from %s", cidr)
	ips, err := BuildIPList(cidr)

	if err != nil {
		t.Errorf("BuildIPList failed: %s", err)
	}

	if len(ips) != 256 {
		t.Errorf("Expected: 256, Actual : %d", len(ips))
	}

	for _, ip := range ips {
		t.Logf("%s", ip)
	}
}

func TestBuildIPListIPRange(t *testing.T) {
	iprange := "202.92.128.0-255"
	t.Logf("Building IP List from %s", iprange)
	ips, err := BuildIPList(iprange)

	if err != nil {
		t.Errorf("BuildIPList failed: %s", err)
	}

	if len(ips) != 256 {
		t.Errorf("Expected: 256, Actual : %d", len(ips))
	}
	for _, ip := range ips {
		t.Logf("%s", ip)
	}
}

// func TestResolveHostname(T *testing.T) {
// 	resolveHostname("pleni.upd.edu.ph")
// 	resolveHostname("202.92.128.181")
// 	resolveHostname("202.92.128.1-254")
// 	resolveHostname("202.92.128.0/24")
// }
