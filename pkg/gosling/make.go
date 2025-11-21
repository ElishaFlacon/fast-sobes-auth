package gosling

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	noUsecase    bool
	noRepository bool
	noHandler    bool
	noApp        bool
	methods      []string
)

var makeCmd = &cobra.Command{
	Use:   "make service [name]",
	Short: "Generate a new service with handler, usecase, and repository",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runMake,
}

func init() {
	makeCmd.Flags().BoolVar(&noUsecase, "no-usecase", false, "Skip usecase generation")
	makeCmd.Flags().BoolVar(&noRepository, "no-repository", false, "Skip repository generation")
	makeCmd.Flags().BoolVar(&noHandler, "no-handler", false, "Skip handler generation")
	makeCmd.Flags().BoolVar(&noApp, "no-app", false, "Skip app provider update")
	makeCmd.Flags().StringArrayVar(&methods, "method", []string{}, "Method names to generate (can be used multiple times)")
}

func runMake(cmd *cobra.Command, args []string) error {
	serviceName := args[0]
	modulePath, err := getModulePath()
	if err != nil {
		return fmt.Errorf("failed to read module path: %w", err)
	}

	generator := NewGenerator(serviceName, modulePath, methods)

	if !noRepository {
		if err := generator.GenerateRepository(); err != nil {
			return fmt.Errorf("failed to generate repository: %w", err)
		}
		if err := generator.UpdateRepositoryInterface(); err != nil {
			return fmt.Errorf("failed to update repository interface: %w", err)
		}
		fmt.Printf("✓ Repository generated for service '%s'\n", serviceName)
	}

	if !noUsecase {
		if err := generator.GenerateUsecase(); err != nil {
			return fmt.Errorf("failed to generate usecase: %w", err)
		}
		if err := generator.UpdateUsecaseInterface(); err != nil {
			return fmt.Errorf("failed to update usecase interface: %w", err)
		}
		fmt.Printf("✓ Usecase generated for service '%s'\n", serviceName)
	}

	if !noHandler {
		if err := generator.GenerateHandler(); err != nil {
			return fmt.Errorf("failed to generate handler: %w", err)
		}
		fmt.Printf("✓ Handler generated for service '%s'\n", serviceName)
	}

	if !noApp {
		if err := generator.UpdateProvider(); err != nil {
			return fmt.Errorf("failed to update provider: %w", err)
		}
		fmt.Printf("✓ Provider updated for service '%s'\n", serviceName)
	}

	fmt.Printf("\n🎉 Service '%s' generated successfully!\n", serviceName)
	fmt.Println("\nNext steps:")
	fmt.Printf("1. Define your protobuf service and regenerate code\n")
	fmt.Printf("2. Register handler in internal/app/app.go: a.provider.%sHandler().RegisterImplementation(a.server)\n", capitalize(serviceName))
	fmt.Printf("3. Implement your business logic\n")

	return nil
}

func getModulePath() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	return "", fmt.Errorf("module declaration not found in go.mod")
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func writeFile(path, content string) error {
	if err := ensureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
