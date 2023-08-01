## Welcome to my dotfiles

There isn't much to see here

# Installing

I like to make symlinks from the repo to the home directory because it's easy to update after with only git

(The one liners will download the latest release of the installer from github if you don't like this you can do it manually or build it yourself)

## If you haven't cloned the repo yet, here is a one liner

```bash
git clone https://github.com/Zelak312/dotfiles.git ~/dotfiles && cd ~/dotfiles && ARCH=$(uname -m); if [ "$ARCH" = "x86_64" ]; then ARCH="amd64"; else ARCH="arm64"; fi; URL=$(curl -s https://api.github.com/repos/Zelak312/dotfiles/releases/latest | jq -r --arg ARCH "$ARCH" '.assets[] | select(.name | contains("dotfiles-linux-"+$ARCH)) | .browser_download_url'); TEMP=$(mktemp) && curl -L -o $TEMP $URL && chmod +x $TEMP && $TEMP && rm $TEMP
```

## if you already cloned the repo, cd inside the root folder

```bash
ARCH=$(uname -m); if [ "$ARCH" = "x86_64" ]; then ARCH="amd64"; else ARCH="arm64"; fi; URL=$(curl -s https://api.github.com/repos/Zelak312/dotfiles/releases/latest | jq -r --arg ARCH "$ARCH" '.assets[] | select(.name | contains("dotfiles-linux-"+$ARCH)) | .browser_download_url'); TEMP=$(mktemp) && curl -L -o $TEMP $URL && chmod +x $TEMP && $TEMP && rm $TEMP
```
