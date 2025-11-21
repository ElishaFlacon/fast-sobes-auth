package gosling

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "text/template"
)

type Generator struct {
    ServiceName     string
    ServiceNameCap  string
    ModulePath      string
    HandlerPath     string
    UsecasePath     string
    RepositoryPath  string
}

func NewGenerator(serviceName, modulePath string) *Generator {
    return &Generator{
        ServiceName:    serviceName,
        ServiceNameCap: capitalize(serviceName),
        ModulePath:     modulePath,
        HandlerPath:    filepath.Join("internal", "handlers", serviceName),
        UsecasePath:    filepath.Join("internal", "usecase", serviceName),
        RepositoryPath: filepath.Join("internal", "repository", serviceName),
    }
}

func (g *Generator) GenerateHandler() error {
    if err := ensureDir(g.HandlerPath); err != nil {
        return err
    }

    handlerContent := g.generateHandlerContent()
    return writeFile(filepath.Join(g.HandlerPath, "handler.go"), handlerContent)
}

func (g *Generator) GenerateUsecase() error {
    if err := ensureDir(g.UsecasePath); err != nil {
        return err
    }

    usecaseContent := g.generateUsecaseContent()
    usecaseImplContent := g.generateUsecaseImplContent()

    if err := writeFile(filepath.Join(g.UsecasePath, "usecase.go"), usecaseContent); err != nil {
        return err
    }

    return writeFile(filepath.Join(g.UsecasePath, g.ServiceName+".go"), usecaseImplContent)
}

func (g *Generator) GenerateRepository() error {
    if err := ensureDir(g.RepositoryPath); err != nil {
        return err
    }

    repoContent := g.generateRepositoryContent()
    return writeFile(filepath.Join(g.RepositoryPath, "repository.go"), repoContent)
}

func (g *Generator) UpdateProvider() error {
    providerPath := filepath.Join("internal", "app", "provider.go")

    data, err := os.ReadFile(providerPath)
    if err != nil {
        return fmt.Errorf("failed to read provider.go: %w", err)
    }

    content := string(data)

    // Check if already exists
    if strings.Contains(content, g.ServiceName+"Usecase") {
        return fmt.Errorf("service '%s' already exists in provider", g.ServiceName)
    }

    // Add imports
    importBlock := fmt.Sprintf(`	%sHandler "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/%s"
	%sUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/%s"`,
        g.ServiceName, g.ServiceName, g.ServiceName, g.ServiceName)

    content = strings.Replace(content,
        `"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"`,
        fmt.Sprintf("%s\n\t\"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase\"", importBlock),
        1)

    // Add fields to Provider struct
    usecaseField := fmt.Sprintf("\n\t// %s\n\t%sUsecase usecase.%sUsecase",
        capitalize(g.ServiceName), g.ServiceName, capitalize(g.ServiceName))
    handlerField := fmt.Sprintf("\t%sHandler *%sHandler.Implementation",
        g.ServiceName, g.ServiceName)

    content = strings.Replace(content,
        "// Usecases",
        fmt.Sprintf("// Usecases%s\n", usecaseField),
        1)

    content = strings.Replace(content,
        "// gRPC handlers",
        fmt.Sprintf("// gRPC handlers\n%s\n", handlerField),
        1)

    // Add methods
    handlerMethod := g.generateProviderHandlerMethod()
    usecaseMethod := g.generateProviderUsecaseMethod()

    content = strings.Replace(content,
        "// ----------------- gRPC HANDLER ------------------",
        fmt.Sprintf("// ----------------- gRPC HANDLER ------------------\n\n%s", handlerMethod),
        1)

    content = strings.Replace(content,
        "// -------------------- SERVICE --------------------",
        fmt.Sprintf("// -------------------- SERVICE --------------------\n\n%s", usecaseMethod),
        1)

    return os.WriteFile(providerPath, []byte(content), 0644)
}

