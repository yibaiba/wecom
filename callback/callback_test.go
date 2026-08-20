package callback

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"strings"
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

func TestPKCS7PadsTo32(t *testing.T) {
	got := pkcs7Pad([]byte{1}, 32)
	if len(got) != 32 || got[31] != 31 {
		t.Fatalf("pad-32 len=%d last=%d", len(got), got[len(got)-1])
	}
}

func TestCallbackOfficialSampleDecrypt(t *testing.T) {
	cb := Callback{
		Token:          "QDG6eK",
		EncodingAESKey: "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C",
		CorpID:         "wx5823bf96d3bd56c7",
	}
	const encrypt = "RypEvHKD8QQKFhvQ6QleEB4J58tiPdvo+rtK1I9qca6aM/wvqnLSV5zEPeusUiX5L5X/0lWfrf0QADHHhGd3QczcdCUpj911L3vg3W/sYYvuJTs3TUUkSUXxaccAS0qhxchrRYt66wiSpGLYL42aM6A8dTT+6k4aSknmPj48kzJs8qLjvd4Xgpue06DOdnLxAUHzM6+kDZ+HMZfJYuR+LtwGc2hgf5gsijff0ekUNXZiqATP7PF5mZxZ3Izoun1s4zG4LUMnvw2r+KqCKIw+3IQH03v+BCA9nMELNqbSf6tiWSrXJB3LAVGUcallcrw8V2t9EL4EhzJWrQUax5wLVMNS0+rUPA3k22Ncx4XXZS9o0MBH27Bo6BpNelZpS+/uh9KsNlY6bHCmJU9p8g7m3fVKn28H3KDYA5Pl/T8Z1ptDAVe0lXdQ2YoyyH2uyPIGHBZZIs2pDBS8R07+qN+E7Q=="
	const msgSig = "477715d11cdb4164915debcba66cb864d751f3e6"
	body := []byte(`<xml><Encrypt><![CDATA[` + encrypt + `]]></Encrypt></xml>`)
	out, err := cb.Decrypt(msgSig, "1409659813", "1372623149", body)
	if err != nil {
		t.Fatal(err)
	}
	plain := string(out)
	if !strings.Contains(plain, "<Content><![CDATA[hello]]></Content>") {
		t.Fatalf("content %s", plain)
	}
	if !strings.Contains(plain, "<FromUserName><![CDATA[mycreate]]></FromUserName>") {
		t.Fatalf("from %s", plain)
	}
}
