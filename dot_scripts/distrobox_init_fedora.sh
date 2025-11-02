# nushell install
echo "[gemfury-nushell]
name=Gemfury Nushell Repo
baseurl=https://yum.fury.io/nushell/
enabled=1
gpgcheck=0
gpgkey=https://yum.fury.io/nushell/gpg.key" | tee /etc/yum.repos.d/fury-nushell.repo
dnf remove -y nu
dnf install -y nushell

# git
dnf install -y git

# vim install
dnf install -y vim

# starship install
curl -sS https://starship.rs/install.sh | sh -s -- -y

# Chezmoi install (as user zelak)
sh -c "$(curl -fsLS get.chezmoi.io)" -- init --apply Zelak312 --branch nushell
chown -R zelak:zelak $HOME
