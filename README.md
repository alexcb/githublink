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

## How is this different from `gh browse`?

1. No need to login.
2. It opens the URL based on the location of HEAD, rather than always pointing to the main branch. This allows you to share URLs to specific lines that won't change (i.e. break).
3. I'll accept PRs to expand this to work with gitlab or other repos.
4. I wrote my initial implementation in python in [2020](https://github.com/alexcb/acbutils/blob/18a4bec7b1f11182ce4ae7cc5f81a60cf9083098/scripts/githublink); `gh browse` was written in [2021](https://github.com/cli/cli/commit/68ce66801b5fb076e449d30c3dcb2867d7cd47b9), the command `githublink` is stuck in my muscle-memory.

## Licensing
githublink is licensed under the Mozilla Public License Version 2.0. See [LICENSE](LICENSE).
