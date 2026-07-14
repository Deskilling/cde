function cde
    if test (count $argv) -gt 0
        if string match -q -- '^-' "$argv[1]"
            cd (cde-bin $argv[1])
        else
            cd (cde-bin $argv)
        end
    else
        cd (cde-bin path)
    end
end
