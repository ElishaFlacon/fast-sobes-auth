package gosling

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rmNoUsecase    bool
	rmNoRepository bool
	rmNoHandler    bool
	rmNoApp        bool
)

var removeCmd = &cobra.Command{
	Use:   "remove service [name]",
	Short: "Remove an existing service",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runRemove,
}

func init() {
	removeCmd.Flags().BoolVar(&rmNoUsecase, "no-usecase", false, "Skip usecase removal")
	removeCmd.Flags().BoolVar(&rmNoRepository, "no-repository", false, "Skip repository removal")
	removeCmd.Flags().BoolVar(&rmNoHandler, "no-handler", false, "Skip handler removal")
	removeCmd.Flags().BoolVar(&rmNoApp, "no-app", false, "Skip app provider cleanup")
}

func runRemove(cmd *cobra.Command, args []string) error {
	serviceName := args[0]
	modulePath, err := getModulePath()
	if err != nil {
		return fmt.Errorf("failed to read module path: %w", err)
	}

	remover := NewRemover(serviceName, modulePath)

	if !rmNoRepository {
		if err := remover.RemoveRepository(); err != nil {
			fmt.Printf("⚠ Warning: failed to remove repository: %v\n", err)
		} else {
			fmt.Printf("✓ Repository removed for service '%s'\n", serviceName)
		}
		if err := remover.RemoveRepositoryInterface(); err != nil {
			fmt.Printf("⚠ Warning: failed to update repository interface: %v\n", err)
		}
	}

	if !rmNoUsecase {
		if err := remover.RemoveUsecase(); err != nil {
			fmt.Printf("⚠ Warning: failed to remove usecase: %v\n", err)
		} else {
			fmt.Printf("✓ Usecase removed for service '%s'\n", serviceName)
		}
		if err := remover.RemoveUsecaseInterface(); err != nil {
			fmt.Printf("⚠ Warning: failed to update usecase interface: %v\n", err)
		}
	}

	if !rmNoHandler {
		if err := remover.RemoveHandler(); err != nil {
			fmt.Printf("⚠ Warning: failed to remove handler: %v\n", err)
		} else {
			fmt.Printf("✓ Handler removed for service '%s'\n", serviceName)
		}
	}

	if !rmNoApp {
		if err := remover.RemoveFromProvider(); err != nil {
			fmt.Printf("⚠ Warning: failed to update provider: %v\n", err)
		} else {
			fmt.Printf("✓ Provider cleaned for service '%s'\n", serviceName)
		}
	}

	fmt.Printf("\n🗑️  Service '%s' removed successfully!\n", serviceName)

	return nil
}

type Remover struct {
	serviceName    string
	serviceNameCap string
	modulePath     string
}

func NewRemover(serviceName, modulePath string) *Remover {
	return &Remover{
		serviceName:    serviceName,
		serviceNameCap: capitalize(serviceName),
		modulePath:     modulePath,
	}
}

func (r *Remover) RemoveHandler() error {
	handlerPath := filepath.Join("internal", "handlers", r.serviceName)
	return os.RemoveAll(handlerPath)
}

func (r *Remover) RemoveUsecase() error {
	usecasePath := filepath.Join("internal", "usecase", r.serviceName)
	return os.RemoveAll(usecasePath)
}

func (r *Remover) RemoveRepository() error {
	repoPath := filepath.Join("internal", "repository", r.serviceName)
	return os.RemoveAll(repoPath)
}

