package output

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hamster-shared/a-line-cli/pkg/logger"
)

type Output struct {
	Name           string
	ID             int
	StartTime      time.Time
	buffer         []string
	f              *os.File
	mu             sync.Mutex
	filename       string
	done           bool
	fileCursor     int
	bufferCursor   int
	stageStartTime map[string]time.Time
}

// New 新建一个 Output 对象，会自动初始化文件，以及定时将内容写入文件
func New(name string, id int) *Output {
	o := &Output{
		Name:      name,
		ID:        id,
		StartTime: time.Now().Local(),
		buffer:    make([]string, 0, 16),
	}

	err := o.initFile()
	if err != nil {
		logger.Errorf("Failed to init output file, err: %s", err)
		return o
	}

	o.timedWriteFile()

	o.WriteLine("[Job] Started on " + o.StartTime.Format("2006-01-02 15:04:05"))

	return o
}

// Duration 返回持续时间
func (o *Output) Duration() time.Duration {
	return time.Since(o.StartTime)
}

// StageDuration 返回某个 Stage 的持续时间
func (o *Output) StageDuration(name string) time.Duration {
	stageStartTime, ok := o.stageStartTime[name]
	if !ok {
		return 0
	}
	if stageStartTime.IsZero() {
		return 0
	}
	return time.Since(o.stageStartTime[name])
}

// Done 标记输出已完成，会将缓存中的内容刷入文件，然后关闭文件
func (o *Output) Done() {
	o.mu.Lock()
	o.done = true
	logger.Trace("output done, flush all, close file")
	o.flush(o.buffer[o.fileCursor:])
	o.flush([]string{"\n\n\n[Job] Finished on " + time.Now().Local().Format("2006-01-02 15:04:05")})
	o.f.Close()
	o.mu.Unlock()
}

// WriteLine 将一行普通内容写入输出
func (o *Output) WriteLine(line string) {
	// 如果不是以换行符结尾，自动添加
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	o.buffer = append(o.buffer, line)
}

// WriteCommandLine 将一行命令行内容写入输出，其实就是在前面加上了一个 "> "
func (o *Output) WriteCommandLine(line string) {
	o.WriteLine("> " + line)
}

// Content 总是返回从起始到现在的所有内容
func (o *Output) Content() string {
	o.bufferCursor = len(o.buffer)
	return strings.Join(o.buffer[:o.bufferCursor], "")
}

// NewContent 总是返回自上次读取后新出现的内容
func (o *Output) NewContent() string {
	if o.bufferCursor >= len(o.buffer) {
		return ""
	}
	endIndex := len(o.buffer)
	result := strings.Join(o.buffer[o.bufferCursor:endIndex], "")
	o.bufferCursor = endIndex
	return result
}

// NewStage 会写入以 [Pipeline] Stage: 开头的一行，表示一个新的 Stage 开始
func (o *Output) NewStage(name string) {
	o.WriteLine("\n")
	o.WriteLine("\n")
	o.WriteLine("\n")
	o.WriteLine("[Pipeline] Stage: " + name)
	o.stageStartTime[name] = time.Now().Local()
}

// 在一个协程中定时刷入文件
func (o *Output) timedWriteFile() {
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
			err := o.flush(o.buffer[o.fileCursor:endIndex])
			if err != nil {
				logger.Error(err)
			}
			o.fileCursor = endIndex
			time.Sleep(1 * time.Second)
		}
	}(endIndex)
}

// 刷入文件
func (o *Output) flush(arr []string) error {
	if o.f == nil {
		return nil
	}
	for _, line := range arr {
		if _, err := o.f.WriteString(line); err != nil {
			logger.Error(err)
			return err
		}
	}
	return nil
}

// 初始化文件
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

// Filename 返回文件名
func (o *Output) Filename() string {
	if o.f == nil {
		return ""
	}
	return o.f.Name()
}

type StageOutput struct {
	Name    string
	Content string
}

// StageOutputList 返回存储了 Stage 输出的列表
func (o *Output) StageOutputList() []StageOutput {
	return parseLogLines(o.buffer[:])
}

// ParseLogFile 解析日志文件，返回存储了 Stage 输出的列表
func ParseLogFile(filename string) ([]StageOutput, error) {
	lines, err := readFileLines(filename)
	if err != nil {
		return nil, err
	}
	result := parseLogLines(lines)
	return result, nil
}

// 读取文件中的行
func readFileLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	return lines, nil
}

func parseLogLines(lines []string) []StageOutput {
	var stageName = "unknown"
	var stageNameList []string

	// 先遍历到 map 里，由于 map 是无序的，所以需要一个数组来记录顺序
	var stageOutputMap = make(map[string][]string)
	for _, line := range lines {
		if strings.HasPrefix(line, "[Job]") || line == "\n" || line == "" {
			continue
		}
		if strings.HasPrefix(line, "[Pipeline] Stage: ") {
			stageName = strings.TrimPrefix(line, "[Pipeline] Stage: ")
			stageOutputMap[stageName] = make([]string, 0)
			stageNameList = append(stageNameList, stageName)
		}
		stageOutputMap[stageName] = append(stageOutputMap[stageName], line)
	}

	var stageOutputList []StageOutput
	for k, v := range stageOutputMap {
		for i := range stageNameList {
			if stageNameList[i] == k {
				stageOutput := StageOutput{
					Name:    k,
					Content: strings.Join(v, "\n"),
				}
				stageOutputList = append(stageOutputList, stageOutput)
			}
		}
	}

	return stageOutputList
}
