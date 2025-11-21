package gosling

import (
    "fmt"
    "os"
    "path/filepath"

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
    }

    if !rmNoUsecase {
        if err := remover.RemoveUsecase(); err != nil {
            fmt.Printf("⚠ Warning: failed to remove usecase: %v\n", err)
        } else {
            fmt.Printf("✓ Usecase removed for service '%s'\n", serviceName)
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
        fmt.Printf("⚠ Manual cleanup required in app/provider.go\n")
    }

    fmt.Printf("\n🗑️  Service '%s' removed successfully!\n", serviceName)
    fmt.Println("\nManual cleanup required:")
    fmt.Printf("1. Remove %sUsecase from internal/usecase/usecase.go\n", capitalize(serviceName))
    if !rmNoRepository {
        fmt.Printf("2. Remove %sRepository from internal/repository/repository.go\n", capitalize(serviceName))
    }
    fmt.Printf("3. Remove handler registration from app/app.go\n")

    return nil
}

type Remover struct {
    serviceName string
    modulePath  string
}

func NewRemover(serviceName, modulePath string) *Remover {
    return &Remover{
        serviceName: serviceName,
        modulePath:  modulePath,
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
