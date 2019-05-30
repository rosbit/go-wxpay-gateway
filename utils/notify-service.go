// +build notify

/**
 * wxpay-notify implementation.
 */
package utils

import (
	"github.com/rosbit/go-wget"
	"os"
	"io"
	"fmt"
	"bufio"
	"regexp"
	"encoding/json"
	"strconv"
	"time"
	"log"
	"github.com/rosbit/go-wxpay-gateway/wx-pay-api"
	"github.com/rosbit/go-wxpay-gateway/conf"
)

const (
	_LINE_PATTERN = "([^\t]*)\t(\\d+)" // YYYY-MM-DD hh:mm:ss\t<bytes_of_item>\n
	_MAX_HEAD_SIZE = 100
)

type _NotifyItem struct {
	timestamp string
	params   *wxpay.PayNotifyParams
}

var (
	_notifyLog *log.Logger
	_re *regexp.Regexp
	_items chan *_NotifyItem
)

func InitNotifyLog(logFile string) error {
	_notifyLogFile, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_notifyLog = log.New(_notifyLogFile, "", log.LstdFlags)
	return nil
}

func init() {
	_re, _ = regexp.Compile(_LINE_PATTERN)
}

// read the content in file `fileName`, notify items one by one
func Notify(fileName string) {
	fp, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer fp.Close()

	r := bufio.NewReader(fp)
	for {
		// read one item, which is a JSON object
		timestamp, item, err := readOneItem(r)
		if err != nil {
			if err != io.EOF {
				_notifyLog.Printf("[warn] %v\n", err)
			}
			break
		}
		_notifyLog.Printf("[info] item read: %s %s ...", timestamp, string(head(item)))

		// parse the JSON to notification parameters
		var params wxpay.PayNotifyParams
		if err = json.Unmarshal(item, &params); err != nil {
			_notifyLog.Printf("[warn] failed to parse item: %v\n", err)
			break
		}

		// notify it
		_items <- &_NotifyItem{timestamp, &params}
	}
}

func head(item []byte) []byte {
	if len(item) < _MAX_HEAD_SIZE {
		return item
	}
	return item[:_MAX_HEAD_SIZE]
}

func readOneItem(r *bufio.Reader) (timestamp string, item []byte, err error) {
	// YYYY-MM-DD hh:mm:ss\t<bytes_of_item>\n
	line, err := r.ReadString('\n')
	if err != nil {
		return "", nil, err
	}
	m := _re.FindStringSubmatch(line)
	if m == nil {
		return "", nil, fmt.Errorf("header line bad format: %s", line)
	}
	timestamp = m[1]
	itemLen, err := strconv.ParseInt(m[2], 10, 32)
	if err != nil {
		return "", nil, fmt.Errorf("parse itemLen(%s) failed: %v", m[2], err)
	}

	// read item
	item = make([]byte, int(itemLen))
	bl, err := r.Read(item)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read %d bytes: %v", itemLen, err)
	}
	if bl != int(itemLen) {
		return "", nil, fmt.Errorf("try to read %d bytes, only %d bytes read", itemLen, bl)
	}
	return
}

func notifyService(params *wxpay.PayNotifyParams) error {
	cbUrl, realParams := params.CbUrl, &params.IPayNotifyParams

	i := 0 // do{} while
	for {
		status, content, _, err := wget.PostJson(cbUrl, "POST", realParams, nil)
		if err != nil {
			i++
			if i >= conf.NotifyConf.RetryCount {
				return err
			}
			time.Sleep(10*time.Second)
			continue
		}
		_notifyLog.Printf("[notity-response] status: %d, content: %s\n", status, content)
		return nil
	}
}

func StartNotifyThreads() {
	_items = make(chan *_NotifyItem, conf.NotifyConf.WorkerNum)
	for i:=0; i<conf.NotifyConf.WorkerNum; i++ {
		go func(i int) {
			for {
				item := <-_items
				if item == nil {
					break
				}
				params := item.params
				if err := notifyService(params); err != nil {
					_notifyLog.Printf("[notify-failed] #%d [%s, %s %s %s]: %v\n", i, item.timestamp, params.AppId, params.AppName, params.CbUrl, err)
				} else {
					_notifyLog.Printf("[notify-ok] #%d [%s %s %s %s]\n", i, item.timestamp, params.AppId, params.AppName, params.CbUrl)
				}
			}
		}(i)
	}
}
