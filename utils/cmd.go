package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCommand executes a shell command and prints output
// 执行Shell命令并打印输出
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommandOutput executes a shell command and returns output
// 执行Shell命令并返回输出
func RunCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// ReadInput reads a line from stdin
// 读取用户输入
func ReadInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// CheckCommandExists checks if a command exists in PATH
// 检查命令是否存在
func CheckCommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
