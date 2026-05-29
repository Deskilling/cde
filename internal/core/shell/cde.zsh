cde() {
 export PATH="$PATH:$HOME/.nix-profile/bin"

 if [[ $# -gt 0 ]]; then
  cde-bin "$@"
 else
  local target
  target=$(cde-bin path) || return 0
  builtin cd "$target"
 fi
}
