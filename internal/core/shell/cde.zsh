cde() {
  local target
  export PATH="$PATH:$HOME/.nix-profile/bin"

  target=$(/usr/local/bin/cde-bin) || return 0
  cd "$target"
}
