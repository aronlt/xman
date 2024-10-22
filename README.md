# xman
shell command assistant tool

## Install
go to the code directory and run the following script:

```bash
go install github.com/aronlt/xman@latest
```

## how to use
#### xman lint
lint all files that added to stash zone, you should install golang-lint first
```bash
xman lint
```
###### option

* lint all files in dir

    --dir=true

#### xman push
automatically complete the local code submission and push the git branch to the remote.

###### option
* add review
 
    --r=true
 
* add commit message
 
  --m="commit message"

#### xman merge_from
automatically complete local code submission, merge from the pointed branch, </br> detect conflicts, and push to the remote branch.

###### option
* add commit message
 
    --m="commit message"
 
* merge from branch
 
  --from="branch"

#### xman merge_to
automatically complete local code submission, merge into the pointed branch, </br> detect conflicts, and push to the remote branch.

###### option
* add commit message

  --m="commit message"

* merge to branch

  --to="branch"

#### xman stash
local code temporarily stored (stash).

#### xman recover
the stash code is restored interactively

#### xman tag
auto tag code, support prefix, suffix

###### option
* tag prefix
 
  --p="debug"

* tag suffix
 
  --s="demo"


#### xman local_branch
list local branch info

#### xman remote_branch
list remote branch info

#### xman list_tags
list tags info

#### xman last_commit
list all branch last commit info

#### xman checkout
checkout to another branch, submit local branch
###### option
* checkout to new branch

    --b="branch_name" 

* checkout new branch from local branch
 
    --cf="from_branch" 







