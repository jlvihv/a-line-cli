package output

import (
	"fmt"
	"testing"

	"github.com/hamster-shared/a-line-cli/pkg/logger"
	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	logger.Init().ToStdoutAndFile().SetLevel(logrus.TraceLevel)
	testOutput := New("test", 10001)

	testOutput.NewStage("第一阶段")
	testOutput.WriteLine("第一行")
	testOutput.WriteLine("第二行")
	testOutput.WriteLine("第三行")
	testOutput.WriteLine("第四行")
	testOutput.WriteLine("第五行")

	testOutput.NewStage("第二阶段")
	testOutput.WriteLine("第一行")
	testOutput.WriteLine("第二行")
	testOutput.WriteLine("第三行")
	testOutput.WriteLine("第四行")
	testOutput.WriteLine("第五行")

	testOutput.Done()

	fmt.Println("文件写入到", testOutput.Filename())
}

func TestParseLogFile(t *testing.T) {
	result, err := ParseLogFile("log/test-10086-2022-11-17-15:34:44.log")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
}

func TestContent(t *testing.T) {
	logger.Init().ToStdoutAndFile().SetLevel(logrus.TraceLevel)
	testOutput := New("test", 10085)

	testOutput.NewStage("第一阶段")
	testOutput.WriteLine("第一行")
	testOutput.WriteLine("第二行")
	testOutput.WriteLine("第三行")
	testOutput.WriteLine("第四行")
	testOutput.WriteLine("第五行")

	fmt.Println("读取所有内容：", testOutput.Content())

	testOutput.NewStage("第二阶段")
	testOutput.WriteLine("第一行")
	testOutput.WriteLine("第二行")
	testOutput.WriteLine("第三行")
	testOutput.WriteLine("第四行")
	testOutput.WriteLine("第五行")

	fmt.Println("读取新出现的内容：", testOutput.NewContent())

	if len(testOutput.NewContent()) != 0 {
		t.Error("new content length error")
	}

	testOutput.Done()
}

func TestStageOutputList(t *testing.T) {
	logger.Init().ToStdoutAndFile().SetLevel(logrus.TraceLevel)
	testOutput := New("test", 10000)

	testOutput.NewStage("第一阶段")
	testOutput.WriteLine("第一行")
	testOutput.WriteLine("第二行")
	testOutput.WriteLine("第三行")
	testOutput.WriteLine("第四行")
	testOutput.WriteLine("第五行")

	testOutput.NewStage("第二阶段")
	testOutput.WriteLine("第一行")
	testOutput.WriteLine("第二行")
	testOutput.WriteLine("第三行")
	testOutput.WriteLine("第四行")
	testOutput.WriteLine("第五行")

	if len(testOutput.StageOutputList()) != 2 {
		t.Error("stage output list length error")
	}

	testOutput.Done()
}
