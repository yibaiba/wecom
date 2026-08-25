package login

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

// This mirrors the host's page-policy helper. It is duplicated rather than
// shared because this module must not depend on the SSO host, and web security
// policy is not a concept this module should own and export.

func scriptHash(body string) string {
	sum := sha256.Sum256([]byte(body))
	return "'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'"
}

// pagePolicy admits only the given inline scripts. Callers must pass the exact
// strings the page ships so the policy cannot drift from the markup.
func pagePolicy(scripts ...string) string {
	hashes := make([]string, 0, len(scripts))
	for _, s := range scripts {
		hashes = append(hashes, scriptHash(s))
	}
	sources := "'none'"
	if len(hashes) > 0 {
		sources = strings.Join(hashes, " ")
	}
	return "default-src 'self'; script-src " + sources +
		"; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'"
}

// jssdkOrigin is the CDN the login panel pins its copy of the WeCom JSSDK to.
// It is spelled out here so the policy and the markup cannot disagree about it.
const jssdkOrigin = "https://wwcdn.weixin.qq.com"

// panelPolicy is the login panel's policy. The panel cannot use pagePolicy: it
// pulls the WeCom JSSDK from a CDN, and that SDK then builds an iframe and talks
// to hosts of its own choosing.
//
// script-src is pinned tightly, because it is both the directive that actually
// stops injected code and the one we can state exactly: two inline scripts plus
// the one CDN this page names. frame-src and connect-src are left at https:
// instead. Their real values belong to the SDK rather than to us, so any list we
// wrote would be a guess that a future SDK release could invalidate — and a wrong
// guess there takes down every login. They also buy little on this page, since
// reaching them at all requires code that script-src already refuses to run.
func panelPolicy(scripts ...string) string {
	hashes := make([]string, 0, len(scripts)+1)
	for _, s := range scripts {
		hashes = append(hashes, scriptHash(s))
	}
	hashes = append(hashes, jssdkOrigin)
	return "default-src 'self'; script-src " + strings.Join(hashes, " ") +
		"; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; connect-src https:; frame-src https:" +
		"; object-src 'none'; form-action 'self'; frame-ancestors 'none'; base-uri 'none'"
}
