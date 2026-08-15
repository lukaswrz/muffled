{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "muffled";
  version = "0.0.0";

  src = lib.cleanSource ./.;

  vendorHash = "sha256-cWUDkXLZLhCsJ4X9v2XMJS6x6K9M0JrTKx6ZiVuGcwk=";

  meta = {
    description = "A ListenBrainz widget";
    license = lib.licenses.agpl3Only;
    mainProgram = "muffled";
  };
}
