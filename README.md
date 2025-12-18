# hato
This is a CLI CheckList tool.

# How to install
If you have `go`:
```bash
$ go install github.com/midry3/hato@latest
```
s
Linux, Mac:
```bash
$ curl https://raw.githubusercontent.com/midry3/hato/main/install.sh | sh
```

or download binary from [here](https://github.com/midry3/hato/releases/latest).

# Usage
```bash
# Initialize
$ hato

# Add an item of checklist
$ hato --add Check A
$ hato --add Check B

# Create `test` checklist and add an item of checklist
$ hato test --add Test Check A
$ hato test --add Test Check B

# Check the list
$ hato
$ hato test
```

Please edit `hato.yml` on current directory.

The way of checking the list is this:
- If ok, `Enter` and next.
- If not ok, `Esc` and stop.

If all checklists are ok, and if you set actions, run the actions.

# Format
```yaml:hato.yml
# https://github.com/midry3/hato

default:
  nargs: 0
  aliases:
    - Alias1
    - Alias2
  inform:
  checklist:
    - Check A
    - Check B
  actions:
    - echo Action1
    - echo Action2

checklist_name:
  nargs: 2
  aliases:
    - Alias3
  inform:
    - echo Information1
    - echo Information2
  checklist:
    # You can recieve an argument by %position
    - Value1 is "%1", ok?
    - Value2 is "%2", ok?
  actions:
    - echo %1
    - echo %(2)
```

And result will be like this:
```bash
$ hato
[1]: Check A => ✅
[2]: Check B => ✅
All of checklist are ok!

Running 1/2: `echo Action1` ...
Action1
Running 2/2: `echo Action2` ...
Action2

✅All actions have been completed!
```

```bash
$ hato checklist_name 123 abc
――――――――――――――――――――――――――――Information――――――――――――――――――――――――――――――――――――――――――
Information1
Information2
―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
[1]: Value1 is "123", ok? => ✅
[2]: Value2 is "abc", ok? => ✅
All of checklist are ok!

Running 1/2: `echo 123` ...
%1
Running 2/2: `echo abc` ...
%(2)

✅All actions have been completed!
```

```bash
$ hato Alias3
This checklist needs just 2 arguments.
```

# Example
```yaml:hato.yml
# https://github.com/midry3/hato

default:
  checklist:
    - Are you happy?

push:
  nargs: 1
  aliases:
    - p
  checklist:
    - Current branch is %1?
    - Updated version?
    - README is ok?
  actions:
    - git pull origin %1
    - git push origin %1

commit:
  nargs: 1
  aliases:
    - c
  inform:
    - git status
  checklist:
    - Checked the stages?
    - Removed test debug codes?
    - Are you ok this commit message?
  actions:
    - git commit -m %1
```
