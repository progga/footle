/**
 * @file
 * Class definition for the Help class.
 */

package help

import (
	"fmt"
	"strings"
)

type helptext struct {
	cmdNAliases []string
	explanation string
}

type helptextPointers map[string]*helptext

type Help struct {
	cmdList            [][]string
	cmdHelptextMapping helptextPointers
}

/**
 * Prepare overall or command specific helptext.
 */
func (h *Help) Me(args []string) string {

	isSeekingHelpForCmd := (len(args) > 0)

	if isSeekingHelpForCmd {
		cmdOrAlias := args[0]
		return h.forOne(cmdOrAlias)
	} else {
		return h.forAll()
	}
}

/**
 * List of all commands and their aliases.
 */
func (h *Help) forAll() (helpOverview string) {

	var allCmdNAliases string

	for _, cmdNAliases := range h.cmdList {
		allCmdNAliases += fmt.Sprintf("%s\n", strings.Join(cmdNAliases, ", "))
	}

	helpOverview = fmt.Sprintf("help [cmd]\nExample: help bye\n\nAvailable commands and their aliases:\n%s", allCmdNAliases)
	return helpOverview
}

/**
 * Explanation of just one command and its aliases.
 */
func (h *Help) forOne(cmd string) (helptextForCmd string) {

	helptextDetailPointer, ok := h.cmdHelptextMapping[cmd]
	if !ok {
		// @see www.yodaspeak.co.uk
		return "Provide a valid command, you must.\n"
	}

	helptextDetail := *helptextDetailPointer
	helptextForCmd = fmt.Sprintf("%s: %s\n", strings.Join(helptextDetail.cmdNAliases, ", "), helptextDetail.explanation)

	return helptextForCmd
}

/**
 * Add the given helptexts into the existing lists.
 */
func (h *Help) add(texts []helptext) {

	for index, helptextDetail := range texts {
		h.cmdList = append(h.cmdList, helptextDetail.cmdNAliases)

		for _, cmdOrAlias := range helptextDetail.cmdNAliases {
			h.cmdHelptextMapping[cmdOrAlias] = &texts[index]
		}
	}
}

/**
 * Load all commands and their explanations.
 */
func (h *Help) prepare(cliCmdList, footleCmdList, DBGpCmdList []helptext) {

	h.cmdHelptextMapping = make(map[string]*helptext)
	h.add(cliCmdList)
	h.add(footleCmdList)
	h.add(DBGpCmdList)
}
