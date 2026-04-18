package casestudy

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type PublishConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	PublicBase string
	RemoteBase string
	Timeout    time.Duration
}

type PublishTarget struct {
	Slug               string
	LocalDir           string
	LocalFile          string
	RemoteDir          string
	RemoteFile         string
	PublicURL          string
	RemoteBase         string
	DisplayRoot        string
	RequestedInputPath string
}

func LoadPublishConfigFromEnv(getenv func(string) string) (PublishConfig, error) {
	config := PublishConfig{
		Host:       strings.TrimSpace(getenv("PF_FTP_HOST")),
		Port:       strings.TrimSpace(getenv("PF_FTP_PORT")),
		User:       strings.TrimSpace(getenv("PF_FTP_USER")),
		Password:   getenv("PF_FTP_PASSWORD"),
		PublicBase: strings.TrimSpace(getenv("PF_PUBLIC_BASE")),
		RemoteBase: strings.TrimSpace(getenv("PF_FTP_REMOTE_BASE")),
		Timeout:    30 * time.Second,
	}

	if config.Port == "" {
		config.Port = "21"
	}
	if config.RemoteBase == "" {
		config.RemoteBase = "/"
	}

	missing := make([]string, 0)
	if config.Host == "" {
		missing = append(missing, "PF_FTP_HOST")
	}
	if config.User == "" {
		missing = append(missing, "PF_FTP_USER")
	}
	if config.Password == "" {
		missing = append(missing, "PF_FTP_PASSWORD")
	}
	if config.PublicBase == "" {
		missing = append(missing, "PF_PUBLIC_BASE")
	}
	if len(missing) > 0 {
		return PublishConfig{}, fmt.Errorf("faltan variables de entorno para canonical-publish: %s", strings.Join(missing, ", "))
	}

	return config, nil
}

func ResolvePublishTarget(inputPath, explicitSlug string, config PublishConfig) (PublishTarget, error) {
	abs, err := filepath.Abs(filepath.Clean(inputPath))
	if err != nil {
		return PublishTarget{}, err
	}

	slugDir, slug, err := ResolveCanonicalSlugDir(abs, explicitSlug)
	if err != nil {
		return PublishTarget{}, err
	}

	localFile := filepath.Join(slugDir, slug+".md")
	info, err := os.Stat(localFile)
	if err != nil {
		return PublishTarget{}, fmt.Errorf("canonical markdown no encontrado: %w", err)
	}
	if info.IsDir() {
		return PublishTarget{}, fmt.Errorf("se esperaba un archivo markdown y se encontró un directorio: %s", localFile)
	}

	remoteDir := CleanRemoteJoin(config.RemoteBase, slug)
	remoteFile := CleanRemoteJoin(remoteDir, slug+".md")
	publicBase := strings.TrimRight(config.PublicBase, "/")

	return PublishTarget{
		Slug:               slug,
		LocalDir:           slugDir,
		LocalFile:          localFile,
		RemoteDir:          remoteDir,
		RemoteFile:         remoteFile,
		PublicURL:          publicBase + "/" + slug + "/" + slug + ".md",
		RemoteBase:         config.RemoteBase,
		DisplayRoot:        abs,
		RequestedInputPath: inputPath,
	}, nil
}

