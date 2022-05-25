with import <nixpkgs> {};
pkgs.buildGoModule {
  name = "mumble-irc-bridge";
  src = ./.;
  vendorSha256 = "sha256-AwJbfWoK8vUvzLF9+zS4rBZ6K9sdSBCqcfKpWDNlwPs=";
  proxyVendor = true;
}
