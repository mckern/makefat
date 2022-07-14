package main_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMakeFat(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("works on darwin only")
	}

	// Make a directory to work in.
	dir, err := ioutil.TempDir("", "makefat")
	if err != nil {
		t.Fatalf("could not create directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// List files we're working with.
	src := filepath.Join(dir, "test.go")
	amd64 := filepath.Join(dir, "amd64")
	arm64 := filepath.Join(dir, "arm64")
	fat := filepath.Join(dir, "fat")

	// Create test source.
	f, err := os.Create(src)
	if err != nil {
		t.Fatalf("could not create source file: %v", err)
	}
	f.Write([]byte(`
package main
import "fmt"
func main() {
	fmt.Println("hello world")
}
`))
	f.Close()

	// Compile test code in both amd64 and arm64.
	cmd := exec.Command("go", "build", "-o", amd64, src)
	cmd.Env = append(os.Environ(), "GOARCH=amd64")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("could not build amd64 target: %v\n%s\n", err, string(out))
	}
	cmd = exec.Command("go", "build", "-o", arm64, src)
	cmd.Env = append(os.Environ(), "GOARCH=arm64")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("could not build arm64 target: %v\n%s\n", err, string(out))
	}

	// Build fat binary.
	cmd = exec.Command("go", "run", "../makefat.go", fat, amd64, arm64)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("could not build fat target: %v\n%s\n", err, string(out))
	}

	// Run fat binary.
	fmt.Printf("fat is %s\n", fat)
	cmd = exec.Command(fat)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("could not run fat target: %v", err)
	}
	if string(out) != "hello world\n" {
		t.Errorf("got=%s, want=hello world\n", string(out))
	}

	// Extract from fat
	cmd = exec.Command("go", "run", "extractfat.go", fat)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("could not extract fat: %v\n%s\n", err, string(out))
	}

	// Compare
	amd64Data, _ := ioutil.ReadFile(amd64)
	amd64Sum := sha256.Sum256(amd64Data)
	arm64Data, _ := ioutil.ReadFile(arm64)
	arm64Sum := sha256.Sum256(arm64Data)
	amd64EData, _ := ioutil.ReadFile(fat + ".CpuAmd64")
	amd64ESum := sha256.Sum256(amd64EData)
	arm64EData, _ := ioutil.ReadFile(fat + ".CpuArm64")
	arm64ESum := sha256.Sum256(arm64EData)
	if amd64Sum != amd64ESum || arm64Sum != arm64ESum {
		fmt.Printf("%s %s\n", hex.EncodeToString(amd64Sum[:]), amd64)
		fmt.Printf("%s %s\n", hex.EncodeToString(amd64ESum[:]), "amd64E")
		fmt.Printf("%s %s\n", hex.EncodeToString(arm64Sum[:]), arm64)
		fmt.Printf("%s %s\n", hex.EncodeToString(arm64ESum[:]), "arm64E")
		t.Fatal("extracted files mismatch")
	}
}
