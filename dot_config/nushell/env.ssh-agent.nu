# deals with ssh agent
do --env {
    let ssh_agent_file = (
        $nu.temp-path | path join $"ssh-agent-(whoami).nuon"
    )

    if ($ssh_agent_file | path exists) {
        let ssh_agent_env = open ($ssh_agent_file)
        if ($"/proc/($ssh_agent_env.SSH_AGENT_PID)" | path exists) {
            load-env $ssh_agent_env
            return
        } else {
            rm $ssh_agent_file
        }
    }

    let ssh_agent_env = ^ssh-agent -c
        | lines
        | first 2
        | parse "setenv {name} {value};"
        | transpose --header-row
        | into record
    load-env $ssh_agent_env
    $ssh_agent_env | save --force $ssh_agent_file

    let ssh_keys = [
        "~/.ssh/id_ed25519"
    ]

    for key in $ssh_keys {
        let expanded_key = ($key | path expand)
        if ($expanded_key | path exists) {
            let already_loaded = (ssh-add -l | complete | get stdout | str contains ($expanded_key | path basename))
            if not $already_loaded {
                ssh-add $expanded_key
            }
        }
    }
}
