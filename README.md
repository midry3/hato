# hato
This is a CLI CheckList tool.

# How to install
If you have `go`:
```bash
$ go install github.com/midry3/hato
```

or download binary from [here](https://github.com/midry3/hato/releases/latest).

# Usage
```bash
# Initialize
$ hato

# Add a checklist
$ hato add Check A
$ hato add Check B

# Check the list
$ hato
# or
$ hato check
```

The way of checking the list is this:
- If ok, `Enter` and next.
- If not ok, `Esc` and stop.

If all checklists are ok, and if you set actions, run the actions.

# Format
```yaml:hato.yml
# https://github.com/midry3/hato

checklist:
  - Check A
  - Check B

actions:
  - echo Action1
  - echo Action2
```
And result will be like this:
```bash:result
$ hato
[1]: Check A => ✅
[2]: Check B => ✅
All checklists are ok!

Running 1/2: `echo Action1` ...
Action1

Running 2/2: `echo Action2` ...
Action2

✅All actions have been completed!
```

# Example
Please edit `hato.yml`.
```yaml:hato.yml
# https://github.com/midry3/hato

checklist:
  - Current branch is main?
  - Merged from dev?
  - Updated version?

actions:
  - git pull origin main
  - git push origin main
```