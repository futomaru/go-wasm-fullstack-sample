package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// add.wasm をバイナリとして埋め込む例
//
//go:embed add.wasm
var addWasm []byte

func main() {
	flag.Parse()

	ctx := context.Background()
	runtime := wazero.NewRuntime(ctx)
	defer runtime.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	module, err := runtime.InstantiateWithConfig(ctx, addWasm, wazero.NewModuleConfig().WithStartFunctions("_initialize"))
	if err != nil {
		log.Fatalf("failed to instantiate module: %v", err)
	}
	defer module.Close(ctx)

	a, b := operands()

	add := module.ExportedFunction("Add")
	if add == nil {
		log.Fatal("Add export not found")
	}

	results, err := add.Call(ctx, uint64(uint32(a)), uint64(uint32(b)))
	if err != nil {
		log.Fatalf("failed to call Add: %v", err)
	}

	fmt.Printf("%d + %d = %d\n", a, b, int32(uint32(results[0])))
}

func operands() (int32, int32) {
	if flag.NArg() < 2 {
		return 1, 41
	}
	return parseOrDefault(flag.Arg(0), 1), parseOrDefault(flag.Arg(1), 41)
}

func parseOrDefault(v string, fallback int32) int32 {
	if v == "" {
		return fallback
	}
	if n, err := strconv.ParseInt(v, 10, 32); err == nil {
		return int32(n)
	}
	return fallback
}
