package gutil

import "encoding/json"
import "crypto/md5"
import "encoding/hex"
import "crypto/sha1"
import b64 "encoding/base64"

/*
一些密文算法
*/

/*
MD5编码
*/
func Md5Encode(str string) string {
	m := md5.New()
	if _, err := m.Write([]byte(str)); CheckSucceed(err) {
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
	if _, err := sha.Write(bits); CheckSucceed(err) {
		return sha.Sum(nil)
	}
	return nil
}

//编码
func JsonEncode(data interface{}) []byte {
	if bits, err := json.Marshal(data); CheckSucceed(err) {
		return bits
	}
	return nil
}

/*
解码
*/
func JsonDecode(val []byte, data interface{}) {
	//jsonMap := make(map[string]string)
	if err := json.Unmarshal(val, &data); err != nil {
		//Json解析出错
	}
}
