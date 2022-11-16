package output

import (
	"fmt"
	"testing"

	"github.com/hamster-shared/a-line-cli/pkg/logger"
	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	logger.Init().ToStdoutAndFile().SetLevel(logrus.TraceLevel)
	testOutput := New("test", 10085)

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

	fmt.Println("文件写入到", testOutput.filename)
}
