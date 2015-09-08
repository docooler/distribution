package reference

import (
	"regexp"
	"strings"
	"testing"
)

type regexpMatch struct {
	input string
	match bool
	subs  []string
}

func checkRegexp(t *testing.T, r *regexp.Regexp, m regexpMatch) {
	matches := r.FindStringSubmatch(m.input)
	if m.match && matches != nil {
		if len(matches) != (r.NumSubexp()+1) || matches[0] != m.input {
			t.Fatalf("Bad match result: %#v", matches)
		}
		if len(matches) < (len(m.subs) + 1) {
			t.Errorf("Expected %d sub matches, only have %d", len(m.subs), len(matches)-1)
		}
		for i := range m.subs {
			if m.subs[i] != matches[i+1] {
				t.Errorf("Unexpected submatch %d: %q, expected %q", i+1, matches[i+1], m.subs[i])
			}
		}
	} else if m.match {
		t.Errorf("Expected match for %q", m.input)
	} else if matches != nil {
		t.Errorf("Unexpected match for %q", m.input)
	}
}

func TestHostRegexp(t *testing.T) {
	hostcases := []regexpMatch{
		{
			input: "test.com",
			match: true,
		},
		{
			input: "test.com:10304",
			match: true,
		},
		{
			input: "test.com:http",
			match: false,
		},
		{
			input: "localhost",
			match: true,
		},
		{
			input: "localhost:8080",
			match: true,
		},
		{
			input: "a",
			match: true,
		},
		{
			input: "a.b",
			match: true,
		},
		{
			input: "ab.cd.com",
			match: true,
		},
		{
			input: "a-b.com",
			match: true,
		},
		{
			input: "-ab.com",
			match: false,
		},
		{
			input: "ab-.com",
			match: false,
		},
		{
			input: "ab.c-om",
			match: true,
		},
		{
			input: "ab.-com",
			match: false,
		},
		{
			input: "ab.com-",
			match: false,
		},
		{
			input: "0101.com",
			match: true, // TODO(dmcgowan): valid if this should be allowed
		},
		{
			input: "001a.com",
			match: true,
		},
		{
			input: "b.gbc.io:443",
			match: true,
		},
		{
			input: "b.gbc.io",
			match: true,
		},
		{
			input: "xn--n3h.com", // ☃.com in punycode
			match: true,
		},
	}
	r := regexp.MustCompile(`^` + hostnameRegexp.String() + `$`)
	for i := range hostcases {
		checkRegexp(t, r, hostcases[i])
	}
}

