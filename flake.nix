{
  description = "Terminal torrent search and download client";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/b5aa0fbd538984f6e3d201be0005b4463d8b09f8";

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
      pkgsFor = system: import nixpkgs { inherit system; };
      version = "0.2.0";
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
          pork = pkgs.callPackage ./packaging/nix/package.nix {
            inherit version;
            source = self;
          };
        in
        {
          inherit pork;
          default = pork;
        }
      );

      apps = forAllSystems (system: {
        default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/pork";
        };
      });

      devShells = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
        in
        {
          default = pkgs.mkShell {
            packages = [ pkgs.go ];
          };
        }
      );
    };
}
