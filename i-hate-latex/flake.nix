{
  description = "Nix flake generating the PhD Manuscript";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };

        buildInputs = with pkgs; [
          coreutils
          fira-code
          fira-code-symbols
          fira-code-nerdfont
          fontconfig
          texlive.combined.scheme-full
          which
          python311Packages.pygments
        ];

        # Environment Variables used in both packages and devShell
        TEXMFHOME = ".cache";
        TEXMFVAR = ".cache/texmf-var";
        SOURCE_DATE_EPOCH = toString self.lastModified;
        OSFONTDIR = "${pkgs.fira-code}/share/fonts";
      in
      with pkgs;
      {
        # formatter: Specify the formatter that will be used by the command `nix fmt`.
        formatter = nixpkgs-fmt;

        # The default devShell is activated when running `nix develop` in the directory containing this flake.
        # It instanciate a shell where the pkg defined in `buildInputs` are loaded and the env var set.
        devShells.default = mkShell {
          inherit buildInputs TEXMFHOME TEXMFVAR SOURCE_DATE_EPOCH OSFONTDIR;
        };

        packages.default = stdenvNoCC.mkDerivation {
          inherit buildInputs TEXMFHOME TEXMFVAR SOURCE_DATE_EPOCH OSFONTDIR;
          name = "document";
          src = self;
          phases = [ "unpackPhase" "buildPhase" "installPhase"];

          buildPhase = ''
            runHook preBuild
            latexmk -lualatex
            runHook postBuild
          '';

          installPhase = ''
            runHook preInstall
            install -m644 -D *.pdf $out/main.pdf
            runHook postInstall
          '';
        };

      });
}