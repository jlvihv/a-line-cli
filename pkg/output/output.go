package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hamster-shared/a-line-cli/pkg/logger"
)

type Output struct {
	Name       string
	ID         int
	buffer     []string
	f          *os.File
	mu         sync.Mutex
	filename   string
	done       bool
	startIndex int
}

func New(name string, id int) *Output {
	o := &Output{
		Name:   name,
		ID:     id,
		buffer: make([]string, 16),
	}

	err := o.initFile()
	if err != nil {
		logger.Errorf("Failed to init output file, err: %s", err)
	}

	err = o.ToFile()
	if err != nil {
		logger.Errorf("Failed to write output to file, err: %s", err)
	}

	o.WriteLine("[Job] Started on " + time.Now().Local().Format("2006-01-02 15:04:05"))

	return o
}

func (o *Output) Destroy() {
	o.buffer = nil
}

func (o *Output) Done() {
	o.mu.Lock()
	o.done = true
	logger.Trace("Output done, flush all, close file")
	o.flush(o.buffer[o.startIndex:])
	o.flush([]string{"\n\n\n[Job] Finished on " + time.Now().Local().Format("2006-01-02 15:04:05")})
	o.f.Close()
	o.mu.Unlock()
}

func (o *Output) WriteLine(line string) {
	// 如果不是以换行符结尾，自动添加
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	o.buffer = append(o.buffer, line)
}

func (o *Output) WriteCommandLine(line string) {
	o.WriteLine("> " + line)
}

func (o *Output) Content() string {
	return strings.Join(o.buffer, "\n")
}

func (o *Output) NewStage(name string) {
	o.buffer = append(o.buffer, "\n\n\n[Pipeline] Stage: "+name+"\n")
}

func (o *Output) ToFile() error {
	endIndex := 0
	go func(endIndex int) {
		for {
			o.mu.Lock()
			if o.done {
				o.mu.Unlock()
				break
			}
			o.mu.Unlock()

			if len(o.buffer) <= endIndex {
				time.Sleep(1 * time.Second)
				continue
			}
			endIndex = len(o.buffer)
			err := o.flush(o.buffer[o.startIndex:endIndex])
			if err != nil {
				logger.Error(err)
			}
			o.startIndex = endIndex
			time.Sleep(1 * time.Second)
		}
	}(endIndex)

	return nil
}

// 刷入文件
func (o *Output) flush(arr []string) error {
	for _, line := range arr {
		if _, err := o.f.WriteString(line); err != nil {
			logger.Error(err)
			return err
		}
	}
	return nil
}

func (o *Output) initFile() error {
	o.mu.Lock()
	if o.f != nil {
		o.mu.Unlock()
		return nil
	}

	if o.filename == "" {
		o.filename = o.Name + "-" + fmt.Sprint(o.ID) + "-" + time.Now().Local().Format("2006-01-02-15:04:05") + ".log"
		o.filename = filepath.Join("log", o.filename)
	}

	f, err := os.OpenFile(o.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Errorf("Failed to create output log file %s, err: %s\n", o.filename, err)
		o.mu.Unlock()
		return err
	}
	o.f = f
	o.mu.Unlock()
	return nil
}

func (o *Output) Filename() string {
	return o.filename
}
