package base

/*
一些密文算法
*/

import "encoding/json"
import "crypto/md5"
import "encoding/hex"
import "crypto/sha1"
import b64 "encoding/base64"

/*
MD5编码
*/
func Md5Encode(str string) string {
	m := md5.New()
	if _, err := m.Write([]byte(str)); err == nil {
		mstr := m.Sum(nil)
		return hex.EncodeToString(mstr)
	}
	return ""
}

/*
Base64编码
*/
func Base64Encode(bits []byte) string {
	return b64.StdEncoding.EncodeToString(bits)
}

/*
Sha1编码
*/
func Sha1Encode(bits []byte) []byte {
	sha := sha1.New()
	if _, err := sha.Write(bits); err == nil {
		return sha.Sum(nil)
	}
	return nil
}

//编码(失败返回空)
func JsonEncode(data interface{}) []byte {
	if bits, err := json.Marshal(data); err == nil {
		return bits
	}
	return []byte{}
}

/*
解码
*/
func JsonDecode(val []byte, data interface{}) error {
	//支持: make(map[string]interface{})
	return json.Unmarshal(val, &data)
}
