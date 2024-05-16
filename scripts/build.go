package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func main() {

	// Comenzar a medir el tiempo
	startTime := time.Now()

	// Establecer las variables de entorno necesarias
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")
	os.Setenv("CGO_ENABLED", "0")

	// Obtener todos los directorios dentro de "lambdas/"
	lambdaFolders, err := filepath.Glob("lambdas/*")
	if err != nil {
		log.Println("Error listing lambda folders:", err)
		return
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(lambdaFolders))

	// Compilar cada lambda en paralelo
	for _, folder := range lambdaFolders {
		wg.Add(1)
		go func(folder string) {
			defer wg.Done()
			err := buildLambda(folder)
			if err != nil {
				errors <- err
			}
		}(folder)
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()
	close(errors)

	// Manejar errores después de que todas las compilaciones han terminado
	for err := range errors {
		if err != nil {
			log.Println("Error:", err)
		}
	}

	// Calcular el tiempo total de ejecución
	duration := time.Since(startTime)
	log.Printf("Total build time: %v\n", duration)
}

func buildLambda(folder string) error {
	folderName := filepath.Base(folder)
	lambdaPath := filepath.Join("lambdas", folderName)

	// Build the lambda executable
	cmd := exec.Command("go", "build", "-tags", "lambda.norpc", "-o","bootstrap")
    cmd.Dir = lambdaPath  // Set the working directory for the command
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("Failed to build in %s: %s", lambdaPath, string(output))
        return err
    }

	// Ensure the bootstrap file path is correct for zipping
    bootstrapPath := filepath.Join(lambdaPath, "bootstrap")

	// Empaquetar el ejecutable en un archivo zip
    zipPath := filepath.Join(lambdaPath, "../../bin", folderName+".zip")
    cmd = exec.Command("zip", "-j", zipPath, bootstrapPath)
    if output, err := cmd.CombinedOutput(); err != nil {
        log.Printf("Failed to zip bootstrap in %s: %s", lambdaPath, string(output))
        return err
    }

	// Eliminar el archivo bootstrap
    os.Remove(bootstrapPath)

	return nil
}