# xman
shell command assistant tool

## Install
go to the code directory and run the following script:

```bash
go install github.com/aronlt/xman@latest
```

## how to use
#### xman mod
replace dependency go module in go.mod file, auto tidy mod file.

#### xman push
automatically complete the local code submission and push the git branch to the remote.

#### xman merge_from
automatically complete local code submission, merge from the pointed branch, </br> detect conflicts, and push to the remote branch.

#### xman merge_to
automatically complete local code submission, merge into the pointed branch, </br> detect conflicts, and push to the remote branch.

#### xman stash
local code temporarily stored (stash).

#### xman recover
the stash code is restored interactively

#### xman tag
auto tag code, support prefix, suffix

## option
* add commit message
```bash
--m="commit message"
```

* skip add review
```bash
--k=true
```

* merge from branch
```bash
--from="branch"
```

* merge to branch
```bash
--to="branch"
```

* tag prefix
```bash
--p="debug"
```

* tag suffix
```bash
--s="demo"
```