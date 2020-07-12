// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	"strings"
	"testing"

	"github.com/marmotedu/apimachinery/pkg/validation/field"
)

func TestIsDNS1123Label(t *testing.T) {
	goodValues := []string{
		"a", "ab", "abc", "a1", "a-1", "a--1--2--b",
		"0", "01", "012", "1a", "1-a", "1--a--b--2",
		strings.Repeat("a", 63),
	}
	for _, val := range goodValues {
		if msgs := IsDNS1123Label(val); len(msgs) != 0 {
			t.Errorf("expected true for '%s': %v", val, msgs)
		}
	}

	badValues := []string{
		"", "A", "ABC", "aBc", "A1", "A-1", "1-A",
		"-", "a-", "-a", "1-", "-1",
		"_", "a_", "_a", "a_b", "1_", "_1", "1_2",
		".", "a.", ".a", "a.b", "1.", ".1", "1.2",
		" ", "a ", " a", "a b", "1 ", " 1", "1 2",
		strings.Repeat("a", 64),
	}
	for _, val := range badValues {
		if msgs := IsDNS1123Label(val); len(msgs) == 0 {
			t.Errorf("expected false for '%s'", val)
		}
	}
}

func TestIsDNS1123Subdomain(t *testing.T) {
	goodValues := []string{
		"a", "ab", "abc", "a1", "a-1", "a--1--2--b",
		"0", "01", "012", "1a", "1-a", "1--a--b--2",
		"a.a", "ab.a", "abc.a", "a1.a", "a-1.a", "a--1--2--b.a",
		"a.1", "ab.1", "abc.1", "a1.1", "a-1.1", "a--1--2--b.1",
		"0.a", "01.a", "012.a", "1a.a", "1-a.a", "1--a--b--2",
		"0.1", "01.1", "012.1", "1a.1", "1-a.1", "1--a--b--2.1",
		"a.b.c.d.e", "aa.bb.cc.dd.ee", "1.2.3.4.5", "11.22.33.44.55",
		strings.Repeat("a", 253),
	}
	for _, val := range goodValues {
		if msgs := IsDNS1123Subdomain(val); len(msgs) != 0 {
			t.Errorf("expected true for '%s': %v", val, msgs)
		}
	}

	badValues := []string{
		"", "A", "ABC", "aBc", "A1", "A-1", "1-A",
		"-", "a-", "-a", "1-", "-1",
		"_", "a_", "_a", "a_b", "1_", "_1", "1_2",
		".", "a.", ".a", "a..b", "1.", ".1", "1..2",
		" ", "a ", " a", "a b", "1 ", " 1", "1 2",
		"A.a", "aB.a", "ab.A", "A1.a", "a1.A",
		"A.1", "aB.1", "A1.1", "1A.1",
		"0.A", "01.A", "012.A", "1A.a", "1a.A",
		"A.B.C.D.E", "AA.BB.CC.DD.EE", "a.B.c.d.e", "aa.bB.cc.dd.ee",
		"a@b", "a,b", "a_b", "a;b",
		"a:b", "a%b", "a?b", "a$b",
		strings.Repeat("a", 254),
	}
	for _, val := range badValues {
		if msgs := IsDNS1123Subdomain(val); len(msgs) == 0 {
			t.Errorf("expected false for '%s'", val)
		}
	}
}

func TestIsValidPortNum(t *testing.T) {
	goodValues := []int{1, 2, 1000, 16384, 32768, 65535}
	for _, val := range goodValues {
		if msgs := IsValidPortNum(val); len(msgs) != 0 {
			t.Errorf("expected true for %d, got %v", val, msgs)
		}
	}

	badValues := []int{0, -1, 65536, 100000}
	for _, val := range badValues {
		if msgs := IsValidPortNum(val); len(msgs) == 0 {
			t.Errorf("expected false for %d", val)
		}
	}
}

func TestIsInRange(t *testing.T) {
	goodValues := []struct {
		value int
		min   int
		max   int
	}{{1, 0, 10}, {5, 5, 20}, {25, 10, 25}}
	for _, val := range goodValues {
		if msgs := IsInRange(val.value, val.min, val.max); len(msgs) > 0 {
			t.Errorf("expected no errors for %#v, but got %v", val, msgs)
		}
	}

	badValues := []struct {
		value int
		min   int
		max   int
	}{{1, 2, 10}, {5, -4, 2}, {25, 100, 120}}
	for _, val := range badValues {
		if msgs := IsInRange(val.value, val.min, val.max); len(msgs) == 0 {
			t.Errorf("expected errors for %#v", val)
		}
	}
}

