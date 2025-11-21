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
	Methods        []Method
}

type Method struct {
	Name    string // lowercase для файлов (auth)
	NameCap string // capitalized для методов (Auth)
}

func NewGenerator(serviceName, modulePath string, methods []string) *Generator {
	normalizedMethods := make([]Method, len(methods))
	for i, m := range methods {
		normalizedMethods[i] = Method{
			Name:    strings.ToLower(m),
			NameCap: capitalize(strings.ToLower(m)),
		}
	}

	return &Generator{
		ServiceName:    strings.ToLower(serviceName),
		ServiceNameCap: capitalize(strings.ToLower(serviceName)),
		ModulePath:     modulePath,
		HandlerPath:    filepath.Join("internal", "handlers", strings.ToLower(serviceName)),
		UsecasePath:    filepath.Join("internal", "usecase", strings.ToLower(serviceName)),
		RepositoryPath: filepath.Join("internal", "repository", strings.ToLower(serviceName)),
		Methods:        normalizedMethods,
	}
}

func (g *Generator) GenerateHandler() error {
	if err := ensureDir(g.HandlerPath); err != nil {
		return err
	}

	// Генерация базового handler.go
	handlerContent := g.generateHandlerBaseContent()
	if err := writeFile(filepath.Join(g.HandlerPath, "handler.go"), handlerContent); err != nil {
		return err
	}

	// Генерация файла для каждого метода или дефолтного файла
	if len(g.Methods) == 0 {
		// Дефолтный метод в отдельном файле
		defaultContent := g.generateHandlerDefaultMethod()
		return writeFile(filepath.Join(g.HandlerPath, g.ServiceName+".go"), defaultContent)
	}

	// Создаем отдельный файл для каждого метода
	for _, method := range g.Methods {
		methodContent := g.generateHandlerMethodContent(method)
		filename := filepath.Join(g.HandlerPath, method.Name+".go")
		if err := writeFile(filename, methodContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) GenerateUsecase() error {
	if err := ensureDir(g.UsecasePath); err != nil {
		return err
	}

	// Базовый usecase.go
	usecaseContent := g.generateUsecaseContent()
	if err := writeFile(filepath.Join(g.UsecasePath, "usecase.go"), usecaseContent); err != nil {
		return err
	}

	// Генерация методов
	if len(g.Methods) == 0 {
		// Дефолтный метод
		defaultContent := g.generateUsecaseDefaultMethod()
		return writeFile(filepath.Join(g.UsecasePath, g.ServiceName+".go"), defaultContent)
	}

	// Отдельный файл для каждого метода
	for _, method := range g.Methods {
		methodContent := g.generateUsecaseMethodContent(method)
		filename := filepath.Join(g.UsecasePath, method.Name+".go")
		if err := writeFile(filename, methodContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) GenerateRepository() error {
	if err := ensureDir(g.RepositoryPath); err != nil {
		return err
	}

	// Базовый repository.go
	repoContent := g.generateRepositoryBaseContent()
	if err := writeFile(filepath.Join(g.RepositoryPath, "repository.go"), repoContent); err != nil {
		return err
	}

	// Генерация методов
	if len(g.Methods) == 0 {
		// Дефолтный метод
		defaultContent := g.generateRepositoryDefaultMethod()
		return writeFile(filepath.Join(g.RepositoryPath, g.ServiceName+".go"), defaultContent)
	}

	// Отдельный файл для каждого метода
	for _, method := range g.Methods {
		methodContent := g.generateRepositoryMethodContent(method)
		filename := filepath.Join(g.RepositoryPath, method.Name+".go")
		if err := writeFile(filename, methodContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) UpdateUsecaseInterface() error {
	usecaseFilePath := filepath.Join("internal", "usecase", "usecase.go")

	data, err := os.ReadFile(usecaseFilePath)
	if err != nil {
		return err
	}

	content := string(data)

	if strings.Contains(content, "type "+g.ServiceNameCap+"Usecase interface") {
		return fmt.Errorf("usecase interface '%s' already exists", g.ServiceNameCap+"Usecase")
	}

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
			methods += fmt.Sprintf("\t%s() error\n", method.NameCap)
		}
		interfaceContent = fmt.Sprintf(`
type %sUsecase interface {
%s}
`, g.ServiceNameCap, methods)
	}

	content = strings.TrimRight(content, "\n") + "\n" + interfaceContent

	return os.WriteFile(usecaseFilePath, []byte(content), 0o644)
}

func (g *Generator) UpdateRepositoryInterface() error {
	repoFilePath := filepath.Join("internal", "repository", "repository.go")

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
			methods += fmt.Sprintf("\t%s() error\n", method.NameCap)
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

	importLines := []string{
		fmt.Sprintf(`   %sHandler "%s/internal/handlers/%s"`, g.ServiceName, g.ModulePath, g.ServiceName),
		fmt.Sprintf(`   %sUsecase "%s/internal/usecase/%s"`, g.ServiceName, g.ModulePath, g.ServiceName),
		fmt.Sprintf(`   %sRepository "%s/internal/repository/%s"`, g.ServiceName, g.ModulePath, g.ServiceName),
	}

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

		if strings.Contains(line, "// Usecases") {
			newLines = append(newLines, fmt.Sprintf("\t%sUsecase usecase.%sUsecase", g.ServiceName, g.ServiceNameCap))
		}

		if strings.Contains(line, "// gRPC handlers") {
			newLines = append(newLines, fmt.Sprintf("\t%sHandler *%sHandler.Implementation", g.ServiceName, g.ServiceName))
		}

		if strings.Contains(line, "// Repositories") {
			newLines = append(newLines, fmt.Sprintf("\t%sRepository repository.%sRepository", g.ServiceName, g.ServiceNameCap))
		}

		if strings.Contains(line, "// ----------------- gRPC HANDLER ------------------") && i+1 < len(lines) {
			handlerMethod := g.generateProviderHandlerMethod()
			newLines = append(newLines, "")
			newLines = append(newLines, handlerMethod)
		}

		if strings.Contains(line, "// -------------------- USECASE --------------------") && i+1 < len(lines) {
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

// Шаблоны для Handler
func (g *Generator) generateHandlerBaseContent() string {
	tmpl := `package {{.ServiceName}}

import (
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
`
	t := template.Must(template.New("handler").Parse(tmpl))
	var buf strings.Builder
	t.Execute(&buf, g)
	return buf.String()
}

func (g *Generator) generateHandlerDefaultMethod() string {
	tmpl := `package {{.ServiceName}}

import "context"

func (i *Implementation) {{.ServiceNameCap}}Handler(ctx context.Context) error {
    return i.usecase.{{.ServiceNameCap}}UseCase()
}
`
	t := template.Must(template.New("handlerDefault").Parse(tmpl))
	var buf strings.Builder
	t.Execute(&buf, g)
	return buf.String()
}

func (g *Generator) generateHandlerMethodContent(method Method) string {
	tmpl := `package {{.ServiceName}}

import "context"

func (i *Implementation) {{.MethodNameCap}}(ctx context.Context) error {
    return i.usecase.{{.MethodNameCap}}()
}
`
	data := struct {
		ServiceName   string
		MethodNameCap string
	}{
		ServiceName:   g.ServiceName,
		MethodNameCap: method.NameCap,
	}

	t := template.Must(template.New("handlerMethod").Parse(tmpl))
	var buf strings.Builder
	t.Execute(&buf, data)
	return buf.String()
}

// Шаблоны для Usecase
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

func (g *Generator) generateUsecaseDefaultMethod() string {
	tmpl := `package {{.ServiceName}}

func (u *usecase) {{.ServiceNameCap}}UseCase() error {
    u.log.Infof("{{.ServiceNameCap}}UseCase called")
    return nil
}
`
	t := template.Must(template.New("usecaseDefault").Parse(tmpl))
	var buf strings.Builder
	t.Execute(&buf, g)
	return buf.String()
}

func (g *Generator) generateUsecaseMethodContent(method Method) string {
	tmpl := `package {{.ServiceName}}

func (u *usecase) {{.MethodNameCap}}() error {
    u.log.Infof("{{.MethodNameCap}} called")
    return nil
}
`
	data := struct {
		ServiceName   string
		MethodNameCap string
	}{
		ServiceName:   g.ServiceName,
		MethodNameCap: method.NameCap,
	}

	t := template.Must(template.New("usecaseMethod").Parse(tmpl))
	var buf strings.Builder
	t.Execute(&buf, data)
	return buf.String()
}

// Шаблоны для Repository
func (g *Generator) generateRepositoryBaseContent() string {
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
`
	t := template.Must(template.New("repository").Parse(tmpl))
	var buf strings.Builder
	t.Execute(&buf, g)
	return buf.String()
}

func (g *Generator) generateRepositoryDefaultMethod() string {
	tmpl := `package {{.ServiceName}}

func (r *repository) {{.ServiceNameCap}}Repository() error {
    return nil
}
`
	t := template.Must(template.New("repositoryDefault").Parse(tmpl))
	var buf strings.Builder
	t.Execute(&buf, g)
	return buf.String()
}

func (g *Generator) generateRepositoryMethodContent(method Method) string {
	tmpl := `package {{.ServiceName}}

func (r *repository) {{.MethodNameCap}}() error {
    return nil
}
`
	data := struct {
		ServiceName   string
		MethodNameCap string
	}{
		ServiceName:   g.ServiceName,
		MethodNameCap: method.NameCap,
	}

	t := template.Must(template.New("repositoryMethod").Parse(tmpl))
	var buf strings.Builder
	t.Execute(&buf, data)
	return buf.String()
}

// Provider методы
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
