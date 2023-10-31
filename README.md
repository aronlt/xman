# xman
shell command assistant tool

## How To Use
go to the code directory and run the following script:
```bash
go build . &&
chmod +x xman &&
mv xman /usr/local/bin
```
## commands
#### xman mod
implementing module replacement for the Go language.
![tidy](static/tidy.png)

#### xman push
automatically complete the local code submission and push the git branch to the remote.
![push](static/push.png)

#### xman merge
automatically complete local code submission, merge into the pointed branch, </br> detect conflicts, and push to the remote branch.
![merge](static/merge.png)

#### xman stash
local code temporarily stored (stash).
![stash](static/stash.png)

#### xman recover
the stash code is restored interactively
![recover](static/recover.png)
