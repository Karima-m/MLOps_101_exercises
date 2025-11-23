package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

func main() {
	// Create a shared context
	ctx := context.Background()

	// Run the stages of the pipeline
	if err := Build(ctx); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
}

func Build(ctx context.Context) error {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx) // Connects to dagger engine
	if err != nil {
		return err
	}
	defer client.Close()

	// Pulls container from docker
	python := client.Container().From("python:3.12.2-bookworm").
		// copies local local "python-files" folder into container at path "python"
		WithDirectory("python", client.Host().Directory("python-files")).
		// executes python script inside the container
		WithExec([]string{"python", "--version"})

	python = python.WithExec([]string{"python", "python/hello.py"})

	_, err = python.
		// exports container's "output" folder back to host filesystem (GitHub actions)
		Directory("output").
		Export(ctx, "output")
	if err != nil {
		return err
	}

	return nil
}
