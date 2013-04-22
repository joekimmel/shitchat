## shitchat

a simple chat server in Go that servers a browser client.

### who should use shitchat

probably nobody. shitchat has no security, and accepts all HTML, CSS, and javascript code injections.

### why did you write shitchat?

i wanted to play with Go and jquery/ajax, so a chat program seemed like a good place to start.

**how can i run shitchat?**
* `go run server.go`
* now point your browser to localhost:2999

**feature TODO list**
* add a backend feature to flush the messages (POSTS)
* add a frontend button to trigger the backend flush
* add a frontend button to clear the frontend's messages
* add a password screen so that a single password can be required to enter the chat (ie no user/pass combo, just a single secret password)


**license**

creative commons share alike
