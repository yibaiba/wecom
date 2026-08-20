package callback

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"testing"
)

func TestCallbackEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	aesKey := base64.StdEncoding.EncodeToString(key)
	cb := Callback{Token: "token", EncodingAESKey: aesKey[:len(aesKey)-1], CorpID: "wwcorp"}
	plain := []byte("<xml><Content>hi</Content></xml>")
	enc, err := cb.Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	got, err := cb.decrypt(enc)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(plain) {
		t.Fatalf("roundtrip %s", got)
	}
	body, _ := xml.Marshal(callbackEnvelope{Encrypt: enc})
	sig := cb.Signature("1", "nonce", enc)
	out, err := cb.Decrypt(sig, "1", "nonce", body)
	if err != nil || string(out) != string(plain) {
		t.Fatalf("decrypt post %s %v", out, err)
	}
	if _, err := cb.Decrypt("bad", "1", "nonce", body); err == nil {
		t.Fatal("expected signature error")
	}
}
