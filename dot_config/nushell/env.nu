# env.nu
#
# Installed by:
# version = "0.106.1"
#
# Previously, environment variables were typically configured in `env.nu`.
# In general, most configuration can and should be performed in `config.nu`
# or one of the autoload directories.
#
# This file is generated for backwards compatibility for now.
# It is loaded before config.nu and login.nu
#
# See https://www.nushell.sh/book/configuration.html
#
# Also see `help config env` for more options.
#
# You can remove these comments if you want or leave
# them for future reference.

#$env.LANG = 'en_US.UTF-8'
#$env.LC_ALL = 'en_US.UTF-8'
$env.PATH ++= ['~/bin']

const empty_path = "~/.config/nushell/empty.nu"
const specific_path = "~/.config/nushell/env.specific.nu"
source (if ($specific_path | path exists) { $specific_path } else { $empty_path })
