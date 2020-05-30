package oauth

import (
	"encoding/binary"
	"encoding/hex"
	"hash/adler32"
	"time"
	"fmt"
)

// -- 产生随机串 32b时间戳+32b checksum--
func GenerateState() (string, uint32) {
	now := uint32(time.Now().Unix())
	expireTime := now + 600 // 10 min
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, expireTime)
	checksum := adler32.Checksum(b[:4])
	binary.LittleEndian.PutUint32(b[4:], checksum)
	return hex.EncodeToString(b), now
}

// -- 验证随机串 ---
func VerifyState(state string) (err error) {
	if len(state) != 16 {
		err = fmt.Errorf("bad size")
		return
	}
	var b []byte
	b, err = hex.DecodeString(state)
	if err != nil {
		return
	}
	checksum := adler32.Checksum(b[:4])
	rChecksum := binary.LittleEndian.Uint32(b[4:])
	if checksum != rChecksum {
		err = fmt.Errorf("bad checksum")
		return
	}
	now := time.Now().Unix()
	ts := binary.LittleEndian.Uint32(b)
	if now > int64(ts){
		err = fmt.Errorf("state expired")
	}
	return
}

