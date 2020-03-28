package builtins_test

import (
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/dalloriam/bish/bish/builtins"
	"github.com/dalloriam/bish/bish/state"
)

type mockExecutor struct{}

func (*mockExecutor) Exec(s string, in io.Reader, stdout, stderr io.Writer) error {
	return nil
}

type validationFunction func(ctx *state.State, args []string) bool

func Test_Registry(t *testing.T) {

	type testCase struct {
		name string
		args []string

		outputPattern string
		expectedApply bool
		wantErr       bool

		validate validationFunction
	}
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.RemoveAll(tempDir)

	cases := []testCase{
		// Alias tests.
		{"alias/single-argument", []string{"alias", "g", "gst"}, "", true, false, validateAlias},
		{"alias/multi-argument", []string{"alias", "g", "gst"}, "", true, false, validateAlias},
		{"alias/missing-body", []string{"alias", "g"}, "", true, true, validateAlias},
		{"alias/no-argument", []string{"alias"}, "", true, true, validateAlias},

		// Chdir tests.
		// TODO: Tests for relative chdir.
		{"cd/default", []string{"cd"}, "", true, false, validateCd},
		{"cd/dir", []string{"cd", tempDir}, "", true, false, validateCd},
		{"cd/too-many-arguments", []string{"cd", "bing", "bong"}, "", true, true, validateCd},

		// Env tests.
		{"env/set", []string{"set", "mykey", "myval"}, "", true, false, validateEnv},
		{"env/key-only", []string{"set", "mykey"}, "", true, true, validateEnv},
		{"env/no-argument", []string{"set"}, "", true, true, validateEnv},
		{"env/too-many-argument", []string{"set", "key", "val", "extra"}, "", true, true, validateEnv},

		// Prompt test.
		{"prompt/set", []string{"prompt", "myprompt"}, "", true, false, validatePrompt},
		{"prompt/no-argument", []string{"prompt"}, "", true, true, validatePrompt},
		{"prompt/too-many-arguments", []string{"prompt", "aksjdasd", "aksdjaks"}, "", true, true, validatePrompt},

		// Version test.
		{"version/basic", []string{"version"}, "(?s)BiSH.*Version:.*Git Commit:", true, false, nil},
		{"version/ignore-args", []string{"version", "some", "args"}, "(?s)BiSH.*Version:.*Git Commit:", true, false, nil},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			executor := &mockExecutor{}
			ctx := state.New(executor.Exec)
			output, didApply, err := builtins.TryBuiltIns(ctx, tCase.args)

			if (err != nil) != tCase.wantErr {
				t.Errorf("expected error: %v, got error: %v", tCase.wantErr, err)
			}
			if err != nil {
				return
			}

			if didApply != tCase.expectedApply {
				t.Errorf("expected apply=%v, got apply=%v", tCase.expectedApply, didApply)
				return
			}

			if tCase.outputPattern != "" {
				// Validate Output.
				pattern := regexp.MustCompile(tCase.outputPattern)
				if !pattern.MatchString(output) {
					t.Error("output does not match provided pattern")
					return
				}
			} else {
				// Output must be empty.
				if output != "" {
					t.Error("output is not empty")
					return
				}
			}

			if didApply && tCase.validate != nil {
				if ok := tCase.validate(ctx, tCase.args); !ok {
					t.Errorf("validation function failed")
					return
				}
			}

		})
	}
}
