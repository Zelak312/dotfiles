# nushell install
# maybe needed but seems to already by
# installed if on host (will see)
#echo "[gemfury-nushell]
#name=Gemfury Nushell Repo
#baseurl=https://yum.fury.io/nushell/
#enabled=1
#gpgcheck=0
#gpgkey=https://yum.fury.io/nushell/gpg.key" | sudo tee /etc/yum.repos.d/fury-nushell.repo
#sudo dnf install -y nushell

# vim install
dnf install -y vim

# starship install
curl -sS https://starship.rs/install.sh | sh -s -- -y

# Chezmoi install (as user zelak)
sh -c "$(curl -fsLS get.chezmoi.io)" -- init --apply Zelak312 --branch nushell
chown -R zelak:zelak $HOME
