# exec2json

```bash
exec2json /bin/bash -c 'echo this is my stdout; echo that is my stderr 1>&2' | jq
```
Outputs:

```json 
{
  "command": [
    "/bin/bash",
    "-c",
    "echo this is my stdout; echo that is my stderr 1>&2"
  ],
  "status": 0,
  "stderr": "that is my stderr\n",
  "stdout": "this is my stdout\n",
  "took": 0.004179875
}
```
