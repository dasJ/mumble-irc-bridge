with import <nixpkgs> {};
pkgs.buildGoModule {
  name = "mumble-irc-bridge";
  src = ./.;
  vendorSha256 = "sha256-wfciT8/1ODHPfSgQhAExJSEMQypfFiGSykF1o/KwNa4=";
}
