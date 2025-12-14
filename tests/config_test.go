package tests

import (
	"fmt"
	"time"
)

type Config struct {
	AppName       string
	Version       string
	Debug         bool
	StartTime     time.Time
	StartDuration time.Duration

	AppTags                 []string
	AppTagVersions          []int
	AppDependenciesVersions map[string]int
	AppDependenciesTags     map[string]string

	Server   ServerConfig
	Database DatabaseConfig
	Users    []User
}

type ServerConfig struct {
	Host     string
	Port     int
	Timeouts TimeoutConfig
}

type TimeoutConfig struct {
	Read  string
	Write string
}

type DatabaseConfig struct {
	Driver         string
	MaxConnections int
	Replicas       []ReplicaConfig
}

type ReplicaConfig struct {
	Host string
	Port int
}
type User struct {
	Name  string
	Roles []string
	Meta  map[string]any
}

func assertConfig(config Config) error {
	err := assertConfigWithoutStructArrays(config)
	if err != nil {
		return err
	}

	// Users
	if len(config.Users) != 2 {
		return fmt.Errorf("Users length: got %d, want %d", len(config.Users), 2)
	}

	// users[0]
	if config.Users[0].Name != "alice" {
		return fmt.Errorf("Users[0].Name: got %q, want %q", config.Users[0].Name, "alice")
	}
	if len(config.Users[0].Roles) != 2 {
		return fmt.Errorf("Users[0].Roles length: got %d, want %d", len(config.Users[0].Roles), 2)
	}
	if config.Users[0].Roles[0] != "admin" {
		return fmt.Errorf("Users[0].Roles[0]: got %q, want %q", config.Users[0].Roles[0], "admin")
	}
	if config.Users[0].Roles[1] != "user" {
		return fmt.Errorf("Users[0].Roles[1]: got %q, want %q", config.Users[0].Roles[1], "user")
	}
	if config.Users[0].Meta == nil {
		return fmt.Errorf("Users[0].Meta: got nil, want non-nil")
	}
	if gotTeam, ok := config.Users[0].Meta["team"]; !ok || gotTeam != "core" {
		return fmt.Errorf("Users[0].Meta[team]: got (%v, ok=%v), want (%v, ok=true)", gotTeam, ok, "core")
	}
	if gotActive, ok := config.Users[0].Meta["active"]; !ok || gotActive != true {
		return fmt.Errorf("Users[0].Meta[active]: got (%v, ok=%v), want (%v, ok=true)", gotActive, ok, true)
	}

	// users[1]
	if config.Users[1].Name != "bob" {
		return fmt.Errorf("Users[1].Name: got %q, want %q", config.Users[1].Name, "bob")
	}
	if len(config.Users[1].Roles) != 1 {
		return fmt.Errorf("Users[1].Roles length: got %d, want %d", len(config.Users[1].Roles), 1)
	}
	if config.Users[1].Roles[0] != "user" {
		return fmt.Errorf("Users[1].Roles[0]: got %q, want %q", config.Users[1].Roles[0], "user")
	}
	if config.Users[1].Meta == nil {
		return fmt.Errorf("Users[1].Meta: got nil, want non-nil")
	}
	if gotTeam, ok := config.Users[1].Meta["team"]; !ok || gotTeam != "ops" {
		return fmt.Errorf("Users[1].Meta[team]: got (%v, ok=%v), want (%v, ok=true)", gotTeam, ok, "ops")
	}
	if gotActive, ok := config.Users[1].Meta["active"]; !ok || gotActive != false {
		return fmt.Errorf("Users[1].Meta[active]: got (%v, ok=%v), want (%v, ok=true)", gotActive, ok, false)
	}

	// Database replicas
	if len(config.Database.Replicas) != 2 {
		return fmt.Errorf("Database.Replicas length: got %d, want %d", len(config.Database.Replicas), 2)
	}
	if config.Database.Replicas[0].Host != "db-replica-1" {
		return fmt.Errorf("Database.Replicas[0].Host: got %q, want %q", config.Database.Replicas[0].Host, "db-replica-1")
	}
	if config.Database.Replicas[0].Port != 5432 {
		return fmt.Errorf("Database.Replicas[0].Port: got %d, want %d", config.Database.Replicas[0].Port, 5432)
	}
	if config.Database.Replicas[1].Host != "db-replica-2" {
		return fmt.Errorf("Database.Replicas[1].Host: got %q, want %q", config.Database.Replicas[1].Host, "db-replica-2")
	}
	if config.Database.Replicas[1].Port != 5432 {
		return fmt.Errorf("Database.Replicas[1].Port: got %d, want %d", config.Database.Replicas[1].Port, 5432)
	}

	return nil
}

func assertConfigWithoutStructArrays(config Config) error {
	// Top-level fields
	if config.AppName != "ExampleApp" {
		return fmt.Errorf("AppName: got %q, want %q", config.AppName, "ExampleApp")
	}
	if config.Version != "1.0.0" {
		return fmt.Errorf("Version: got %q, want %q", config.Version, "1.0.0")
	}
	if config.Debug != true {
		return fmt.Errorf("Debug: got %v, want %v", config.Debug, true)
	}
	wantStartTime, err := time.Parse(time.RFC3339, "2025-01-01T12:30:45Z")
	if err != nil {
		return fmt.Errorf("parse want start time: %v", err)
	}
	if !config.StartTime.Equal(wantStartTime) {
		return fmt.Errorf("StartTime: got %s, want %s", config.StartTime.UTC().Format(time.RFC3339Nano), wantStartTime.UTC().Format(time.RFC3339Nano))
	}
	if config.StartDuration != 90*time.Minute {
		return fmt.Errorf("StartDuration: got %v, want %v", config.StartDuration, 90*time.Minute)
	}

	// Server
	if config.Server.Host != "0.0.0.0" {
		return fmt.Errorf("Server.Host: got %q, want %q", config.Server.Host, "0.0.0.0")
	}
	if config.Server.Port != 8080 {
		return fmt.Errorf("Server.Port: got %d, want %d", config.Server.Port, 8080)
	}
	if config.Server.Timeouts.Read != "5s" {
		return fmt.Errorf("Server.Timeouts.Read: got %q, want %q", config.Server.Timeouts.Read, "5s")
	}
	if config.Server.Timeouts.Write != "10s" {
		return fmt.Errorf("Server.Timeouts.Write: got %q, want %q", config.Server.Timeouts.Write, "10s")
	}

	// Database
	if config.Database.Driver != "postgres" {
		return fmt.Errorf("Database.Driver: got %q, want %q", config.Database.Driver, "postgres")
	}
	if config.Database.MaxConnections != 20 {
		return fmt.Errorf("Database.MaxConnections: got %d, want %d", config.Database.MaxConnections, 20)
	}
	return nil
}
