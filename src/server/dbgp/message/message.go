/**
 * Turn DBGp messages into usable data structures.
 */

package message

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

/**
 * Parse dbgp XML message.
 */
func Decode(xmlContent string) (message Message, err error) {

	has_response := (-1 != strings.LastIndex(xmlContent, "</response>"))
	if has_response {
		response, err := decodeResponse(xmlContent)

		if nil == err {
			message = prepareResponseMessage(response)
		}
	}

	has_init := (-1 != strings.LastIndex(xmlContent, "</init>"))
	if has_init {
		init, err := decodeInit(xmlContent)

		if nil == err {
			message = prepareInitMessage(init)
		}
	}

	if !has_response && !has_init {
		err = fmt.Errorf("Unknown message: %s", xmlContent)
	}

	return message, err
}

/**
 * Prepare a message structure based on DBGp engine's initialization attempt.
 */
func prepareInitMessage(init Init) (message Message) {

	message.MessageType = "init"
	message.State = "starting"
	message.Properties.Filename = init.FileURI

	return message
}

/**
 * Prepare a message structure based on DBGp engine's response.
 */
func prepareResponseMessage(response Response) (message Message) {

	message.MessageType = "response"
	message.State = response.Status
	message.Content = response.Content
	message.Properties.Filename = response.Message.Filename
	message.Properties.LineNumber = response.Message.LineNo
	message.Properties.ErrorMessage = response.Error.Message
	message.Properties.ErrorCode = response.Error.Code
	message.Properties.TxId = response.TransactionId
	message.Properties.Command = response.Command
	message.Properties.BreakpointId = response.Id

	if len(response.Breakpoints) > 0 {
		message.Breakpoints = make(map[int]Breakpoint)

		for _, Breakpoint := range response.Breakpoints {
			message.Breakpoints[Breakpoint.Id] = Breakpoint
		}
	}

	message.Context.Local = prepareVariables(response.Variables)

	return message
}

/**
 * Extract variable values.
 */
func prepareVariables(vars []VariableDetails) (variables map[string]Variable) {

	if len(vars) == 0 {
		return
	}

	variables = make(map[string]Variable)

	for _, varDetails := range vars {
		children := prepareVariables(varDetails.Variables)

		hasLoadedChildren := false
		if varDetails.HasChildren {
			// Note: for empty arrays, we always say that we have loaded the children.
			hasLoadedChildren = (varDetails.NumChildren == 0 || len(children) > 0)
		}

		varValue, isBase64 := extractVariableValue(varDetails)

		variables[varDetails.Fullname] = Variable{
			DisplayName:       varDetails.Name,
			VarType:           varDetails.VarType,
			Value:             varValue, // Useful for basic types only.
			AccessModifier:    varDetails.Facet,
			IsCompositeType:   varDetails.HasChildren,
			Children:          children,
			ChildCount:        varDetails.NumChildren,
			HasLoadedChildren: hasLoadedChildren,
			IsBase64:          isBase64,
		}
	}

	return
}

/**
 * Determine variable value and its encoding.
 */
func extractVariableValue(varDetails VariableDetails) (varValue string, isBase64 bool) {

	varValue = varDetails.Value
	isBase64 = (varDetails.Encoding == "base64")

	if !isBase64 {
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(varDetails.Value)

	if err != nil {
		return
	}

	mimeType := http.DetectContentType(decoded)
	isText := len(mimeType) > 4 && mimeType[0:4] == "text"

	if isText {
		varValue = string(decoded)
		isBase64 = false
	}

	return varValue, isBase64
}
