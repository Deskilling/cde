{
  description = "cde";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    # generates outputs for all common systems
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "cde-bin";
          version = "0.0.7";
          src = ./.;

          vendorHash = "sha256-nee7yiwXsSfklE478gD5biuQjoL/d9X4N6fnnx3L1OY=";

          subPackages = [ "cmd/cde-bin" ];

          env.CGO_ENABLED = "0";

          flags = [ "-trimpath" ];
          ldflags = [
            "-s"
            "-w"
          ];
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            gnumake
            go
            gopls
          ];
        };
      }
    );
}
