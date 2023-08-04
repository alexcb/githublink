## githublink

Opens a file that is tracked in git (and stored on github.com) in the webbrowser

    usage: githublink <path> [<line-number>]

## Demo

Check out the following 30 second demo video to see it in action:

[![Demo Video](https://img.youtube.com/vi/1kcVQ0uegRM/0.jpg)](https://www.youtube.com/watch?v=1kcVQ0uegRM)

## Installation

intel-linux users:

    sudo /bin/sh -c ' wget https://github.com/alexcb/githublink/releases/latest/download/githublink-linux-amd64 -O /usr/local/bin/githublink && chmod +x /usr/local/bin/githublink'

raspberrypi-v4-linux users:

    sudo /bin/sh -c ' wget https://github.com/alexcb/githublink/releases/latest/download/githublink-linux-arm64 -O /usr/local/bin/githublink && chmod +x /usr/local/bin/githublink'

intel-mac users:

    sudo /bin/sh -c ' wget https://github.com/alexcb/githublink/releases/latest/download/githublink-darwin-amd64 -O /usr/local/bin/githublink && chmod +x /usr/local/bin/githublink'

m1/2-mac users:

    sudo /bin/sh -c ' wget https://github.com/alexcb/githublink/releases/latest/download/githublink-darwin-arm64 -O /usr/local/bin/githublink && chmod +x /usr/local/bin/githublink'


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
    execute "!" . "githublink" . " " . bufname("%") . " " . ( line(".") )
```

## Licensing
githublink is licensed under the Mozilla Public License Version 2.0. See [LICENSE](LICENSE).
