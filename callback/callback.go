package callback

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Callback verifies and decrypts WeCom application receive-message XML.
type Callback struct {
	Token          string
	EncodingAESKey string
	CorpID         string
}

func (c Callback) aesKey() ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(c.EncodingAESKey + "=")
	if err != nil {
		return nil, fmt.Errorf("wecom encodingAESKey: %w", err)
	}
	if len(raw) != 32 {
		return nil, fmt.Errorf("wecom encodingAESKey must decode to 32 bytes")
	}
	return raw, nil
}

// Signature is the official SHA1 of sorted token/timestamp/nonce/encrypt.
func (c Callback) Signature(timestamp, nonce, encrypt string) string {
	parts := []string{c.Token, timestamp, nonce, encrypt}
	sort.Strings(parts)
	sum := sha1.Sum([]byte(strings.Join(parts, "")))
	return fmt.Sprintf("%x", sum)
}

// VerifyURL decrypts the GET echostr used when configuring the receive URL.
func (c Callback) VerifyURL(msgSig, timestamp, nonce, echostr string) (string, error) {
	if c.Signature(timestamp, nonce, echostr) != msgSig {
		return "", fmt.Errorf("wecom callback signature mismatch")
	}
	plain, err := c.decrypt(echostr)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

type callbackEnvelope struct {
	XMLName xml.Name `xml:"xml"`
	Encrypt string   `xml:"Encrypt"`
}

// Decrypt reads a POST receive-message body and returns the inner XML.
func (c Callback) Decrypt(msgSig, timestamp, nonce string, postBody []byte) ([]byte, error) {
	var env callbackEnvelope
	if err := xml.Unmarshal(postBody, &env); err != nil {
		return nil, fmt.Errorf("wecom callback xml: %w", err)
	}
	if c.Signature(timestamp, nonce, env.Encrypt) != msgSig {
		return nil, fmt.Errorf("wecom callback signature mismatch")
	}
	return c.decrypt(env.Encrypt)
}

// Encrypt packs plaintext for a reply (random + length + msg + corpid).
func (c Callback) Encrypt(plain []byte) (string, error) {
	key, err := c.aesKey()
	if err != nil {
		return "", err
	}
	rand16 := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, rand16); err != nil {
		return "", err
	}
	buf := make([]byte, 0, 16+4+len(plain)+len(c.CorpID))
	buf = append(buf, rand16...)
	lenbuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenbuf, uint32(len(plain)))
	buf = append(buf, lenbuf...)
	buf = append(buf, plain...)
	buf = append(buf, []byte(c.CorpID)...)
	padded := pkcs7Pad(buf, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
	out := make([]byte, len(padded))
	mode.CryptBlocks(out, padded)
	return base64.StdEncoding.EncodeToString(out), nil
}

func (c Callback) decrypt(b64 string) ([]byte, error) {
	key, err := c.aesKey()
	if err != nil {
		return nil, err
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	if len(raw)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("wecom callback ciphertext length")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plain := make([]byte, len(raw))
	cipher.NewCBCDecrypter(block, key[:aes.BlockSize]).CryptBlocks(plain, raw)
	plain, err = pkcs7Unpad(plain, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	if len(plain) < 20 {
		return nil, fmt.Errorf("wecom callback plaintext short")
	}
	msgLen := binary.BigEndian.Uint32(plain[16:20])
	end := 20 + int(msgLen)
	if end > len(plain) {
		return nil, fmt.Errorf("wecom callback message length")
	}
	if string(plain[end:]) != c.CorpID {
		return nil, fmt.Errorf("wecom callback corpid mismatch")
	}
	out := make([]byte, msgLen)
	copy(out, plain[20:end])
	return out, nil
}

func pkcs7Pad(b []byte, block int) []byte {
	n := block - (len(b) % block)
	pad := bytesRepeat(byte(n), n)
	return append(b, pad...)
}

func pkcs7Unpad(b []byte, block int) ([]byte, error) {
	if len(b) == 0 || len(b)%block != 0 {
		return nil, fmt.Errorf("wecom pkcs7 length")
	}
	n := int(b[len(b)-1])
	if n == 0 || n > block || n > len(b) {
		return nil, fmt.Errorf("wecom pkcs7 pad")
	}
	return b[:len(b)-n], nil
}

func bytesRepeat(b byte, n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = b
	}
	return out
}