func (r *Remover) RemoveUsecaseInterface() error {
	usecaseFilePath := filepath.Join("internal", "usecase", "usecase.go")
	data, err := os.ReadFile(usecaseFilePath)
	if err != nil {
		return err
	}

	content := string(data)
	interfaceName := r.serviceNameCap + "Usecase"

	// Remove interface declaration
	lines := strings.Split(content, "\n")
	var newLines []string
	skip := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Start of interface
		if strings.HasPrefix(trimmed, "type "+interfaceName+" interface") {
			skip = true
			continue
		}

		// End of interface
		if skip && trimmed == "}" {
			skip = false
			// Remove empty line after if exists
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "" {
				continue
			}
			continue
		}

		if !skip {
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(usecaseFilePath, []byte(strings.Join(newLines, "\n")), 0o644)
}

func (r *Remover) RemoveRepositoryInterface() error {
	repoFilePath := filepath.Join("internal", "repository", "repository.go")
	if _, err := os.Stat(repoFilePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(repoFilePath)
	if err != nil {
		return err
	}

	content := string(data)
	interfaceName := r.serviceNameCap + "Repository"

	lines := strings.Split(content, "\n")
	var newLines []string
	skip := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "type "+interfaceName+" interface") {
			skip = true
			continue
		}

		if skip && trimmed == "}" {
			skip = false
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "" {
				continue
			}
			continue
		}

		if !skip {
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(repoFilePath, []byte(strings.Join(newLines, "\n")), 0o644)
}

func (r *Remover) RemoveFromProvider() error {
	providerPath := filepath.Join("internal", "app", "provider.go")
	data, err := os.ReadFile(providerPath)
	if err != nil {
		return err
	}

	content := string(data)

	// Remove imports
	importToRemove1 := fmt.Sprintf(`%sHandler "%s/internal/handlers/%s"`, r.serviceName, r.modulePath, r.serviceName)
	importToRemove2 := fmt.Sprintf(`%sUsecase "%s/internal/usecase/%s"`, r.serviceName, r.modulePath, r.serviceName)
	importToRemove3 := fmt.Sprintf(`%sRepository "%s/internal/repository/%s"`, r.serviceName, r.modulePath, r.serviceName)

	content = removeLineContaining(content, importToRemove1)
	content = removeLineContaining(content, importToRemove2)
	content = removeLineContaining(content, importToRemove3)

	// Remove struct fields more precisely
	content = removeStructField(content, r.serviceName+"Usecase")
	content = removeStructField(content, r.serviceName+"Handler")
	content = removeStructField(content, r.serviceName+"Repository")

	// Remove methods
	content = removeMethod(content, fmt.Sprintf("func (p *Provider) %sHandler()", r.serviceNameCap))
	content = removeMethod(content, fmt.Sprintf("func (p *Provider) %sUsecase()", r.serviceNameCap))
	content = removeMethod(content, fmt.Sprintf("func (p *Provider) %sRepository()", r.serviceNameCap))

	return os.WriteFile(providerPath, []byte(content), 0o644)
}

func removeStructField(content, fieldName string) string {
	lines := strings.Split(content, "\n")
	var newLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Ищем строки вида: "fieldName type" или "fieldName *type" или "fieldName type.Type"
		if strings.HasPrefix(trimmed, fieldName+" ") || strings.HasPrefix(trimmed, fieldName+"\t") {
			// Пропускаем эту строку
			continue
		}
		newLines = append(newLines, line)
	}

	return strings.Join(newLines, "\n")
}

func removeLineContaining(content, substr string) string {
	lines := strings.Split(content, "\n")
	var newLines []string

	for _, line := range lines {
		if !strings.Contains(line, substr) {
			newLines = append(newLines, line)
		}
	}

	return strings.Join(newLines, "\n")
}

func removeMethod(content, methodSignature string) string {
	lines := strings.Split(content, "\n")
	var newLines []string
	skip := false
	braceCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, methodSignature) {
			skip = true
			braceCount = 0
		}

		if skip {
			braceCount += strings.Count(line, "{")
			braceCount -= strings.Count(line, "}")

			if braceCount == 0 && strings.Contains(line, "}") {
				skip = false
				continue
			}
			continue
		}

		newLines = append(newLines, line)
	}

	return strings.Join(newLines, "\n")
}
