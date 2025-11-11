package solver

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	attemptDTO "problum/internal/attempt/service/dto"
	"problum/internal/solver/dto"
	templateDTO "problum/internal/template/service/dto"
	testDTO "problum/internal/test/service/dto"
	"problum/internal/utils"

	"github.com/bytedance/sonic"
	"github.com/flosch/pongo2/v6"
	"github.com/rs/zerolog/log"
	"golang.org/x/tools/imports"
)

var (
	errWrongAnswer    = errors.New("wrong answer")
	errNotNilExitCode = errors.New("not nil exit code")
)

type TestService interface {
	GetByProblemID(context.Context, int) (*testDTO.Test, error)
}

type TemplateService interface {
	GetByProblemIDAndLanguage(context.Context, int, string) (*templateDTO.Template, error)
}

type Solver struct {
	testSvc     TestService
	templateSvc TemplateService
}

type runIsolateConfig struct {
	StdinFile  string
	StdoutFile string
	StderrFile string
	MetaFile   string
	BoxFile    string
	Time       string
	WallTime   string
	Mem        string
	Processes  string
	RunCommand []string
}

func New(testSvc TestService, templateSvc TemplateService) *Solver {
	return &Solver{
		testSvc:     testSvc,
		templateSvc: templateSvc,
	}
}

func (s *Solver) Solve(ctx context.Context, attempt *attemptDTO.Attempt) (*dto.Result, error) {
	test, err := s.testSvc.GetByProblemID(ctx, attempt.ProblemID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get test for problem")
		return nil, fmt.Errorf("failed to get test for problem: %w", err)
	}

	template, err := s.templateSvc.GetByProblemIDAndLanguage(ctx, attempt.ProblemID, attempt.Language)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get template for problem")
		return nil, fmt.Errorf("failed to get template for problem: %w", err)
	}

	metadata, err := parseTemplateMetadata(template.Metadata)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse template metadata")
		return nil, fmt.Errorf("failed to parse template metadata: %w", err)
	}
	metadata["code"] = attempt.Code
	log.Info().Interface("metadata", metadata).Msg("metadata")

	switch attempt.Language {
	case "python":
		return s.solvePython(test, metadata)
	case "go":
		return s.solveGolang(test, metadata)
	default:
		log.Error().Err(err).Msg("Failed to get test for problem")
		return nil, fmt.Errorf("unsupported language")
	}
}

func (s *Solver) solvePython(test *testDTO.Test, md map[string]any) (*dto.Result, error) {
	result := &dto.Result{}

	if err := s.renderTemplate("code.py.j2", md); err != nil {
		return nil, err
	}

	path, err := initIsolate()
	if err != nil {
		log.Error().Err(err).Msg("Failed to init isolate")
		return nil, fmt.Errorf("failed to init isolate: %w", err)
	}

	if err := runTests(path, "python", test.Tests, result); err != nil {
		log.Error().Err(err).Msg("Failed to run tests")
		return nil, fmt.Errorf("failed to run tests: %w", err)
	}

	if err := cleanupIsolate(); err != nil {
		log.Error().Err(err).Msg("Failed to cleanup isolate")
		return nil, fmt.Errorf("failed to cleanup isolate: %w", err)
	}

	return result, nil
}

func (s *Solver) solveGolang(test *testDTO.Test, md map[string]any) (*dto.Result, error) {
	result := &dto.Result{}

	if err := s.renderTemplate("code.go.j2", md); err != nil {
		return nil, err
	}

	if err := s.renderTemplate("harness.go.j2", md); err != nil {
		return nil, err
	}

	if errorMsg, err := s.compileGolang(); err != nil {
		result.Status = "CE"
		if errorMsg != nil {
			result.ErrorMessage = errorMsg
		} else {
			result.ErrorMessage = utils.Ptr("Compile error")
		}

		return result, nil
	}

	path, err := initIsolate()
	if err != nil {
		log.Error().Err(err).Msg("Failed to init isolate")
		return nil, fmt.Errorf("failed to init isolate: %w", err)
	}

	if err := runTests(path, "go", test.Tests, result); err != nil {
		log.Error().Err(err).Msg("Failed to run tests")
		return nil, fmt.Errorf("failed to run tests: %w", err)
	}

	if err := cleanupIsolate(); err != nil {
		log.Error().Err(err).Msg("Failed to cleanup isolate")
		return nil, fmt.Errorf("failed to cleanup isolate: %w", err)
	}

	return result, nil
}

func initIsolate() (string, error) {
	initCmd := exec.Command(
		"isolate",
		"--box-id", "0",
		// "--cg",
		"--init",
	)
	var initOutput bytes.Buffer
	var initStderr bytes.Buffer
	initCmd.Stdout = &initOutput
	initCmd.Stderr = &initStderr

	if err := initCmd.Run(); err != nil {
		log.Error().Str("stdout", initOutput.String()).Str("stderr", initStderr.String()).Err(err).Msg("Failed to init isolate")
		return "", err
	}

	return initOutput.String(), nil
}

