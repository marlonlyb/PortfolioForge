package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
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
	"github.com/marlonlyb/portfolioforge/infrastructure/casestudy"
)

type canonicalPublishTarget struct {
	Slug        string
	LocalDir    string
	LocalFile   string
	RemoteDir   string
	RemoteFile  string
	PublicURL   string
	RemoteBase  string
	DisplayRoot string
}

func runCanonicalPublish(args []string) error {
	fs := flag.NewFlagSet("canonical-publish", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	caseDir := fs.String("case-dir", "", "Ruta del caso de estudio o del directorio 90. dev_portfolioforge")
	slug := fs.String("slug", "", "Slug explícito del proyecto si el directorio contiene más de un canonical")
	dryRun := fs.Bool("dry-run", false, "Solo imprime qué se publicaría sin subir archivos")

	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: go run ./cmd canonical-publish --case-dir <ruta> [--slug <slug>] [--dry-run]")
		fmt.Fprintln(fs.Output(), "")
		fmt.Fprintln(fs.Output(), "Resuelve el canonical local dentro de 90. dev_portfolioforge/<slug>/ y publica toda la carpeta por FTPS.")
		fmt.Fprintln(fs.Output(), "La URL pública esperada es https://mlbautomation.com/dev/portfolioforge/<slug>/<slug>.md")
		fmt.Fprintln(fs.Output(), "")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*caseDir) == "" {
		fs.Usage()
		return errors.New("case-dir es obligatorio")
	}

	config, err := casestudy.LoadPublishConfigFromEnv(os.Getenv)
	if err != nil {
		return err
	}

	target, err := casestudy.ResolvePublishTarget(*caseDir, *slug, casestudy.PublishConfig(config))
	if err != nil {
		return err
	}

	files, err := casestudy.CollectPublishFiles(target.LocalDir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no se encontraron archivos para publicar en %s", target.LocalDir)
	}

	fmt.Printf("canonical-publish: slug=%s\n", target.Slug)
	fmt.Printf("canonical-publish: local=%s\n", target.LocalDir)
	fmt.Printf("canonical-publish: remote=%s\n", target.RemoteDir)
	fmt.Printf("canonical-publish: url=%s\n", target.PublicURL)

	if *dryRun {
		for _, file := range files {
			rel, _ := filepath.Rel(target.LocalDir, file)
			fmt.Printf("dry-run: %s -> %s\n", file, path.Join(target.RemoteDir, filepath.ToSlash(rel)))
		}
		return nil
	}

	if err := casestudy.PublishFiles(target, files, casestudy.PublishConfig(config)); err != nil {
		return err
	}

	if err := casestudy.VerifyPublishedCanonical(context.Background(), target.PublicURL); err != nil {
		return err
	}

	fmt.Printf("canonical-publish: publicado OK %s\n", target.PublicURL)
	return nil
}

type canonicalPublishConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	PublicBase string
	RemoteBase string
	Timeout    time.Duration
}

func loadCanonicalPublishConfigFromEnv() (canonicalPublishConfig, error) {
	config := canonicalPublishConfig{
		Host:       strings.TrimSpace(os.Getenv("PF_FTP_HOST")),
		Port:       strings.TrimSpace(os.Getenv("PF_FTP_PORT")),
		User:       strings.TrimSpace(os.Getenv("PF_FTP_USER")),
		Password:   os.Getenv("PF_FTP_PASSWORD"),
		PublicBase: strings.TrimSpace(os.Getenv("PF_PUBLIC_BASE")),
		RemoteBase: strings.TrimSpace(os.Getenv("PF_FTP_REMOTE_BASE")),
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
		return canonicalPublishConfig{}, fmt.Errorf("faltan variables de entorno para canonical-publish: %s", strings.Join(missing, ", "))
	}

	return config, nil
}

func resolveCanonicalPublishTarget(inputPath, explicitSlug string, config canonicalPublishConfig) (canonicalPublishTarget, error) {
	abs, err := filepath.Abs(filepath.Clean(inputPath))
	if err != nil {
		return canonicalPublishTarget{}, err
	}

	slugDir, slug, err := resolveCanonicalSlugDir(abs, explicitSlug)
	if err != nil {
		return canonicalPublishTarget{}, err
	}

	localFile := filepath.Join(slugDir, slug+".md")
	info, err := os.Stat(localFile)
	if err != nil {
		return canonicalPublishTarget{}, fmt.Errorf("canonical markdown no encontrado: %w", err)
	}
	if info.IsDir() {
		return canonicalPublishTarget{}, fmt.Errorf("se esperaba un archivo markdown y se encontró un directorio: %s", localFile)
	}

	remoteDir := cleanRemoteJoin(config.RemoteBase, slug)
	remoteFile := cleanRemoteJoin(remoteDir, slug+".md")
	publicBase := strings.TrimRight(config.PublicBase, "/")

	return canonicalPublishTarget{
		Slug:        slug,
		LocalDir:    slugDir,
		LocalFile:   localFile,
		RemoteDir:   remoteDir,
		RemoteFile:  remoteFile,
		PublicURL:   publicBase + "/" + slug + "/" + slug + ".md",
		RemoteBase:  config.RemoteBase,
		DisplayRoot: abs,
	}, nil
}

func resolveCanonicalSlugDir(absPath, explicitSlug string) (string, string, error) {
	trimmedSlug := strings.TrimSpace(explicitSlug)

	if trimmedSlug != "" {
		if isSlugDir(absPath, trimmedSlug) {
			return absPath, trimmedSlug, nil
		}
		if filepath.Base(absPath) == "90. dev_portfolioforge" {
			candidate := filepath.Join(absPath, trimmedSlug)
			if isSlugDir(candidate, trimmedSlug) {
				return candidate, trimmedSlug, nil
			}
		}
		candidate := filepath.Join(absPath, "90. dev_portfolioforge", trimmedSlug)
		if isSlugDir(candidate, trimmedSlug) {
			return candidate, trimmedSlug, nil
		}
		return "", "", fmt.Errorf("no se encontró el directorio del slug %q dentro de %s", trimmedSlug, absPath)
	}

	if slug := filepath.Base(absPath); isSlugDir(absPath, slug) {
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
		if isSlugDir(filepath.Join(canonicalRoot, slug), slug) {
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

func isSlugDir(dirPath, slug string) bool {
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

func cleanRemoteJoin(parts ...string) string {
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

func collectCanonicalPublishFiles(localDir string) ([]string, error) {
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

func publishCanonicalFiles(target canonicalPublishTarget, files []string, config canonicalPublishConfig) error {
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
		remotePath := cleanRemoteJoin(target.RemoteDir, filepath.ToSlash(rel))
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
	cleaned := cleanRemoteJoin(remoteDir)
	parts := strings.Split(strings.Trim(cleaned, "/"), "/")
	current := "/"
	for _, part := range parts {
		if part == "" {
			continue
		}
		current = cleanRemoteJoin(current, part)
		if err := conn.MakeDir(current); err != nil && !isFTPIgnorableDirError(err) {
			return fmt.Errorf("crear directorio remoto %s: %w", current, err)
		}
	}
	return nil
}

func isFTPIgnorableDirError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "file exists") || strings.Contains(message, "directory already exists") || strings.Contains(message, "550")
}

func verifyPublishedCanonical(publicURL string) error {
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, publicURL, nil)
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
