let pkgs = import <nixpkgs> {};
in with pkgs;
stdenv.mkDerivation {
	name = "goagainst";

	buildInputs = [ go ];
}
