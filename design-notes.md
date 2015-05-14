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

	Name = "github"
	Version = "0.9 alpha"

	[[Command]]
		Name = "dance"
		Short = "makes a dancing banana appear"
		Long = """
		usage:
			dance [options]

			-t=<seconds> 	display for n seconds
		"""
	[[Context]]
	    Name = "bug"
	    Plural = "bugs"
	    Short = "These commands interact with bugs from github."
	    Long = """
	    Bug commands interact with bug from github.  You will need to set up 
	    authentication in the plugin's configuration.
	    """

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
		Name = "github"
		Version = "0.9 alpha"

**Manifest Definition**

This is the specification for what is expected to be in a Q plugin manifest.
These are the Go structs that would be used to generate the manifest.

A note on Versions... versions are expected to match the following regular
expression `[0-9].[0-9]+(.[0-9]+)?( [a-z]+)?`  i.e. major.minor[.patch] [tag]
Where tag is something like "alpha" or "rc1".  When versions are compared, they
are compared by major, then minor, then patch (a missing patch is the same as
patch level 0).  An empty tag is considered higher than a non-empty tag.
Otherwise, tags are compared lexigraphically, so rc1 > beta > alpha.  Note that
this means you have to be careful with tags, since beta2 > beta10.

	type Manifest struct {
		Name string 		// plugin name
		Version string 		// plugin version
		Command []Command 	// commands w/o context
		Context []Context 	// contexts 
		Service []Service 	// services
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

	type Service struct {
		Name string 	// name of the service
		Version string 	// version of the service API
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

