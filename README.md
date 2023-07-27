## githublink

Opens a file that is tracked in git (and stored on github.com) in the webbrowser

    usage: githublink <path> [<line-number>]

## Building

First download earthly, then run one of the corresponding targets which matches your platform:

    earthly +githublink-linux-amd64
    earthly +githublink-linux-arm64
    earthly +githublink-darwin-amd64
    earthly +githublink-darwin-arm64

This will output a binary under `./build/...`.

## Vim integration

Place the following under your `~/.vimrc` (or `~/.config/nvim/init.vim`):

```VimL
let g:go_fmt_command = "gofmt"
command! GitHubLink call GitHubLink()
function! GitHubLink()
    execute "!" . "githublink" . " " . bufname("%") . " " . ( line(".") ) ." &"
```

## Licensing
githublink is licensed under the Mozilla Public License Version 2.0. See [LICENSE](LICENSE).
