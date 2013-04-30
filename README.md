## shitchat

a simple chat server in Go that serves a browser client.

### who should use shitchat

probably nobody. shitchat has no security, and accepts all HTML, CSS, and javascript code injections.

### why did you write shitchat?

i wanted to play with Go and jquery/ajax, so a chat program seemed like a good place to start.

**how can i run shitchat?**
* `go run server.go`
* now point your browser to localhost:2999

**feature TODO list**
* read the port in from a config file
* add a password screen so that a single password can be required to enter the chat (ie no user/pass combo, just a single secret password)
* read the password in from a config file
* set the initial username from the password screen so people don't have to start as WebUser5000
* bots

**license**

creative commons share alike
