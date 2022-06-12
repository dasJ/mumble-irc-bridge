with import <nixpkgs> {};
pkgs.buildGoModule {
  name = "mumble-irc-bridge";
  src = ./.;
  vendorSha256 = "sha256-Ptettm3Z6VSjuatD56dVU9Gvg0liIx1v3PFrL3oOiwc=";
  proxyVendor = true;
}
