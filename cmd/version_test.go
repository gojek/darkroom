package cmd

import (
	"bytes"
	"github.com/gojek/darkroom/internal/version"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type VersionCmdSuite struct {
	suite.Suite
	rootCmd *cobra.Command
	buf     *bytes.Buffer
}

func TestVersionCmd(t *testing.T) {
	suite.Run(t, new(VersionCmdSuite))
}

func (s *VersionCmdSuite) SetupSuite() {
	version.Build = version.BuildInfo{
		Version:   "0.1.0",
		GitTag:    "v0.1.0",
		GitCommit: "c910e75b573b48961c7dcc1dd1063a543164d963",
		BuildDate: "2020-03-03T10:59:06Z",
	}
}

func (s *VersionCmdSuite) SetupTest() {
	s.rootCmd = &cobra.Command{
		Use: "app",
	}
	s.buf = &bytes.Buffer{}
	s.rootCmd.SetOut(s.buf)
}

func (s *VersionCmdSuite) TestVersionOutput() {
	// given
	s.rootCmd.SetArgs([]string{"version"})
	// when
	err := s.rootCmd.Execute()
	// then
	s.NoError(err)
	// and
	s.Equal(strings.TrimSpace(`0.1.0`), strings.TrimSpace(s.buf.String()))
}
