# next-day-todo

A little helper script to run a todo list based on markdown files.

The todo list files are expected to be in a directory. Each file has a name like `YYYY-MM-DD-Day.md` (that's `2006-01-02-Mon.md` in golang time format).

Each file has the following layout:

```markdown
# 2006-01-02 (Monday)

## TODO

* list
* of
* tasks


## Done

* things
* completed
* today
```

The `next-day-todo` script finds the latest file (by date found in the name) in the directory, copies the TODO list, sets the Done list to `* nothing`, and writes this out to a new file for today's date. It also copies the file's full path to the OS X clipboard.

I usually start my day with

```bash
$ next-day-todo
$ vim ⌘-V
```

Originally inspired by https://dev.to/jlhcoder/the-power-of-the-todo-list