func cleanupIsolate() error {
	cleanupCmd := exec.Command(
		"isolate",
		"--box-id", "0",
		// "--cg",
		"--cleanup",
	)
	if err := cleanupCmd.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to cleanup isolate")
		return fmt.Errorf("failed to cleanup isolate: %w", err)
	}

	return nil
}

func getIsolateConfig(path, language string) (*runIsolateConfig, error) {
	cfg := &runIsolateConfig{}

	boxPath := strings.TrimSpace(path)
	sandBoxPath := filepath.Join(boxPath, "box")

	cfg.StdinFile = filepath.Join(sandBoxPath, "stdin.txt")
	cfg.StdoutFile = filepath.Join(sandBoxPath, "stdout.txt")
	cfg.StderrFile = filepath.Join(sandBoxPath, "stderr.txt")
	cfg.MetaFile = filepath.Join(boxPath, "meta.txt")

	switch language {
	case "go":
		goFile, err := os.ReadFile("solve")
		if err != nil {
			log.Error().Err(err).Msg("Failed to read go file")
			return nil, err
		}

		goBoxFile := filepath.Join(sandBoxPath, "solve")
		if err := os.WriteFile(goBoxFile, goFile, 0o755); err != nil {
			log.Error().Err(err).Msg("Failed to write go file")
			return nil, err
		}

		cfg.BoxFile = goBoxFile
		cfg.Time = fmt.Sprintf("%d", 5)
		cfg.WallTime = fmt.Sprintf("%d", 5)
		cfg.Mem = fmt.Sprintf("%d", 4096*1024)
		cfg.Processes = fmt.Sprintf("--processes=%d", 64)
		cfg.RunCommand = []string{"--run", "--", "./solve"}
	case "python":
		pythonFile, err := os.ReadFile("code.py")
		if err != nil {
			log.Error().Err(err).Msg("Failed to read python file")
			return nil, err
		}

		pythonBoxFile := filepath.Join(sandBoxPath, "code.py")
		if err := os.WriteFile(pythonBoxFile, pythonFile, 0o755); err != nil {
			log.Error().Err(err).Msg("Failed to write python file")
			return nil, err
		}

		cfg.BoxFile = pythonBoxFile
		cfg.Time = fmt.Sprintf("%d", 5)
		cfg.WallTime = fmt.Sprintf("%d", 5)
		cfg.Mem = fmt.Sprintf("%d", 4096*1024)
		cfg.RunCommand = []string{"--run", "--", "/usr/bin/python3", "./code.py"}
	default:
		return nil, fmt.Errorf("unsupported language")
	}

	return cfg, nil
}

func runIsolate(cfg *runIsolateConfig, test *testDTO.TestCase, result *dto.Result) error {
	if err := os.WriteFile(cfg.StdinFile, test.Input, 0o644); err != nil {
		log.Error().Err(err).Msg("Failed to write stdin")
		return err
	}

	args := []string{
		"--box-id", "0",
		// "--cg",
		"--meta", cfg.MetaFile,
		"--stdin", "stdin.txt",
		"--stdout", "stdout.txt",
		"--stderr", "stderr.txt",
		"--time", cfg.Time,
		"--wall-time", cfg.WallTime,
		"--mem", cfg.Mem,
	}
	if cfg.Processes != "" {
		args = append(args, cfg.Processes)
	}
	args = append(args, cfg.RunCommand...)
	runCmd := exec.Command("isolate", args...)

	var runStdout bytes.Buffer
	var runStderr bytes.Buffer
	runCmd.Stdout = &runStdout
	runCmd.Stderr = &runStderr

	if err := runCmd.Run(); err != nil {
		log.Error().Err(err).Str("stdout", runStdout.String()).Str("stderr", runStderr.String()).Msg("Failed to run isolate")
		// return err
	}

	stdoutData, _ := os.ReadFile(cfg.StdoutFile)
	stderrData, _ := os.ReadFile(cfg.StderrFile)
	metaData, _ := os.ReadFile(cfg.MetaFile)

	metadata := parseMetadata(metaData)

	exitCode := 1
	if exitCodeStr, ok := metadata["exitcode"]; ok {
		if code, err := strconv.Atoi(exitCodeStr); err == nil {
			exitCode = code
		}
	}

	log.Info().Str("stdout", string(stdoutData)).Msg("stdout")
	log.Info().Str("stderr", string(stderrData)).Msg("stderr")
	log.Info().Int("exit_code", exitCode).Msg("exit_code")
	log.Info().Interface("metadata", metadata).Msg("metadata")

	result.Duration = max(result.Duration, getDuration(metadata))
	if exitCode != 0 {
		if status, ok := metadata["status"]; ok {
			result.Status = status
		} else {
			result.Status = "RE"
		}

		if string(stderrData) != "" {
			result.ErrorMessage = utils.Ptr(string(stderrData))
		} else if message, ok := metadata["message"]; ok {
			result.ErrorMessage = utils.Ptr(message)
		} else {
			result.ErrorMessage = utils.Ptr("Runtime error")
		}

		return errNotNilExitCode
	}
	result.MemoryUsage = max(result.MemoryUsage, getMemoryUsage(metadata))

	if string(stdoutData) != string(test.Output) {
		result.Status = "WA"
		result.ErrorMessage = utils.Ptr("Wrong answer")
		return errWrongAnswer
	}

	return nil
}

