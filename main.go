package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	useflagly "github.com/useflagly/sdk-go"
	"github.com/useflagly/sdk-go/models"
)

func loadDotEnv(paths ...string) {
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				if key != "" && os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
		return
	}
}

func prettyPrint(label string, v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s: %s\n\n", label, string(b))
}

func ptr[T any](v T) *T { return &v }

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	loadDotEnv(".env", "../.env")

	apiKey := os.Getenv("FLAGLY_API_KEY")
	if apiKey == "" {
		log.Fatal("FLAGLY_API_KEY não definida")
	}
	identifier := os.Getenv("FLAGLY_IDENTIFIER")
	if identifier == "" {
		log.Fatal("FLAGLY_IDENTIFIER não definida")
	}
	slug := getEnv("FLAGLY_SLUG", "teste-1")
	environment := getEnv("FLAGLY_ENVIRONMENT", "HML")

	client := useflagly.New(useflagly.Options{Token: apiKey})
	ctx := context.Background()

	// --- Health Check ---
	health, err := client.HealthCheck(ctx)
	if err != nil {
		log.Printf("Health erro: %v", err)
	} else {
		prettyPrint("Health", health)
	}

	// --- Initialize ---
	_, err = client.Initialize(ctx, models.ReceiveMessage{
		Identifier: identifier,
		Slug:       slug,
	}, environment)
	// Initialize retorna um número (session id), não um JSON object — erro de unmarshal é esperado
	if err != nil && !strings.Contains(err.Error(), "cannot unmarshal") {
		log.Printf("Initialize erro: %v", err)
	} else {
		fmt.Printf("Initialize: ok\n\n")
	}

	// --- Result ---
	result, err := client.GetResult(ctx, identifier)
	if err != nil {
		log.Fatalf("GetResult erro: %v", err)
	}
	prettyPrint("Result", result)

	// --- Validar usando os slugs do resultado ---
	data, ok := result["data"].(map[string]any)
	if !ok {
		return
	}

	for flowSlug, flowVal := range data {
		flowParts, ok := flowVal.(map[string]any)
		if !ok {
			continue // ignora entradas não-objeto (ex: flags diretos no data)
		}

		flowResult, err := client.ValidateFlow(ctx, flowSlug, models.ValidateBody{Identifier: ptr(identifier)}, environment)
		if err != nil {
			log.Printf("ValidateFlow (%s) erro: %v", flowSlug, err)
		} else {
			prettyPrint("ValidateFlow ("+flowSlug+")", flowResult)
		}

		for fpSlug, fpVal := range flowParts {
			flags, ok := fpVal.(map[string]any)
			if !ok {
				continue
			}

			fpResult, err := client.ValidateFlowPart(ctx, fpSlug, models.ValidateBody{Identifier: ptr(identifier)}, environment)
			if err != nil {
				log.Printf("ValidateFlowPart (%s) erro: %v", fpSlug, err)
			} else {
				prettyPrint("ValidateFlowPart ("+fpSlug+")", fpResult)
			}

			for flagSlug := range flags {
				flagResult, err := client.ValidateFlag(ctx, flagSlug, models.ValidateBody{Identifier: ptr(identifier)}, environment)
				if err != nil {
					log.Printf("ValidateFlag (%s) erro: %v", flagSlug, err)
				} else {
					prettyPrint("ValidateFlag ("+flagSlug+")", flagResult)
				}
			}
		}
	}
}

