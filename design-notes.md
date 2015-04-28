# Design notes

This is probably out of date, but keeping it for posterity in case there's any
good ideas in here.

## Plugins

Plugins are any executable called Q-&lt;something&gt;. Plugins are installed by
calling q install &lt;plugin-name&gt; which calls the `commands` command on the
plugin executable, which should return metadata about what commands and events
the plugin supports.  If this metadata doesn't conflict with any existing
plugins, Q will merge it with the existing metadata and then call the `install`
command for the plugin, which should do any required setup for the plugin.

### Required meta-commands

* meta - return metadata listing supported commands and events
* install - perform any install steps
* uninstall - perform any uninstall steps
* help - print general help for the plugin

### Metadata Format

Metadata returned from the meta command must be returned as JSON.

description - (string) single sentence describing the plugin

commands - array of command objects

#### Commands

Commands are exposed on the base Q command as subcommands, so for example, if
your plugin exposes the "get" command, users will execute that command as `q
get`.  Commands should more or less follow English grammar.  So `q show notes`
rather than `q notes show`.

Q has a few standard commands that can be appended to by plugins - these are for
object creation.  The commands move, new, show, del, and tag are built-in and
should be reused where appropriate.

#### Command metadata format

```
command - (string) name of the command (e.g. "new")
subcommand - (string, optional) name of required subcommand (e.g. "todo")
help - (string) single line description of the command
subhelp - (string) single line description of the subcommand
```

The plugin will be expected to respond to the command help &lt;command&gt;
&lt;subcommand&gt;, so for example, a plugin that exposes the "new todo" command
will need to be able to respond to a `help new todo` command with appopriate
help output.

