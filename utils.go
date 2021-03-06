package tcpx

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"io"
	"net"
	"strings"
)

type H map[string]interface{}

func Debug(src interface{}) string {
	buf, e := json.MarshalIndent(src, "  ", "  ")
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
	}
	return string(buf)
}

// Whether s in arr
// Support %%
func In(s string, arr []string) bool {
	for _, v := range arr {
		if strings.Contains(v, "%") {
			if strings.HasPrefix(v, "%") && strings.HasSuffix(v, "%") {
				if strings.Contains(s, string(v[1:len(v)-1])) {
					return true
				}
			} else if strings.HasPrefix(v, "%") {
				if strings.HasSuffix(s, string(v[1:])) {
					return true
				}
			} else if strings.HasSuffix(v, "%") {
				if strings.HasPrefix(s, string(v[:len(v)-1])) {
					return true
				}
			}
		} else {
			if v == s {
				return true
			}
		}
	}
	return false
}

// Defer eliminates all panic cases and handle panic reason by handlePanicError
func Defer(f func(), handlePanicError ...func(interface{})) {
	defer func() {
		if e := recover(); e != nil {
			for _, handler := range handlePanicError {
				handler(e)
			}
		}
	}()
	f()
}

// CloseChanel(func(){close(chan)})
func CloseChanel(f func()) {
	defer func() {
		if e := recover(); e != nil {
			// when close(chan) panic from 'close of closed chan' do nothing
		}
	}()
	f()
}
func MD5(rawMsg string) string {
	data := []byte(rawMsg)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has)
	return strings.ToUpper(md5str1)
}

// Write full buf
// In case buf is too big and conn can't write once.
//
//
/*
   if len(buf)>65535 {
       connLock.Lock()
       WriteConn(buf, conn)
       connLock.Unlock()
    } else {
       conn.Write(buf)
   }
*/
//
func WriteConn(buf []byte, conn net.Conn) error {
	var sum = 0
	for {
		n, e := conn.Write(buf)

		if e != nil {
			if e == io.EOF {
				return io.EOF
				break
			}
			return errorx.Wrap(e)
			break
		}
		sum += n
		if sum >= len(buf) {
			break
		}
	}
	return nil
}