func (g *Generator) generateHandlerContent() string {
    tmpl := `package {{.ServiceName}}

import (
	"context"

	"{{.ModulePath}}/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Implementation struct {
	// TODO: Add UnimplementedServer from your protobuf package
	// Example: pb.Unimplemented{{.ServiceNameCap}}Server
	usecase usecase.{{.ServiceNameCap}}Usecase
}

func NewImplementation(uc usecase.{{.ServiceNameCap}}Usecase) *Implementation {
	return &Implementation{
		usecase: uc,
	}
}

func (i *Implementation) RegisterImplementation(grpcServer *grpc.Server) {
	// TODO: Register your service with gRPC server
	// Example: pb.Register{{.ServiceNameCap}}Server(grpcServer, i)
}

// TODO: Implement your gRPC methods here
// Example:
// func (i *Implementation) SomeMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
//     result := i.usecase.SomeMethod(ctx, req)
//     return &pb.Response{Data: result}, nil
// }
`
    t := template.Must(template.New("handler").Parse(tmpl))
    var buf strings.Builder
    t.Execute(&buf, g)
    return buf.String()
}

func (g *Generator) generateUsecaseContent() string {
    tmpl := `package {{.ServiceName}}

import (
	"{{.ModulePath}}/internal/domain"
	def "{{.ModulePath}}/internal/usecase"
	"{{.ModulePath}}/internal/repository"
)

var _ def.{{.ServiceNameCap}}Usecase = (*usecase)(nil)

type usecase struct {
	log        domain.Logger
	repository repository.{{.ServiceNameCap}}Repository
}

func NewUsecase(
	log domain.Logger,
	repository repository.{{.ServiceNameCap}}Repository,
) *usecase {
	return &usecase{
		log:        log,
		repository: repository,
	}
}
`
    t := template.Must(template.New("usecase").Parse(tmpl))
    var buf strings.Builder
    t.Execute(&buf, g)
    return buf.String()
}

func (g *Generator) generateUsecaseImplContent() string {
    tmpl := `package {{.ServiceName}}

import (
	"context"
)

// TODO: Implement your business logic methods here
// Example:
// func (u *usecase) GetData(ctx context.Context, id string) (string, error) {
//     u.log.Infof("Getting data for id: %s", id)
//     return u.repository.FindByID(ctx, id)
// }
`
    t := template.Must(template.New("usecaseImpl").Parse(tmpl))
    var buf strings.Builder
    t.Execute(&buf, g)
    return buf.String()
}

func (g *Generator) generateRepositoryContent() string {
    tmpl := `package {{.ServiceName}}

import (
	"context"
	def "{{.ModulePath}}/internal/repository"
)

var _ def.{{.ServiceNameCap}}Repository = (*repository)(nil)

type repository struct {
	// TODO: Add database connection
	// db *sql.DB or gorm.DB
}

func NewRepository( /* TODO: add dependencies */) *repository {
	return &repository{
		// Initialize fields
	}
}

// TODO: Implement repository methods here
// Example:
// func (r *repository) FindByID(ctx context.Context, id string) (string, error) {
//     // Database query logic
//     return "", nil
// }
`
    t := template.Must(template.New("repository").Parse(tmpl))
    var buf strings.Builder
    t.Execute(&buf, g)
    return buf.String()
}

func (g *Generator) generateProviderHandlerMethod() string {
    return fmt.Sprintf(`func (p *Provider) %sHandler() *%sHandler.Implementation {
	if p.%sHandler == nil {
		p.%sHandler = %sHandler.NewImplementation(p.%sUsecase())
	}
	return p.%sHandler
}
`, g.ServiceNameCap, g.ServiceName, g.ServiceName, g.ServiceName, g.ServiceName, g.ServiceNameCap, g.ServiceName)
}

func (g *Generator) generateProviderUsecaseMethod() string {
    return fmt.Sprintf(`func (p *Provider) %sUsecase() usecase.%sUsecase {
	if p.%sUsecase == nil {
		p.%sUsecase = %sUsecase.NewUsecase(p.log, p.%sRepository())
	}
	return p.%sUsecase
}
`, g.ServiceNameCap, g.ServiceNameCap, g.ServiceName, g.ServiceName, g.ServiceName, g.ServiceNameCap, g.ServiceName)
}
