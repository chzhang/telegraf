package snmp

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type mockedCommandResult struct {
	stdout    string
	stderr    string
	exitError bool
}

func mockExecCommand(arg0 string, args ...string) *exec.Cmd {
	args = append([]string{"-test.run=TestMockExecCommand", "--", arg0}, args...)
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stderr = os.Stderr // so the test output shows errors
	return cmd
}

// This is not a real test. This is just a way of mocking out commands.
//
// Idea based on https://github.com/golang/go/blob/7c31043/src/os/exec/exec_test.go#L568
func TestMockExecCommand(_ *testing.T) {
	var cmd []string
	for _, arg := range os.Args {
		if arg == "--" {
			cmd = []string{}
			continue
		}
		if cmd == nil {
			continue
		}
		cmd = append(cmd, arg)
	}
	if cmd == nil {
		return
	}

	cmd0 := strings.Join(cmd, "\000")
	mcr, ok := mockedCommandResults[cmd0]
	if !ok {
		cv := fmt.Sprintf("%#v", cmd)[8:] // trim `[]string` prefix
		//nolint:errcheck,revive
		fmt.Fprintf(os.Stderr, "Unmocked command. Please add the following to `mockedCommands` in snmp_mocks_generate.go, and then run `go generate`:\n\t%s,\n", cv)
		//nolint:revive // error code is important for this "test"
		os.Exit(1)
	}
	//nolint:errcheck,revive
	fmt.Printf("%s", mcr.stdout)
	//nolint:errcheck,revive
	fmt.Fprintf(os.Stderr, "%s", mcr.stderr)
	if mcr.exitError {
		//nolint:revive // error code is important for this "test"
		os.Exit(1)
	}
	//nolint:revive // error code is important for this "test"
	os.Exit(0)
}

func init() {
	execCommand = mockExecCommand
}

// BEGIN GO GENERATE CONTENT
var mockedCommandResults = map[string]mockedCommandResult{
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x00.1.0.0.0":                    {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x00.1.0.0.1.1":                  {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x00.1.0.0.1.2":                  {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x001.0.0.1.1":                   {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x00.1.0.0.0.1.1":                {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x00.1.0.0.0.1.1.0":              {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x00.1.0.0.0.1.5":                {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x00.1.2.3":                      {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x00.1.0.0.0.1.7":                {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00.iso.2.3":                                 {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00-m\x00all\x00.999":                        {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00TEST::server":                             {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00TEST::server.0":                           {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00TEST::testTable":                          {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00TEST::connections":                        {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00TEST::latency":                            {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00TEST::description":                        {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00TEST::hostname":                           {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00IF-MIB::ifPhysAddress.1":                  {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00BRIDGE-MIB::dot1dTpFdbAddress.1":          {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00-Ob\x00TCP-MIB::tcpConnectionLocalAddress.1":     {stdout: "", stderr: "", exitError: true},
	"snmptranslate\x00-Td\x00TEST::testTable.1":                               {stdout: "", stderr: "", exitError: true},
	"snmptable\x00-Ch\x00-Cl\x00-c\x00public\x00127.0.0.1\x00TEST::testTable": {stdout: "", stderr: "", exitError: true},
}
