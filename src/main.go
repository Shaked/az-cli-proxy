package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	http.HandleFunc("/token", tokenHandler)

	port := "8080"
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		log.Printf("Resource parameter is required\n")
		http.Error(w, "Resource parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("Getting token for resource: %s\n", resource)
	// Create an Azure CLI Credential
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		log.Printf("Failed to create Azure CLI Credential: %s\n", err.Error())
		http.Error(w, "Failed to create Azure CLI Credential: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetching the token
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{resource}})
	if err != nil {
		log.Printf("Failed to get token: %s\n", err.Error())
		http.Error(w, "Failed to get token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	expiresOn := strconv.FormatInt(token.ExpiresOn.Unix(), 10)
	response := map[string]string{
		"access_token": token.Token,
		"expires_on":   expiresOn,
		"resource":     resource,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %s\n", err.Error())
		http.Error(w, "Failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Token fetched successfully\n")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
