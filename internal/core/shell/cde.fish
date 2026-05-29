function cde
    set -x PATH $PATH $HOME/.nix-profile/bin
    if test (count $argv) -gt 0
        cde-bin $argv
    else
        set target (cde-bin path) or return 0
        builtin cd $target
    end
end
