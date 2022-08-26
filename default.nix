with import <nixpkgs> {};
pkgs.buildGoModule {
  name = "mumble-irc-bridge";
  src = ./.;
  vendorSha256 = "sha256-jCHs2b2aPou/5Y1ZcfLqudkrxBVlrNofDa2U5Qrqe1M=";
  proxyVendor = true;
}
