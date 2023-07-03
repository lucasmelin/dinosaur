<h1 align="center"> dinosaur </h1>

<p align="center">
A toy DNS resolver, inspired by <a href="https://implement-dns.wizardzines.com/index.html">Implement DNS in a weekend</a> by Julia Evans.
<br/>
<br/>
<br/>
<img src="assets/logo.png" alt="dinosaur logo">
</p>

## Usage

![Dinosaur demo video](assets/demo.mp4)

## Development

[Install Nix on your system](https://github.com/DeterminateSystems/nix-installer), then run `nix build` to build the package.
The resulting application can be run from `result/bin/dinosaur`.
With `direnv` installed, you can run `direnv allow` to set up a new shell containing all the development tools.

### Adding dependencies

[Add the import](https://go.dev/blog/using-go-modules) to your `.go` file, then run `go mod tidy` to trigger an update of the `go.mod` file.
Update the `gomod2nix.toml` and Nix store using `gomod2nix generate`.
