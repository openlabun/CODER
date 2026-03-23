package roble_infrastructure

import (
	"fmt"
	"os"
	"strings"
)



type RobleDatabaseAdapter struct {
	client      *RobleClient
	accessToken string
}

func NewRobleDatabaseAdapter(client *RobleClient) *RobleDatabaseAdapter {
	return &RobleDatabaseAdapter{
		client:      client,
		accessToken: strings.TrimSpace(os.Getenv("ROBLE_ACCESS_TOKEN")),
	}
}

func NewRobleDatabaseAdapterWithToken(client *RobleClient, accessToken string) *RobleDatabaseAdapter {
	return &RobleDatabaseAdapter{client: client, accessToken: strings.TrimSpace(accessToken)}
}

func (a *RobleDatabaseAdapter) SetAccessToken(accessToken string) {
	a.accessToken = strings.TrimSpace(accessToken)
}

func (a *RobleDatabaseAdapter) Insert(tableName string, records []map[string]any) (map[string]any, error) {
	token, err := a.requireAccessToken()
	if err != nil {
		return nil, err
	}

	return a.client.Insert(tableName, records, token)
}

func (a *RobleDatabaseAdapter) Read(tableName string, conditions map[string]string) (map[string]any, error) {
	token, err := a.requireAccessToken()
	if err != nil {
		return nil, err
	}

	return a.client.Read(tableName, conditions, token)
}

func (a *RobleDatabaseAdapter) Update(tableName, idColumn, idValue string, updates map[string]any) (map[string]any, error) {
	token, err := a.requireAccessToken()
	if err != nil {
		return nil, err
	}

	return a.client.Update(tableName, idColumn, idValue, updates, token)
}

func (a *RobleDatabaseAdapter) Delete(tableName, idColumn, idValue string) (map[string]any, error) {
	token, err := a.requireAccessToken()
	if err != nil {
		return nil, err
	}

	return a.client.Delete(tableName, idColumn, idValue, token)
}

func (a *RobleDatabaseAdapter) GetClient() *RobleClient {
	return a.client
}

func (a *RobleDatabaseAdapter) requireAccessToken() (string, error) {
	if strings.TrimSpace(a.accessToken) == "" {
		return "", fmt.Errorf("ROBLE_ACCESS_TOKEN is required for database requests")
	}

	return a.accessToken, nil
}
