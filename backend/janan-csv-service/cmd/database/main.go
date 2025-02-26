package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"janan_csv_service/internal/models"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
)

func loadEnums(sb *strings.Builder) {
	enums := []string{
		`CREATE TYPE gender AS ENUM (c
	  'MALE',
	  'FEMALE'
	 );`,
	}
	for _, enum := range enums {
		sb.WriteString(enum)
		sb.WriteString(";\n")
	}
	enums = []string{
		`CREATE TYPE language AS ENUM (
	  'en',
	  'fr',
	  'es'
	 );`,
	}
	for _, enum := range enums {
		sb.WriteString(enum)
		sb.WriteString(";\n")
	}
	enums = []string{
		`CREATE TYPE RestMethods AS ENUM (
	  'Get',
	  'Post'
	 );`,
	}
	for _, enum := range enums {
		sb.WriteString(enum)
		sb.WriteString(";\n")
	}
}

func loadModels(sb *strings.Builder) {
	models := []interface{}{
		&models.Student{}}
	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	sb.WriteString(stmts)
	sb.WriteString(";\n")
}

func main() {
	sb := &strings.Builder{}
	loadEnums(sb)
	loadModels(sb)

	_, err := io.WriteString(os.Stdout, sb.String())
	if err != nil {
		log.Fatal(err)
	}
}
