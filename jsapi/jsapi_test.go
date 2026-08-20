package jsapi

import "testing"

func TestSignature(t *testing.T) {
	got := Signature("sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90e5h", "Wm3WZYTPz0wzccnW", "http://mp.weixin.qq.com?params=value", 1414587457)
	if got == "" || len(got) != 40 {
		t.Fatalf("sig %s", got)
	}
}
