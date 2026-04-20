package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	useflagly "github.com/useflagly/sdk-go"
	"github.com/useflagly/sdk-go/models"
)

func prettyPrint(label string, v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s: %s\n\n", label, string(b))
}

func main() {
	apiKey := os.Getenv("FLAGLY_API_KEY")
	if apiKey == "" {
		log.Fatal("FLAGLY_API_KEY não definida")
	}

	client := useflagly.New(useflagly.Options{
		Token: apiKey,
		// BaseURL: "https://api.useflagly.com.br", // opcional
	})

	ctx := context.Background()

	// --- Health Check ---
	health, err := client.HealthCheck(ctx)
	if err != nil {
		log.Fatalf("HealthCheck: %v", err)
	}
	prettyPrint("Health", health)

	// --- Validar Feature Flag ---
	flagResult, err := client.ValidateFlag(ctx, "minha-feature", models.ValidateBody{
		Identifier: ptr("user-123"),
		Context:    map[string]any{"plano": "premium", "pais": "BR"},
	}, "production")
	if err != nil {
		log.Printf("ValidateFlag erro: %v", err)
	} else {
		prettyPrint("Flag resultado", flagResult)
	}

	// --- Validar Flow ---
	flowResult, err := client.ValidateFlow(ctx, "meu-fluxo", models.ValidateBody{
		Identifier: ptr("user-123"),
	}, "production")
	if err != nil {
		log.Printf("ValidateFlow erro: %v", err)
	} else {
		prettyPrint("Flow resultado", flowResult)
	}

	// --- Validar Cenário ---
	scenarioResult, err := client.ValidateScenario(ctx, "meu-cenario", models.ValidateBody{
		Identifier: ptr("user-123"),
		Context:    map[string]any{"plano": "free"},
	}, "")
	if err != nil {
		log.Printf("ValidateScenario erro: %v", err)
	} else {
		prettyPrint("Cenário resultado", scenarioResult)
	}

	// --- Validar parte de Flow ---
	flowPartResult, err := client.ValidateFlowPart(ctx, "meu-fluxo-parte", models.ValidateBody{
		Identifier: ptr("user-123"),
	}, "")
	if err != nil {
		log.Printf("ValidateFlowPart erro: %v", err)
	} else {
		prettyPrint("FlowPart resultado", flowPartResult)
	}

	// --- Cache do flag ---
	cached, err := client.GetFlagCache(ctx, "minha-feature", "user-123")
	if err != nil {
		log.Printf("GetFlagCache erro: %v", err)
	} else {
		prettyPrint("Cache do flag", cached)
	}
}

func ptr[T any](v T) *T { return &v }