func TestIsQualifiedName(t *testing.T) {
	successCases := []string{
		"simple",
		"now-with-dashes",
		"1-starts-with-num",
		"1234",
		"simple/simple",
		"now-with-dashes/simple",
		"now-with-dashes/now-with-dashes",
		"now.with.dots/simple",
		"now-with.dashes-and.dots/simple",
		"1-num.2-num/3-num",
		"1234/5678",
		"1.2.3.4/5678",
		"Uppercase_Is_OK_123",
		"example.com/Uppercase_Is_OK_123",
		"requests.storage-foo",
		strings.Repeat("a", 63),
		strings.Repeat("a", 253) + "/" + strings.Repeat("b", 63),
	}
	for i := range successCases {
		if errs := IsQualifiedName(successCases[i]); len(errs) != 0 {
			t.Errorf("case[%d]: %q: expected success: %v", i, successCases[i], errs)
		}
	}

	errorCases := []string{
		"nospecialchars%^=@",
		"cantendwithadash-",
		"-cantstartwithadash-",
		"only/one/slash",
		"Example.com/abc",
		"example_com/abc",
		"example.com/",
		"/simple",
		strings.Repeat("a", 64),
		strings.Repeat("a", 254) + "/abc",
	}
	for i := range errorCases {
		if errs := IsQualifiedName(errorCases[i]); len(errs) == 0 {
			t.Errorf("case[%d]: %q: expected failure", i, errorCases[i])
		}
	}
}

func TestIsValidIP(t *testing.T) {
	goodValues := []string{
		"::1",
		"2a00:79e0:2:0:f1c3:e797:93c1:df80",
		"::",
		"2001:4860:4860::8888",
		"::fff:1.1.1.1",
		"1.1.1.1",
		"1.1.1.01",
		"255.0.0.1",
		"1.0.0.0",
		"0.0.0.0",
	}
	for _, val := range goodValues {
		if msgs := IsValidIP(val); len(msgs) != 0 {
			t.Errorf("expected true for %q: %v", val, msgs)
		}
	}

	badValues := []string{
		"[2001:db8:0:1]:80",
		"myhost.mydomain",
		"-1.0.0.0",
		"[2001:db8:0:1]",
		"a",
	}
	for _, val := range badValues {
		if msgs := IsValidIP(val); len(msgs) == 0 {
			t.Errorf("expected false for %q", val)
		}
	}
}

func TestIsValidIPv4Address(t *testing.T) {
	goodValues := []string{
		"1.1.1.1",
		"1.1.1.01",
		"255.0.0.1",
		"1.0.0.0",
		"0.0.0.0",
	}
	for _, val := range goodValues {
		if msgs := IsValidIPv4Address(field.NewPath(""), val); len(msgs) != 0 {
			t.Errorf("expected %q to be valid IPv4 address: %v", val, msgs)
		}
	}

	badValues := []string{
		"[2001:db8:0:1]:80",
		"myhost.mydomain",
		"-1.0.0.0",
		"[2001:db8:0:1]",
		"a",
		"2001:4860:4860::8888",
		"::fff:1.1.1.1",
		"::1",
		"2a00:79e0:2:0:f1c3:e797:93c1:df80",
		"::",
	}
	for _, val := range badValues {
		if msgs := IsValidIPv4Address(field.NewPath(""), val); len(msgs) == 0 {
			t.Errorf("expected %q to be invalid IPv4 address", val)
		}
	}
}

func TestIsValidIPv6Address(t *testing.T) {
	goodValues := []string{
		"2001:4860:4860::8888",
		"2a00:79e0:2:0:f1c3:e797:93c1:df80",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"::fff:1.1.1.1",
		"::1",
		"::",
	}

	for _, val := range goodValues {
		if msgs := IsValidIPv6Address(field.NewPath(""), val); len(msgs) != 0 {
			t.Errorf("expected %q to be valid IPv6 address: %v", val, msgs)
		}
	}

	badValues := []string{
		"1.1.1.1",
		"1.1.1.01",
		"255.0.0.1",
		"1.0.0.0",
		"0.0.0.0",
		"[2001:db8:0:1]:80",
		"myhost.mydomain",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334:2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"-1.0.0.0",
		"[2001:db8:0:1]",
		"a",
	}
	for _, val := range badValues {
		if msgs := IsValidIPv6Address(field.NewPath(""), val); len(msgs) == 0 {
			t.Errorf("expected %q to be invalid IPv6 address", val)
		}
	}
}

func TestIsValidPercent(t *testing.T) {
	goodValues := []string{
		"0%",
		"00000%",
		"1%",
		"01%",
		"99%",
		"100%",
		"101%",
	}
	for _, val := range goodValues {
		if msgs := IsValidPercent(val); len(msgs) != 0 {
			t.Errorf("expected true for %q: %v", val, msgs)
		}
	}

	badValues := []string{
		"",
		"0",
		"100",
		"0.0%",
		"99.9%",
		"hundred",
		" 1%",
		"1% ",
		"-0%",
		"-1%",
		"+1%",
	}
	for _, val := range badValues {
		if msgs := IsValidPercent(val); len(msgs) == 0 {
			t.Errorf("expected false for %q", val)
		}
	}
}
