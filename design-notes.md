## Design Notes:

## Plugin types:

### Commands and Contexts

Q command lines take the form `q <command> [context] [...]`  When a plugin is
installed, it can register commands and optionally new contexts.  Commands
are the verbs of the command line, for example "list" or "add".  Commands often
take a context, which specializes the command for this particular plugin.
Examples of contexts might be bugs or todos, or even something like git.  Put
together, a todos plugin might register "add", "remove", and "list" commands,
each with the context "todos".  With that plugin installed, a user could type `q
list todos` and get back a list of todo items.  If that user also had a bugs
plugin installed, they might also be able to do `q list bugs` to get back a list
of bugs.

Plugins may register commands without associated contexts, but this means that
no other plugins can reuse that command (since they would conflict), so this
should be used sparingly.

If two plugins would register the same context, during the install of the second
plugin, the install would fail, citing conflicting contexts.  In the future, it
will be possible to rename contexts so that they do not conflict.


### Servers

Servers allow other plugins to manipulate information through their service.
Often times this will take the form of API calls to a web service, such as
github or a kanban board.


## Creating a plugin

Q plugins are simply executables. They may be written in any language that can
read from stdin and write to stdout.  They must respond to the following command
line commands:

### manifest

When given the manifest command, the plugin must write its TOML manifest to
stdout.  The manifest tells Q what features the plugin exposes - commands and
their associated contexts, and any named services.

```toml
	name = "launchpad"
	# commands at the root are commands that do not take a context
    [[command]]
    	name = "dance"
    	short = "makes a dancing banana appear"
    	long = """
usage:
	dance [options]

	-t=<seconds> 	display for n seconds
"""
    [[context]]
        name = "bug"
        plural = "bugs"
    	[[context.command]]
    		name = "list"
    		short = "show all bugs"
    		long = """
usage:
	list bugs [options]

	-a=<assignee> 	show only bugs assigned to assignee
"""

    	[[context.command]]
    		name = "add"
    		short = "create a new bugs"
    		long = """
usage:
	add bug [options]

	-t=<title> 	use the given string as the title of the bug
	-b=<body> 	use the given string as the body of the bug
"""


```



## Use Cases:

### Work on a Bug

I want to tell Q about a bug that I am going to start working on.  

**Q should:**

- connect to the bug tracker and assign the bug to me
- make a branch of the code in source control with an obvious name
- create a card on my kanban board with the appropriate type and a link to the bug

**example CLI:**

q take [bug url] [vcs dir]

q take lp:1424892 gh:juju/juju

**implications:**

Need a string resolvers that can translate custom uri types to full strings, i.e. translate lp:1424892 to https://bugs.launchpad.net/juju-core/+bug/1424892 or translate gh:juju/juju to $GOPATH/src/github.com/juju/juju

Need plugins to handle vcs, bug trackers, kanban.
