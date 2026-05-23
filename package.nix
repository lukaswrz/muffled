{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "muffled";
  version = "0.0.0";

  src = lib.cleanSource ./.;

  vendorHash = "sha256-+/XSQWXx+SRoes7x1Gqj3J9z19CpQj05OFqv63ZP5yo=";

  meta = {
    description = "A ListenBrainz widget";
    license = lib.licenses.agpl3Only;
    mainProgram = "muffled";
  };
}