func runTests(path, language string, tests []testDTO.TestCase, result *dto.Result) error {
	cfg, err := getIsolateConfig(path, language)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get isolate config")
		return fmt.Errorf("failed to get isolate config: %w", err)
	}
	result.Status = "AC"

	for _, t := range tests {
		if err := runIsolate(cfg, &t, result); err != nil {
			log.Error().Err(err).Msg("Failed to run test")
			break
		}
	}

	return nil
}

func (s *Solver) renderTemplate(filename string, context map[string]any) error {
	template, err := pongo2.FromFile(filename)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get template from file")
		return fmt.Errorf("failed to get template from file: %w", err)
	}

	ctx := pongo2.Context{
		"code":          context["code"],
		"function_name": context["function_name"],
		"parameters":    context["parameters"],
	}

	out, err := template.Execute(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute template")
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(strings.TrimSuffix(filename, ".j2"), []byte(out), 0o644); err != nil {
		log.Error().Err(err).Msg("Failed to write file")
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *Solver) compileGolang() (*string, error) {
	codeGo, _ := os.ReadFile("code.go")
	harnessGo, _ := os.ReadFile("harness.go")

	codeGoFormatted, err := imports.Process("code.go", codeGo, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to format code.go")
		return utils.Ptr(err.Error()), fmt.Errorf("failed to format code.go: %w", err)
	}

	harnessGoFormatted, err := imports.Process("harness.go", harnessGo, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to format harness.go")
		return nil, fmt.Errorf("failed to format harness.go: %w", err)
	}

	if err := os.WriteFile("code.go", codeGoFormatted, 0o644); err != nil {
		log.Error().Err(err).Msg("Failed to write formatted code.go")
		return nil, fmt.Errorf("failed to write formatted code.go: %w", err)
	}

	if err := os.WriteFile("harness.go", harnessGoFormatted, 0o644); err != nil {
		log.Error().Err(err).Msg("Failed to write formatted harness.go")
		return nil, fmt.Errorf("failed to write formatted harness.go: %w", err)
	}

	runCmd := exec.Command(
		"go", "build", "-o", "solve", "code.go", "harness.go",
	)

	var runStdout bytes.Buffer
	var runStderr bytes.Buffer
	runCmd.Stdout = &runStdout
	runCmd.Stderr = &runStderr

	if err := runCmd.Run(); err != nil {
		log.Error().Str("stdout", runStdout.String()).Str("stderr", runStderr.String()).Err(err).Msg("Failed to compile golang")
		return nil, err
	}

	return nil, nil
}

func parseMetadata(b []byte) map[string]string {
	metadata := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(b)))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			metadata[key] = value
		}
	}

	return metadata
}

func parseTemplateMetadata(metadata json.RawMessage) (map[string]any, error) {
	md := make(map[string]any)
	if err := sonic.Unmarshal(metadata, &md); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal metadata: %w", err)
	}

	return md, nil
}

func getDuration(metadata map[string]string) time.Duration {
	timeStr, ok1 := metadata["time"]
	timeWallStr, ok2 := metadata["time-wall"]

	if !ok1 && !ok2 {
		return 0
	}

	timeSeconds, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse time seconds")
	}

	timeWallSeconds, err := strconv.ParseFloat(timeWallStr, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse time-wall seconds")
	}

	ans := time.Duration(0)
	if ok1 && ok2 {
		ans = min(time.Duration(timeSeconds)*time.Nanosecond, time.Duration(timeWallSeconds)*time.Nanosecond)
	} else if ok1 {
		ans = time.Duration(timeSeconds) * time.Nanosecond
	} else {
		ans = time.Duration(timeWallSeconds) * time.Nanosecond
	}

	return ans
}

func getMemoryUsage(metadata map[string]string) int64 {
	memoryStr, ok := metadata["max-rss"]
	if !ok {
		return 0
	}

	memory, err := strconv.ParseInt(memoryStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse int from max-rss")
		return 0
	}

	// to bytes
	return memory * 1024
}
