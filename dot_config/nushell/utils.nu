def confirm [message: string = ""] {
  let ans = input $"($message) Are you sure? [y/N] "
  return ($ans in [Y y])
}

def "docker-compose" [] {
  "Docker Compose Utilities"
}

def "docker-compose up-log" [service?: string] {
  match $service {
    null => { docker compose up -d; docker compose logs -f }
    _ => { docker compose up -d $service; docker compose logs -f $service }
  }
}

def "docker-compose update" [service?: string] {
  if not (confirm "This will restart containers, continue?") {
    return
  }
 
  match $service {
    null => { docker compose pull; docker compose down }
    _ => { docker compose pull $service; docker compose down $service }
  }
}

def "docker-compose clean-restart" [service?: string] {
  if not (confirm "This will restart containers, continue?") {
    return
  }

  match $service {
    null => { docker compose down; dc up-log }
    _ => { docker compose down $service; dc up-log $service }
  }
}

def "docker-build" [] {
  "Docker Build Utilities"
}

def "docker-build reg-build" [
  name_tag: string
  docker_context: string = "."
] {
  docker build -t $"registry.i.zelak.dev/($name_tag)" $"($docker_context)"
}

def "docker-build reg-push" [
  name_tag: string
] {
  docker push $"registry.i.zelak.dev/($name_tag)"
}

def --env "cd ls" [
  path: string
] {
  cd $path
  ls
}

def --env mkcd [
  path: string
] {
  mkdir $path
  cd $path
}

def backup [
  --name (-n): string
  ...files: string
] {
  if ($files | length) < 1 {
    print "At least one file needs to be provided"
    return
  }

  let b_name = match $name {
    null => { $"(date now | format date "%Y%m%d_%H%M%S").tar.gz" }
    _ => {
      if ($name | str ends-with ".tar.gz") {
        $name
      } else {
        $"($name).tar.gz"
      }
    }
  }

  let cmd_tar = (tar -czf $b_name ...$files | complete)
  if ($cmd_tar | get exit_code) == 1 {
    print "Backup failed"
    return
  }

  print $"Backup created: ($b_name)"
}

def "check ports" [
  range: string
] {
  ss -tuln 
  | lines 
  | str trim 
  | split column -c " " 
  | rename ...($in | get 0 | values) 
  | skip 1 
  | where ($it.Local | split row ":" | last | str starts-with $range)
}

def "git chore" [
  branch_name: string
] {
  git switch -c $"chore/($branch_name)"
}

def "git feature" [
  branch_name: string
] {
  git switch -c $"feature/($branch_name)"
}

def "git stash-diff" [
  stash_number?: int
] {
  match $stash_number {
    null => { git stash show -p }
    _ => { git stash show -p $"stash@{($stash_number)}" }
  }
}

def --env "git clone" [
  --no-cd (-n)
  --base (-b): string
  url: string
  ...git_args: string
] {
  let root_base = $"($nu.home-path)/git"

  let parsed_path = $url | path parse
  mut repo_owner = $parsed_path | get parent | path basename
  let repo_name = $parsed_path | get stem
  
  if ($repo_owner | str contains ":") {
    # Git ssh url maybe
    $repo_owner = $repo_owner | split row ":" | last
  }

  mut r_base = $root_base
  if $base != null {
    $r_base = $root_base | path join $base
  } else {
    $r_base = $root_base | path join $repo_owner
  }
  
  $r_base = $r_base | path join $repo_name
  mkdir $r_base
  ^git clone $url $r_base ...$git_args

  if not $no_cd {
    export-env {
      cd $r_base
    }
  }
}

def "zip extract" [
  file_path: string
] {
  if not ($file_path | path exists) {
    print $"path: ($file_path) doesn't exist"
    return
  }

  let base_name = $file_path | path basename | split row "." | $in.0
  let extension = $file_path | path basename | split row "." | if ($in | length) > 2 { select 1 2 | str join "." } else { $in.1 }
  let command = match $extension {
    "tar.bz2" | "tbz2" => "tar xjf"
    "tar.gz" | "tgz" => "tar xzf"
    "tar.xz" | "tar" => "tar xf"
    "bz2" => "bunzip2"
    "gz" => "gunzip"
    "rar" =>  "unrar x"
    "zip" => "unzip"
    "Z" => "uncompress"
    "7z" => "7z x"
    "deb" => "ar x"
    "tar.zst" => "unzstd"
  }

  let output = ^($command | split row " " | append $file_path) | complete
  if $output.exit_code == 0 {
    print "Extraction done"
  } else {
    print "Extraction failed"
    print $output.stderr 
  }
}

def --env "path save" [
  name: string = "default"
] {
  $env.SAVE_DIRS = $env.SAVE_DIRS? | default {} | upsert $name (pwd)
  print $"Saved current directory as ($name)"
}

def --env "path restore" [
  name: string = "default"
] {
  let path = $env | get --optional ([SAVE_DIRS, $name] | into cell-path)
  if $path == null {
    print $"No saved directory found for ($name)"
    return
  }

  cd $path
}

def "path list" [] {
  $env.SAVE_DIRS?
}
