
# Re-create tmux sessions

<br>

### To init a new configuration file
`tmuxer init <session-name> [path]`
<br>
```toml
# tmuxer basic config file
name = "tmux session name"
root = "root directory of the tmux session"


# To declare a window
[[window]] # All fields are optional
name = "editor"
layout = "horizontal" # `horizontal`, and `vertical` for now
commands = ["nvim",""] # Each command will run for each pane

[[window]] # This is valid

[[window]]
...
```

### To start in the directory with the config file
`tmuxer` 


### Start session and specify config directory
`tmuxer [config-path]`