func TestFullNameRegexp(t *testing.T) {
	testcases := []regexpMatch{
		{
			input: "",
			match: false,
		},
		{
			input: "short",
			match: true,
			subs:  []string{"", "short"},
		},
		{
			input: "simple/name",
			match: true,
			subs:  []string{"simple", "name"},
		},
		{
			input: "library/ubuntu",
			match: true,
			subs:  []string{"library", "ubuntu"},
		},
		{
			input: "docker/stevvooe/app",
			match: true,
			subs:  []string{"docker", "stevvooe/app"},
		},
		{
			input: "aa/aa/aa/aa/aa/aa/aa/aa/aa/bb/bb/bb/bb/bb/bb",
			match: true,
			subs:  []string{"aa", "aa/aa/aa/aa/aa/aa/aa/aa/bb/bb/bb/bb/bb/bb"},
		},
		{
			input: "aa/aa/bb/bb/bb",
			match: true,
			subs:  []string{"aa", "aa/bb/bb/bb"},
		},
		{
			input: "a/a/a/a",
			match: true,
			subs:  []string{"a", "a/a/a"},
		},
		{
			input: "a/a/a/a/",
			match: false,
		},
		{
			input: "a//a/a",
			match: false,
		},
		{
			input: "a",
			match: true,
			subs:  []string{"", "a"},
		},
		{
			input: "a/aa",
			match: true,
			subs:  []string{"a", "aa"},
		},
		{
			input: "a/aa/a",
			match: true,
			subs:  []string{"a", "aa/a"},
		},
		{
			input: "foo.com",
			match: true,
			subs:  []string{"", "foo.com"},
		},
		{
			input: "foo.com/",
			match: false,
		},
		{
			input: "foo.com:8080/bar",
			match: true,
			subs:  []string{"foo.com:8080", "bar"},
		},
		{
			input: "foo.com:http/bar",
			match: false,
		},
		{
			input: "foo.com/bar",
			match: true,
			subs:  []string{"foo.com", "bar"},
		},
		{
			input: "foo.com/bar/baz",
			match: true,
			subs:  []string{"foo.com", "bar/baz"},
		},
		{
			input: "localhost:8080/bar",
			match: true,
			subs:  []string{"localhost:8080", "bar"},
		},
		{
			input: "sub-dom1.foo.com/bar/baz/quux",
			match: true,
			subs:  []string{"sub-dom1.foo.com", "bar/baz/quux"},
		},
		{
			input: "blog.foo.com/bar/baz",
			match: true,
			subs:  []string{"blog.foo.com", "bar/baz"},
		},
		{
			input: "a^a",
			match: false,
		},
		{
			input: "aa/asdf$$^/aa",
			match: false,
		},
		{
			input: "asdf$$^/aa",
			match: false,
		},
		{
			input: "aa-a/a",
			match: true,
			subs:  []string{"aa-a", "a"},
		},
		{
			input: strings.Repeat("a/", 128) + "a",
			match: true,
			subs:  []string{"a", strings.Repeat("a/", 127) + "a"},
		},
		{
			input: "a-/a/a/a",
			match: false,
		},
		{
			input: "foo.com/a-/a/a",
			match: false,
		},
		{
			input: "-foo/bar",
			match: false,
		},
		{
			input: "foo/bar-",
			match: false,
		},
		{
			input: "foo-/bar",
			match: false,
		},
		{
			input: "foo/-bar",
			match: false,
		},
		{
			input: "_foo/bar",
			match: false,
		},
		{
			input: "foo_bar",
			match: true,
			subs:  []string{"", "foo_bar"},
		},
		{
			input: "foo_bar.com",
			match: true,
			subs:  []string{"", "foo_bar.com"},
		},
		{
			input: "foo_bar.com:8080",
			match: false,
		},
		{
			input: "foo_bar.com:8080/app",
			match: false,
		},
		{
			input: "foo.com/foo_bar",
			match: true,
			subs:  []string{"foo.com", "foo_bar"},
		},
		{
			input: "____/____",
			match: false,
		},
		{
			input: "_docker/_docker",
			match: false,
		},
		{
			input: "docker_/docker_",
			match: false,
		},
		{
			input: "b.gcr.io/test.example.com/my-app",
			match: true,
			subs:  []string{"b.gcr.io", "test.example.com/my-app"},
		},
		{
			input: "xn--n3h.com/myimage", // ☃.com in punycode
			match: true,
			subs:  []string{"xn--n3h.com", "myimage"},
		},
		{
			input: "xn--7o8h.com/myimage", // 🐳.com in punycode
			match: true,
			subs:  []string{"xn--7o8h.com", "myimage"},
		},
	}
	for i := range testcases {
		checkRegexp(t, anchoredNameRegexp, testcases[i])
	}
}

func TestReferenceRegexp(t *testing.T) {
	testcases := []regexpMatch{
		{
			input: "registry.com:8080/myapp:tag",
			match: true,
			subs:  []string{"registry.com:8080/myapp", "tag", ""},
		},
		{
			input: "registry.com:8080/myapp@sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912",
			match: true,
			subs:  []string{"registry.com:8080/myapp", "", "sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912"},
		},
		{
			input: "registry.com:8080/myapp:tag2@sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912",
			match: true,
			subs:  []string{"registry.com:8080/myapp", "tag2", "sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912"},
		},
		{
			input: "registry.com:8080/myapp@sha256:badbadbadbad",
			match: false,
		},
		{
			input: "registry.com:8080/myapp:invalid~tag",
			match: false,
		},
		{
			input: "bad_hostname.com:8080/myapp:tag",
			match: false,
		},
		{
			input:// localhost treated as name, missing tag with 8080 as tag
			"localhost:8080@sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912",
			match: true,
			subs:  []string{"localhost", "8080", "sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912"},
		},
		{
			input: "localhost:8080/name@sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912",
			match: true,
			subs:  []string{"localhost:8080/name", "", "sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912"},
		},
		{
			input: "localhost:http/name@sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912",
			match: false,
		},
		{
			// localhost will be treated as an image name without a host
			input: "localhost@sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912",
			match: true,
			subs:  []string{"localhost", "", "sha256:be178c0543eb17f5f3043021c9e5fcf30285e557a4fc309cce97ff9ca6182912"},
		},
		{
			input: "registry.com:8080/myapp@bad",
			match: false,
		},
		{
			input: "registry.com:8080/myapp@2bad",
			match: false, // TODO(dmcgowan): Support this as valid
		},
	}

	for i := range testcases {
		checkRegexp(t, ReferenceRegexp, testcases[i])
	}

}
