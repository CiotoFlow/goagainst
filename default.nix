with import <nixpkgs> {};
let
  mkDep = name: goRoot: src:
    stdenv.mkDerivation rec {
      inherit name src;
      
      buildCommand = ''
        ensureDir $out/src/`dirname ${goRoot}`
        ln -s ${src} $out/src/`dirname ${goRoot}/tiedot`
      '';
    };
    
  mkGithubDep = name: attrs:
    mkDep name "github.com/${attrs.owner}/${attrs.repo}" (fetchFromGitHub attrs);

  mkBitbucketDep = name: attrs:
    mkDep name "bitbucket.org/kardianos/osext" (fetchzip {
        name = "${name}-${attrs.rev}-src";
        url = "https://bitbucket.org/${attrs.owner}/${attrs.repo}/get/${attrs.rev}.zip";
        inherit (attrs) sha256;
    });

  deps = [
    
    (mkGithubDep "go-rice" {
      owner = "GeertJohan";
      repo = "go.rice";
      rev = "a4d0b5624c673fef4b517f350272136ced6bb5b1";
      sha256 = "0bksi0w08hc997wxqqw90nwpz1y3cf94bdiy2kgq8yvw3pq4xzc7";
    })
  
    (mkGithubDep "tiedot" {
      owner = "HouzuoGuo";
      repo = "tiedot";
      rev = "33b05aa560b2f8d40b9285293e3ba8f8b3a22dff";
      sha256 = "1l4ra6jvg1l2szk8zziy87q4x903gd34vabiwbpydyhf4a2m2js7";
    })

    (mkBitbucketDep "osext" {
      owner = "kardianos";
      repo = "osext";
      rev = "5d3ddcf53a50";
      sha256 = "0ggh93kbvs9xwgrabfzix1gsj1fb97ckf4llryrnp5bncrwl0gh5";
    })

    (mkGithubDep "go-zipexe" {
      owner = "daaku";
      repo = "go.zipexe";
      rev = "44882fc939f4c58d87a60de34796c6cfb9623269";
      sha256 = "1z6ivb1pla65hn1wz49paigx15p417xfipnnw137yn7hg90f7y8b";
    })
    
  ];

in
stdenv.mkDerivation {
  name = "goagainst";

  buildInputs = [ go ];

  GOPATH = lib.makeLibraryPath deps;
  
  buildPhase = ''
    go build -o goagainst ./src
  '';

  checkPhase = ''
    go test ./src/trollan
  '';
}
