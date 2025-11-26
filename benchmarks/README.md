# TinyJSON Benchmarks

This directory contains benchmarking tools to compare TinyJSON against the standard library `encoding/json` when compiled to WebAssembly using TinyGo.

## Build Script

The `build.sh` script compiles the WASM binaries with different JSON implementations.

### Usage

```bash
# Build with TinyJSON (default)
./build.sh

# Build with encoding/json (stdlib)
./build.sh stlib
```

# Run the benchmark server

```bash
# Start the local server to serve the compiled WASM
go run ./web/server.go
```

You can then open `http://localhost:6060` in a browser to see the benchmark results.

### Output

The compiled WASM binary is output to `web/public/main.wasm`.

The script will display:
- Uncompressed binary size
- Gzipped size (for realistic deployment comparison)

### Example Output

```
Building with TinyJSON...
Compiling clients/tinyjson/main.go...
✓ Build complete: web/public/main.wasm
  Size: 94 KB
  Gzipped: 34 KB

Building with encoding/json (stdlib)...
Compiling clients/stdlib/main.go...
✓ Build complete: web/public/main.wasm
  Size: 398 KB
  Gzipped: 150 KB
```


## Source Files

- `clients/tinyjson/main.go` - Implementation using TinyJSON
- `clients/stdlib/main.go` - Implementation using encoding/json (stdlib)

Both files implement the same functionality to ensure fair comparison.

## Results

| Implementation | Binary Size (WASM + Gzip) |
| :--- | :--- |
| **TinyJSON** | **27.2 KB** |
| encoding/json (stdlib) | 119 KB |

See the [main README](../README.md#benchmarks) for detailed benchmark results and screenshots.
