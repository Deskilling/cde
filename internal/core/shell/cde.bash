cde() {
    if [[ $# -gt 0 ]]; then
        if [[ "$1" == -* ]]; then
            target=$(cde-bin -${1#-}) || return 0
            builtin cd "$target"
        else
            cde-bin "$@"
    fi
    else
        local target
        target=$(cde-bin path) || return 0
        builtin cd "$target"
    fi
}
