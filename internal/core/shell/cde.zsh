cde() {
  local target
  export PATH="$PATH:$HOME/.nix-profile/bin"

  target=$(cde-bin) || return 0
  cd "$target"
}