func ResolveCanonicalSlugDir(absPath, explicitSlug string) (string, string, error) {
	trimmedSlug := strings.TrimSpace(explicitSlug)

	if trimmedSlug != "" {
		if IsSlugDir(absPath, trimmedSlug) {
			return absPath, trimmedSlug, nil
		}
		if filepath.Base(absPath) == "90. dev_portfolioforge" {
			candidate := filepath.Join(absPath, trimmedSlug)
			if IsSlugDir(candidate, trimmedSlug) {
				return candidate, trimmedSlug, nil
			}
		}
		candidate := filepath.Join(absPath, "90. dev_portfolioforge", trimmedSlug)
		if IsSlugDir(candidate, trimmedSlug) {
			return candidate, trimmedSlug, nil
		}
		return "", "", fmt.Errorf("no se encontró el directorio del slug %q dentro de %s", trimmedSlug, absPath)
	}

	if slug := filepath.Base(absPath); IsSlugDir(absPath, slug) {
		return absPath, slug, nil
	}

	canonicalRoot := absPath
	if filepath.Base(absPath) != "90. dev_portfolioforge" {
		canonicalRoot = filepath.Join(absPath, "90. dev_portfolioforge")
	}

	entries, err := os.ReadDir(canonicalRoot)
	if err != nil {
		return "", "", fmt.Errorf("no se pudo leer 90. dev_portfolioforge: %w", err)
	}

	slugDirs := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		slug := entry.Name()
		if IsSlugDir(filepath.Join(canonicalRoot, slug), slug) {
			slugDirs = append(slugDirs, slug)
		}
	}

	if len(slugDirs) == 0 {
		return "", "", fmt.Errorf("no se encontró ningún directorio canonical válido dentro de %s", canonicalRoot)
	}
	if len(slugDirs) > 1 {
		sort.Strings(slugDirs)
		return "", "", fmt.Errorf("hay múltiples directorios canonical en %s; usa --slug (%s)", canonicalRoot, strings.Join(slugDirs, ", "))
	}

	return filepath.Join(canonicalRoot, slugDirs[0]), slugDirs[0], nil
}

func IsSlugDir(dirPath, slug string) bool {
	info, err := os.Stat(dirPath)
	if err != nil || !info.IsDir() {
		return false
	}
	fileInfo, err := os.Stat(filepath.Join(dirPath, slug+".md"))
	if err != nil || fileInfo.IsDir() {
		return false
	}
	return true
}

func CleanRemoteJoin(parts ...string) string {
	built := "/"
	for _, part := range parts {
		trimmed := strings.Trim(part, "/")
		if trimmed == "" {
			continue
		}
		built = path.Join(built, trimmed)
	}
	return built
}

func CollectPublishFiles(localDir string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(localDir, func(current string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, current)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func PublishFiles(target PublishTarget, files []string, config PublishConfig) error {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	tlsConfig := &tls.Config{ServerName: config.Host, MinVersion: tls.VersionTLS12}

	conn, err := ftp.Dial(addr,
		ftp.DialWithTimeout(config.Timeout),
		ftp.DialWithExplicitTLS(tlsConfig),
	)
	if err != nil {
		return fmt.Errorf("ftps dial: %w", err)
	}
	defer conn.Quit()

	if err := conn.Login(config.User, config.Password); err != nil {
		return fmt.Errorf("ftps login: %w", err)
	}

	if err := ensureRemoteDir(conn, target.RemoteDir); err != nil {
		return err
	}

	for _, file := range files {
		rel, err := filepath.Rel(target.LocalDir, file)
		if err != nil {
			return err
		}
		remotePath := CleanRemoteJoin(target.RemoteDir, filepath.ToSlash(rel))
		if err := ensureRemoteDir(conn, path.Dir(remotePath)); err != nil {
			return err
		}

		handle, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("abrir archivo local %s: %w", file, err)
		}

		if err := conn.Stor(remotePath, handle); err != nil {
			handle.Close()
			return fmt.Errorf("subir archivo %s -> %s: %w", file, remotePath, err)
		}
		if err := handle.Close(); err != nil {
			return fmt.Errorf("cerrar archivo local %s: %w", file, err)
		}
	}

	return nil
}

func ensureRemoteDir(conn *ftp.ServerConn, remoteDir string) error {
	cleaned := CleanRemoteJoin(remoteDir)
	parts := strings.Split(strings.Trim(cleaned, "/"), "/")
	current := "/"
	for _, part := range parts {
		if part == "" {
			continue
		}
		current = CleanRemoteJoin(current, part)
		if err := conn.MakeDir(current); err != nil && !isIgnorableDirError(err) {
			return fmt.Errorf("crear directorio remoto %s: %w", current, err)
		}
	}
	return nil
}

func isIgnorableDirError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "file exists") || strings.Contains(message, "directory already exists") || strings.Contains(message, "550")
}

func VerifyPublishedCanonical(ctx context.Context, publicURL string) error {
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, publicURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("verificar URL publicada %s: %w", publicURL, err)
	}
	defer resp.Body.Close()
	_, _ = io.CopyN(io.Discard, resp.Body, 1024)

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("la URL publicada respondió %d para %s", resp.StatusCode, publicURL)
	}
	return nil
}
