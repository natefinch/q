# Design Notes

### Plugin types

#### Commands and Contexts

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


#### Servers

Servers allow other plugins to manipulate information through their service.
Often times this will take the form of API calls to a web service, such as
github or a kanban board.


### Creating a plugin

Q plugins are simply executables. They may be written in any language that can
read from stdin and write to stdout and expose/consume JSON-RPC.  They must
respond to the following command line commands:

#### manifest

When given the manifest command, the plugin must write its TOML manifest to
stdout.  The manifest tells Q what features the plugin exposes - commands and
their associated contexts, and any named services.

**Example Manifest**

	# the name of the plugin
	Name = "github"
	[[Command]]
		# commands at the root of the manifest are commands that do not take a context
		# the name of the command
		Name = "dance"
		# short is a single line help shown when help displays a list of commands
		Short = "makes a dancing banana appear"
		# long is what is shown when the user does q help <command>
		Long = """
		usage:
			dance [options]

			-t=<seconds> 	display for n seconds
		"""
	[[Context]]
		# the name of the context
	    Name = "bug"
	    # the plural version of the context name
	    Plural = "bugs"
	    Short = "These commands interact with bugs from github."
	    Long = """
	    Bug commands interact with bug from github.  You will need to set up 
	    authentication in the plugin's configuration.
	    """

	    # contexts can register commands as well, these will be accessed by
	    # q <command> <context>
		[[Context.Command]]
			Name = "list"
			Short = "show all bugs"
			Long = """
			usage:
				list bugs [options]

				-a=<assignee> 	show only bugs assigned to assignee
			"""

		[[Context.Command]]
			Name = "add"
			Short = "create a new bug"
			Long = """
			usage:
				add bug [options]

				-t=<title> 	use the given string as the title of the bug
				-b=<body> 	use the given string as the body of the bug
			"""
	[[Service]]
		# services expose functionality that other plugins can use.
		Name = "github"

**Manifest Definition**

	type Manifest struct {
		Name string 		// plugin name
		Command []Command 	// commands w/o context
		Context []Context 	// contexts 
	}

	type Command struct {
		Name string 	// command name
		Short string 	// single line help text
		Long string 	// multi-line help text
	}

	type Context struct {
		Name string 		// context name
		Plural string 		// plural version of context name
		Short string 		// single line help text
		Long string 		// multi-line help text
		Command []Command 	// commands that use this context
	}

### server

When given the server command, the plugin should start its JSON-RPC server
listening on stdin and writing responses to stdout.  Any logging the plugin
wishes to do should be sent through stderr.  The plugin should interpret the
Interrupt signal as a request for it to shut down its process.

### command [args...]

When given the command command, any further command line args given by the user
will be passed through to the plugin.  At this point, Q will start a JSON-RPC
server over the plugin's stdin and stdout, and will relay any RPC requests to
the appropriate plugin servers.  

## Installing plugins





## Example Use Cases:

### Work on a Bug

I want to tell Q about a bug that I am going to start working on.  

**Q should:**

- connect to the bug tracker and assign the bug to me
- make a branch of the code in source control with an obvious name
- create a card on my kanban board with the appropriate type and a link to the bug

**example CLI:**

q take [bug url] [vcs dir]

**implications:**

Need a string resolvers that can translate custom uri types to full strings, i.e. translate lp:1424892 to https://bugs.launchpad.net/juju-core/+bug/1424892 or translate gh:juju/juju to $GOPATH/src/github.com/juju/juju

Need plugins to handle vcs, bug trackers, kanban.
