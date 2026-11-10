{
  description = "template";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
  };

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      packages = forAllSystems (pkgs: {
        default = pkgs.buildGoModule {
          pname = "cde-bin";
          version = "0.1.1";
          src = ./.;

          vendorHash = "sha256-nee7yiwXsSfklE478gD5biuQjoL/d9X4N6fnnx3L1OY=";

          subPackages = [ "cmd/cde-bin" ];

          env.CGO_ENABLED = "0";

          ldflags = [
            "-s"
            "-w"
          ];
        };
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShellNoCC {
          buildInputs = with pkgs; [
            gnumake
            go
            gopls
            golangci-lint
          ];
        };
      });
    };
}
