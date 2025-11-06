package solver

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
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
	return nil, nil
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
		return nil, err
	}

	boxPath := strings.TrimSpace(initOutput.String())
	sandBoxPath := filepath.Join(boxPath, "box")
	stdinFile := filepath.Join(sandBoxPath, "stdin.txt")
	stdoutFile := filepath.Join(sandBoxPath, "stdout.txt")
	stderrFile := filepath.Join(sandBoxPath, "stderr.txt")
	metaFile := filepath.Join(boxPath, "meta.txt")

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

	for _, t := range test.Tests {
		if err := os.WriteFile(stdinFile, t.Input, 0o644); err != nil {
			log.Error().Err(err).Msg("Failed to write stdin")
			return nil, err
		}

		runCmd := exec.Command(
			"isolate",
			"--box-id", "0",
			// "--cg",
			"--meta", metaFile,
			"--stdin", "stdin.txt",
			"--stdout", "stdout.txt",
			"--stderr", "stderr.txt",
			"--time", fmt.Sprintf("%d", 5),
			"--mem", fmt.Sprintf("%d", 4096*1024),
			"--processes=64",
			"--run", "--", "./solve",
		)
		var runStdout bytes.Buffer
		var runStderr bytes.Buffer
		runCmd.Stdout = &runStdout
		runCmd.Stderr = &runStderr

		if err := runCmd.Run(); err != nil {
			log.Error().Err(err).Str("stdout", runStdout.String()).Str("stderr", runStderr.String()).Msg("Failed to run isolate")
			// return err
		}

		stdoutData, _ := os.ReadFile(stdoutFile)
		stderrData, _ := os.ReadFile(stderrFile)
		metaData, _ := os.ReadFile(metaFile)

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
			result.ErrorMessage = utils.Ptr(string(stderrData))
			break
		}
		result.MemoryUsage = max(result.MemoryUsage, getMemoryUsage(metadata))

		if string(stdoutData) != string(t.Output) {
			result.Status = "WA"
			result.ErrorMessage = utils.Ptr("Wrong answer")
			break
		}
	}

	if result.ErrorMessage == nil {
		result.Status = "AC"
	}

	cleanupCmd := exec.Command(
		"isolate",
		"--box-id", "0",
		// "--cg",
		"--cleanup",
	)
	if err := cleanupCmd.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to run isolate")
		return nil, err
	}

	return result, nil
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
