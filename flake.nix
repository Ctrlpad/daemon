{
  description = "Ctrlpad daemon - listens for Ctrlpad button actions over BLE and executes them";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      forAllSystems =
        f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      packages = forAllSystems (pkgs: rec {
        ctrlpad-daemon = pkgs.callPackage ./nix/package.nix {
          version = "0.0.7";
          commit = self.rev or self.dirtyRev or "none";
          date = self.lastModifiedDate;
        };
        default = ctrlpad-daemon;
      });

      overlays.default = final: _prev: {
        ctrlpad-daemon = final.callPackage ./nix/package.nix { version = "0.0.7"; };
      };

      nixosModules = rec {
        ctrlpad-daemon = import ./nix/module.nix self;
        default = ctrlpad-daemon;
      };

      formatter = forAllSystems (pkgs: pkgs.nixfmt-tree);
    };
}
