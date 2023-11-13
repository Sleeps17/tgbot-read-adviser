package telegram

const msgHelp = `I can save and keep you pages. Also I can offer you them to read.

In order to save the page, just send me al link to it.

In order to get a random page from your list, send me command /rnd.

In order to remove a page from your list, send me command /remove and link on this page.

In order to get all pages from your list? send me command /all.
`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedPages   = "You have no saved pages 🙊"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
	msgDeleted        = "Deleted 💀"
	msgNoThisPage     = "You don't have this page in your list 🎯"
	msgEmptyList      = "Your pages list is empty 🧩"
)
