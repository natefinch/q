# Q

Q is a command line tool for programmers to make their lives better in any way
they see fit.  At its base is a todo list and note taking application, but its
support for plugins makes it infinitely extendible in any way imaginable.

## Why is it called Q?

<img align=right width=100 src="https://cloud.githubusercontent.com/assets/3185864/4443856/dd44f740-47eb-11e4-8429-0f2b9fb79d1c.png">Q is a character on Star the Next Generation with apparent omnipotence.  This fits well with the design goals of the Q application, because it is intended to be extensible to do anything.  He also
happens to have a single-letter name, which makes the command very easy to type.  

## Built-ins

### Tasks 

Tasks are items in a todo list.  They have a title, an optional description, a
body, an ID, and a weight.  The Task with the highest weight is the one that
should be done next.  Tasks without a weight are ordered by date of creation,
with the oldest being the one you should do next.  Tasks generally have three
states - todo, doing, done.

You may add tags and custom metadata to tasks and task lists to make filtering
and searching easier, as well as to support plugin use cases.

### Notes

Notes are just that - notes... they can have tags and metadata just like tasks,
but don't have the concept of weight or states like todo, doing, or done.

## Plugins

Q is extendible via executable plugins.  Plugins must be called q-[name] and are
invoked like `q [name] [options]`.  

## Existing plugins

### gh

gh is a plugin for integrating with github to translate issues and pull requests
into tasks in your list.

### kb

kb is a plugin for turning your task lists into one or more Kanban boards which
may optionally be shared with others.



