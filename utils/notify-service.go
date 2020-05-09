// +build notify

/**
 * wxpay-notify implementation.
 */
package utils

import (
	"gopkg.in/natefinch/lumberjack.v2"
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
	"go-wxpay-gateway/conf"
)

const (
	_LINE_PATTERN = "([^\t]*)\t(\\d+)" // YYYY-MM-DD hh:mm:ss\t<bytes_of_item>\n
	_MAX_HEAD_SIZE = 100
)

type _LogItem struct {
	AppName string `json:"app_name"`
	CbUrl   string `json:"cb_url"`
	Params  map[string]interface{} `json:"params"`
}

type _NotifyItem struct {
	timestamp string
	appId     string
	logItem  *_LogItem
}

var (
	_notifyLog *log.Logger
	_re *regexp.Regexp
	_items chan *_NotifyItem
)

func InitNotifyLog(logFile string) error {
	_notifyLogFile := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   false, // disabled by default
	}
	_notifyLog = log.New(_notifyLogFile, "", log.LstdFlags)
	return nil
}

func init() {
	_re, _ = regexp.Compile(_LINE_PATTERN)
}

type _M  map[string]interface{}

func (m _M) getString(val *string, name string) error {
	if v, ok := m[name]; !ok {
		return fmt.Errorf("param %s not found", name)
	} else {
		switch v.(type) {
		case string:
			*val = v.(string)
			return nil
		default:
			return fmt.Errorf("string expected for param %s", name)
		}
	}
}

func parseAppId(params map[string]interface{}) (appId string, err error) {
	if err = _M(params).getString(&appId, "app_id"); err != nil {
		return
	}
	delete(params, "app_id")
	return
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
		var logItem _LogItem
		if err = json.Unmarshal(item, &logItem); err != nil {
			_notifyLog.Printf("[warn] failed to parse item: %v\n", err)
			break
		}

		appId, err := parseAppId(logItem.Params)
		if err != nil {
			_notifyLog.Printf("[warn] failed to parse item: %v\n", err)
		}
		// notify it
		_items <- &_NotifyItem{timestamp, appId, &logItem}
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
	totalBytes := int(itemLen)
	item = make([]byte, totalBytes)
	bytesRead := 0
	for bytesRead < totalBytes {
		bl, err := r.Read(item[bytesRead:])
		if err != nil {
			return "", nil, fmt.Errorf("failed to read %d bytes: %v", totalBytes, err)
		}
		bytesRead += bl
	}
	return
}

func notifyService(threadNo int, item *_NotifyItem) error {
	i := 0 // do{} while
	for {
		logItem := item.logItem
		status, content, _, err := wget.PostJson(logItem.CbUrl, "POST", logItem.Params, nil)
		if err != nil {
			i++
			if i >= conf.NotifyConf.RetryCount {
				return err
			}
			time.Sleep(10*time.Second)
			continue
		}
		_notifyLog.Printf("[notity-response] #%d [%s %s %s %s]: status: %d, content: %s\n",
			threadNo, item.timestamp, item.appId, logItem.AppName, logItem.CbUrl,
			status, content,
		)
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
				logItem := item.logItem
				if err := notifyService(i, item); err != nil {
					_notifyLog.Printf("[notify-failed] #%d [%s %s %s %s]: %v\n", i, item.timestamp, item.appId, logItem.AppName, logItem.CbUrl, err)
				} else {
					_notifyLog.Printf("[notify-ok] #%d [%s %s %s %s]\n", i, item.timestamp, item.appId, logItem.AppName, logItem.CbUrl)
				}
			}
		}(i)
	}
}
