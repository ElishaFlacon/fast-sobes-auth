package gosling

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Generator struct {
	ServiceName    string
	ServiceNameCap string
	ModulePath     string
	HandlerPath    string
	UsecasePath    string
	RepositoryPath string
	Methods        []string
}

func NewGenerator(serviceName, modulePath string, methods []string) *Generator {
	return &Generator{
		ServiceName:    serviceName,
		ServiceNameCap: capitalize(serviceName),
		ModulePath:     modulePath,
		HandlerPath:    filepath.Join("internal", "handlers", serviceName),
		UsecasePath:    filepath.Join("internal", "usecase", serviceName),
		RepositoryPath: filepath.Join("internal", "repository", serviceName),
		Methods:        methods,
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

func (g *Generator) UpdateUsecaseInterface() error {
	usecaseFilePath := filepath.Join("internal", "usecase", "usecase.go")

	// Read existing file
	data, err := os.ReadFile(usecaseFilePath)
	if err != nil {
		return err
	}

	content := string(data)

	// Check if interface already exists
	if strings.Contains(content, "type "+g.ServiceNameCap+"Usecase interface") {
		return fmt.Errorf("usecase interface '%s' already exists", g.ServiceNameCap+"Usecase")
	}

	// Generate interface
	var interfaceContent string
	if len(g.Methods) == 0 {
		interfaceContent = fmt.Sprintf(`
type %sUsecase interface {
	%sUseCase() error
}
`, g.ServiceNameCap, g.ServiceNameCap)
	} else {
		methods := ""
		for _, method := range g.Methods {
			methodCap := capitalize(method)
			methods += fmt.Sprintf("\t%s(ctx context.Context) error\n", methodCap)
		}
		interfaceContent = fmt.Sprintf(`
type %sUsecase interface {
%s}
`, g.ServiceNameCap, methods)
	}

	// Append to file
	content = strings.TrimRight(content, "\n") + "\n" + interfaceContent

	return os.WriteFile(usecaseFilePath, []byte(content), 0o644)
}

func (g *Generator) UpdateRepositoryInterface() error {
	repoFilePath := filepath.Join("internal", "repository", "repository.go")

	// Create file if not exists
	if _, err := os.Stat(repoFilePath); os.IsNotExist(err) {
		content := "package repository\n"
		if err := writeFile(repoFilePath, content); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(repoFilePath)
	if err != nil {
		return err
	}

	content := string(data)

	if strings.Contains(content, "type "+g.ServiceNameCap+"Repository interface") {
		return fmt.Errorf("repository interface '%s' already exists", g.ServiceNameCap+"Repository")
	}

	var interfaceContent string
	if len(g.Methods) == 0 {
		interfaceContent = fmt.Sprintf(`
type %sRepository interface {
	%sRepository() error
}
`, g.ServiceNameCap, g.ServiceNameCap)
	} else {
		methods := ""
		for _, method := range g.Methods {
			methodCap := capitalize(method)
			methods += fmt.Sprintf("\t%s(ctx context.Context) error\n", methodCap)
		}
		interfaceContent = fmt.Sprintf(`
type %sRepository interface {
%s}
`, g.ServiceNameCap, methods)
	}

	content = strings.TrimRight(content, "\n") + "\n" + interfaceContent

	return os.WriteFile(repoFilePath, []byte(content), 0o644)
}

func (g *Generator) UpdateProvider() error {
	providerPath := filepath.Join("internal", "app", "provider.go")

	data, err := os.ReadFile(providerPath)
	if err != nil {
		return fmt.Errorf("failed to read provider.go: %w", err)
	}

	content := string(data)

	if strings.Contains(content, g.ServiceName+"Usecase usecase."+g.ServiceNameCap+"Usecase") {
		return fmt.Errorf("service '%s' already exists in provider", g.ServiceName)
	}

	// Add imports
	importLines := []string{
		fmt.Sprintf(`	%sHandler "%s/internal/handlers/%s"`, g.ServiceName, g.ModulePath, g.ServiceName),
		fmt.Sprintf(`	%sUsecase "%s/internal/usecase/%s"`, g.ServiceName, g.ModulePath, g.ServiceName),
		fmt.Sprintf(`	%sRepository "%s/internal/repository/%s"`, g.ServiceName, g.ModulePath, g.ServiceName),
	}

	// Find last import
	lines := strings.Split(content, "\n")
	var newLines []string
	importInserted := false

	for i, line := range lines {
		newLines = append(newLines, line)

		if !importInserted && strings.Contains(line, `"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"`) {
			for _, imp := range importLines {
				newLines = append(newLines, imp)
			}
			importInserted = true
		}

		// Add to Provider struct
		if strings.Contains(line, "// Usecases") {
			newLines = append(newLines, fmt.Sprintf("\t%sUsecase usecase.%sUsecase", g.ServiceName, g.ServiceNameCap))
		}

		if strings.Contains(line, "// gRPC handlers") {
			newLines = append(newLines, fmt.Sprintf("\t%sHandler *%sHandler.Implementation", g.ServiceName, g.ServiceName))
		}

		if strings.Contains(line, "// Repositories") {
			newLines = append(newLines, fmt.Sprintf("\t%sRepository repository.%sRepository", g.ServiceName, g.ServiceNameCap))
		}

		// Add methods
		if strings.Contains(line, "// ----------------- gRPC HANDLER ------------------") && i+1 < len(lines) {
			handlerMethod := g.generateProviderHandlerMethod()
			newLines = append(newLines, "")
			newLines = append(newLines, handlerMethod)
		}

		if strings.Contains(line, "// -------------------- SERVICE --------------------") && i+1 < len(lines) {
			usecaseMethod := g.generateProviderUsecaseMethod()
			newLines = append(newLines, "")
			newLines = append(newLines, usecaseMethod)
		}

		if strings.Contains(line, "// ------------------ REPOSITORY -------------------") && i+1 < len(lines) {
			repoMethod := g.generateProviderRepositoryMethod()
			newLines = append(newLines, "")
			newLines = append(newLines, repoMethod)
		}
	}

	return os.WriteFile(providerPath, []byte(strings.Join(newLines, "\n")), 0o644)
}

func (g *Generator) generateHandlerContent() string {
	tmpl := `package {{.ServiceName}}

import (
	"context"

	"{{.ModulePath}}/internal/usecase"

	"google.golang.org/grpc"
)

type Implementation struct {
	usecase usecase.{{.ServiceNameCap}}Usecase
}

func NewImplementation(uc usecase.{{.ServiceNameCap}}Usecase) *Implementation {
	return &Implementation{
		usecase: uc,
	}
}

func (i *Implementation) RegisterImplementation(grpcServer *grpc.Server) {
	// Register your gRPC service here
}
{{if .Methods}}
{{range .Methods}}
func (i *Implementation) {{. | Capitalize}}(ctx context.Context) error {
	return i.usecase.{{. | Capitalize}}(ctx)
}
{{end}}
{{else}}
func (i *Implementation) {{.ServiceNameCap}}Handler(ctx context.Context) error {
	return i.usecase.{{.ServiceNameCap}}UseCase()
}
{{end}}
`
	funcMap := template.FuncMap{
		"Capitalize": capitalize,
	}

	t := template.Must(template.New("handler").Funcs(funcMap).Parse(tmpl))
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
	if len(g.Methods) == 0 {
		tmpl := `package {{.ServiceName}}

func (u *usecase) {{.ServiceNameCap}}UseCase() error {
	u.log.Infof("{{.ServiceNameCap}}UseCase called")
	return nil
}
`
		t := template.Must(template.New("usecaseImpl").Parse(tmpl))
		var buf strings.Builder
		t.Execute(&buf, g)
		return buf.String()
	}

	tmpl := `package {{.ServiceName}}

import (
	"context"
)
{{range .Methods}}
func (u *usecase) {{. | Capitalize}}(ctx context.Context) error {
	u.log.Infof("{{. | Capitalize}} called")
	return nil
}
{{end}}
`
	funcMap := template.FuncMap{
		"Capitalize": capitalize,
	}

	t := template.Must(template.New("usecaseImpl").Funcs(funcMap).Parse(tmpl))
	var buf strings.Builder
	t.Execute(&buf, g)
	return buf.String()
}

func (g *Generator) generateRepositoryContent() string {
	if len(g.Methods) == 0 {
		tmpl := `package {{.ServiceName}}

import (
	def "{{.ModulePath}}/internal/repository"
)

var _ def.{{.ServiceNameCap}}Repository = (*repository)(nil)

type repository struct {
}

func NewRepository() *repository {
	return &repository{}
}

func (r *repository) {{.ServiceNameCap}}Repository() error {
	return nil
}
`
		t := template.Must(template.New("repository").Parse(tmpl))
		var buf strings.Builder
		t.Execute(&buf, g)
		return buf.String()
	}

	tmpl := `package {{.ServiceName}}

import (
	"context"
	def "{{.ModulePath}}/internal/repository"
)

var _ def.{{.ServiceNameCap}}Repository = (*repository)(nil)

type repository struct {
}

func NewRepository() *repository {
	return &repository{}
}
{{range .Methods}}
func (r *repository) {{. | Capitalize}}(ctx context.Context) error {
	return nil
}
{{end}}
`
	funcMap := template.FuncMap{
		"Capitalize": capitalize,
	}

	t := template.Must(template.New("repository").Funcs(funcMap).Parse(tmpl))
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
}`, g.ServiceNameCap, g.ServiceName, g.ServiceName, g.ServiceName, g.ServiceName, g.ServiceNameCap, g.ServiceName)
}

func (g *Generator) generateProviderUsecaseMethod() string {
	return fmt.Sprintf(`func (p *Provider) %sUsecase() usecase.%sUsecase {
	if p.%sUsecase == nil {
		p.%sUsecase = %sUsecase.NewUsecase(p.log, p.%sRepository())
	}
	return p.%sUsecase
}`, g.ServiceNameCap, g.ServiceNameCap, g.ServiceName, g.ServiceName, g.ServiceName, g.ServiceNameCap, g.ServiceName)
}

func (g *Generator) generateProviderRepositoryMethod() string {
	return fmt.Sprintf(`func (p *Provider) %sRepository() repository.%sRepository {
	if p.%sRepository == nil {
		p.%sRepository = %sRepository.NewRepository()
	}
	return p.%sRepository
}`, g.ServiceNameCap, g.ServiceNameCap, g.ServiceName, g.ServiceName, g.ServiceName, g.ServiceName)
}
