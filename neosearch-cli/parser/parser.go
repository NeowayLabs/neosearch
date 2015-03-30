package parser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/NeowayLabs/neosearch/engine"
	"github.com/NeowayLabs/neosearch/index"
	"github.com/NeowayLabs/neosearch/utils"

	"github.com/iNamik/go_lexer"
)

type parserState struct {
	IsUsing                     bool
	IsCommand                   bool
	IsValue                     bool
	IsDoubleQuotedString        bool
	IsSingleQuotedString        bool
	IsEscapedDoubleQuotedString bool
	IsEscapedSingleQuotedString bool
}

// We define our lexer tokens starting from the pre-defined EOF token
const (
	TokenEOF = lexer.TokenTypeEOF
	TokenNil = lexer.TokenTypeEOF + iota
	TokenSpace
	TokenNewline
	TokenDoubleQuotedString
	TokenSingleQuotedString
	TokenEscapedDoubleQuotedString
	TokenEscapedSingleQuotedString
	TokenSemiColon
	TokenWord
	TokenNumbers
	TokenUsing
	TokenSet
	TokenGet
	TokenMergeSet
	TokenDelete
)

var bytesNonWord = []byte{' ', '\t', '\f', '\v', '\n', '\r', ';', '"', '\'', '\\', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

var bytesIntegers = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

var bytesSpace = []byte{' ', '\t', '\f', '\v'}

var bytesDoubleQuotedStrings = []byte{'"'}

var bytesSingleQuotedStrings = []byte{'\''}

var bytesEscapedDoubleQuotedString = []byte{'\\', '"'}

var bytesEscapedSingleQuotedString = []byte{'\\', '\''}

const charNewLine = '\n'

const charReturn = '\r'

const charSemicolon = ';'

var commandsAvailable = []string{
	"set",
	"mergeset",
	"get",
	"delete",
	"batch",
	"flushbatch",
}

// Checks if the given command is valid.
func isValidCommand(cmd string) bool {
	for i := 0; i < len(commandsAvailable); i++ {
		if cmd == commandsAvailable[i] {
			return true
		}
	}

	return false
}

func setQuotedString(token string, command *engine.Command, pState *parserState) {
	if pState.IsUsing {
		command.Index += token
	} else if pState.IsCommand {
		command.Key = []byte(string(command.Key) + token)
	} else if pState.IsValue {
		command.Value = []byte(string(command.Value) + token)
	}
}

func validateBatch(cmd engine.Command) bool {
	if cmd.Command == "batch" && cmd.Index != "" &&
		cmd.Key == nil && cmd.Value == nil {
		return true
	}

	return false
}

func validateFlushBatch(cmd engine.Command) bool {
	if cmd.Command == "flushbatch" && cmd.Index != "" &&
		cmd.Key == nil && cmd.Value == nil {
		return true
	}

	return false
}

func validateSetters(cmd engine.Command) bool {
	if cmd.Command == "set" || cmd.Command == "mergeset" {
		if cmd.Index != "" && cmd.Key != nil &&
			cmd.Value != nil {
			return true
		}
	}

	return false
}

func validateGetters(cmd engine.Command) bool {
	if cmd.Command == "get" || cmd.Command == "delete" {
		if cmd.Index != "" && cmd.Key != nil &&
			cmd.Value == nil {
			return true
		}
	}

	return false
}

func validateCommand(cmd engine.Command) bool {
	if cmd.Command == "set" {
		return validateSetters(cmd)
	} else if cmd.Command == "get" {
		return validateGetters(cmd)
	} else if cmd.Command == "mergeset" {
		return validateSetters(cmd)
	} else if cmd.Command == "delete" {
		return validateGetters(cmd)
	} else if cmd.Command == "batch" {
		return validateBatch(cmd)
	} else if cmd.Command == "flushbatch" {
		return validateFlushBatch(cmd)
	}

	return false
}

// FromString parses the cmdline string and returns an array of engine.Command
func FromString(cmdline string, listCommands *[]engine.Command) error {
	return FromReader(strings.NewReader(cmdline), listCommands)
}

// FromReader parse the file
func FromReader(file io.Reader, listCommands *[]engine.Command) error {
	var command engine.Command

	pState := parserState{}

	// Create our lexer
	// NewSize(startState, reader, readerBufLen, channelCap)
	lex := lexer.NewSize(lexFunc, file, 100, 1)
	var lastTokenType = TokenNil

	// Process lexer-emitted tokens
	for t := lex.NextToken(); lexer.TokenTypeEOF != t.Type(); t = lex.NextToken() {
		switch t.Type() {
		case TokenWord:
			if lastTokenType != TokenWord {
				tokenValue := string(t.Bytes())

				// TokenWord is a double or single quoted string?
				if pState.IsDoubleQuotedString || pState.IsSingleQuotedString {
					if !pState.IsUsing && !pState.IsCommand && !pState.IsValue {
						return errors.New("Invalid quoted string: " + tokenValue)
					}

					setQuotedString(tokenValue, &command, &pState)

					// TokenWord is the Index name?
					// using <TokenWord> ...
				} else if pState.IsUsing {
					command.Index = tokenValue
					pState.IsUsing = false

					// TokenWord is the key of command?
					// using document.db mergeset <TokenWord> ...
				} else if pState.IsCommand {
					command.Key = []byte(tokenValue)
					pState.IsCommand = false
					pState.IsValue = true

					// TokenWord is the command value?
					// using document.db mergeset name <TokenWord>
				} else if pState.IsValue {
					command.Value = []byte(tokenValue)
					pState.IsValue = false
				} else {
					// Here we handle the available KEYWORDS

					// Keyword USING
					if !pState.IsUsing && strings.ToLower(tokenValue) == "using" {
						pState.IsUsing = true
					} else if !pState.IsCommand {
						// Must be a command KEYWORD
						// see commandsAvailable

						if !isValidCommand(tokenValue) {
							return errors.New("Invalid keyword '" + tokenValue + "'")
						}

						pState.IsCommand = true
						pState.IsUsing = false
						command.Command = strings.ToLower(tokenValue)
					}
				}
			}
		case TokenDoubleQuotedString:
			if pState.IsDoubleQuotedString {
				if pState.IsCommand {
					pState.IsCommand = false
					pState.IsValue = true
				} else if pState.IsUsing {
					pState.IsUsing = false
				} else if pState.IsValue {
					pState.IsValue = false
				}
			}

			if pState.IsSingleQuotedString {
				setQuotedString(string(t.Bytes()), &command, &pState)
			} else {
				pState.IsDoubleQuotedString = !pState.IsDoubleQuotedString
			}
		case TokenSingleQuotedString:
			if pState.IsSingleQuotedString {
				if pState.IsCommand {
					pState.IsCommand = false
					pState.IsValue = true
				} else if pState.IsUsing {
					pState.IsUsing = false
				} else if pState.IsValue {
					pState.IsValue = false
				}
			}

			if pState.IsDoubleQuotedString {
				setQuotedString(string(t.Bytes()), &command, &pState)
			} else {
				pState.IsSingleQuotedString = !pState.IsSingleQuotedString
			}

		case TokenEscapedDoubleQuotedString:
			if pState.IsSingleQuotedString {
				panic("Escaped double quoted string inside single quoted string...")
			} else if pState.IsDoubleQuotedString {
				setQuotedString(string(t.Bytes()[1:]), &command, &pState)
			}

			pState.IsEscapedDoubleQuotedString = !pState.IsEscapedDoubleQuotedString

		case TokenEscapedSingleQuotedString:
			if pState.IsDoubleQuotedString {
				return errors.New("Escaped single quoted string inside double quoted string...")
			} else if pState.IsSingleQuotedString {
				setQuotedString(string(t.Bytes()), &command, &pState)
			}
		case TokenSpace:
			// Spaces only makes difference inside quotes
			if pState.IsDoubleQuotedString || pState.IsSingleQuotedString {
				setQuotedString(string(t.Bytes()), &command, &pState)
			}

		case TokenNewline:
			// New lines only makes difference inside quotes
			if pState.IsDoubleQuotedString || pState.IsSingleQuotedString {
				setQuotedString(string(t.Bytes()), &command, &pState)
			}
		case TokenSemiColon:
			if pState.IsSingleQuotedString || pState.IsDoubleQuotedString {
				setQuotedString(string(t.Bytes()), &command, &pState)
			} else {
				*listCommands = append(*listCommands, command)
				command = engine.Command{}
				pState = parserState{}
			}
		case TokenNumbers:
			tokenValue := string(t.Bytes())
			tokenIntValue, err := strconv.Atoi(tokenValue)

			if err != nil {
				return fmt.Errorf("Failed to convert %s to integer", tokenValue)
			}

			if pState.IsSingleQuotedString || pState.IsDoubleQuotedString {
				setQuotedString(tokenValue, &command, &pState)
			} else if pState.IsUsing {
				if index.ValidateIndexName(tokenValue) {
					command.Index = tokenValue
					pState.IsUsing = false

					// TokenNumbers is the key of command?
					// using document.db mergeset <TokenNumbers> ...
				}
			} else if pState.IsCommand {
				command.Key = utils.Int64ToBytes(int64(tokenIntValue))
				pState.IsCommand = false
				pState.IsValue = true

				// TokenNumbers is the command value?
				// using document.db mergeset name <TokenNumbers>
			} else if pState.IsValue {
				tokenIntValue, err := strconv.Atoi(tokenValue)

				if err != nil {
					return err
				}

				if command.Command == "mergeset" {
					command.Value = utils.Uint64ToBytes(uint64(tokenIntValue))
				} else {
					command.Value = utils.Int64ToBytes(int64(tokenIntValue))
				}

				pState.IsValue = false
			}
		default:
			return errors.New("Failed to parse line at '" + string(t.Bytes()) + "'")
		}

		lastTokenType = t.Type()
	}

	// Checks if the last command was correctly parsed but
	// doesn't have the semicolon at the end...
	if validateCommand(command) {
		*listCommands = append(*listCommands, command)
		command = engine.Command{}
		pState = parserState{}

		// Checks if exists a invalid partial command
	} else if command.Index != "" || command.Command != "" ||
		command.Key != nil || command.Value != nil {
		return fmt.Errorf("The last command wasn't correctly finished nor have the semicolon at end: %v", command)
	}

	return nil
}

func lexFunc(l lexer.Lexer) lexer.StateFn {
	// EOF
	if l.MatchEOF() {
		l.EmitEOF()
		return nil // We're done here
	}

	if l.MatchMinMaxBytes(bytesEscapedDoubleQuotedString, 2, 2) {
		l.EmitTokenWithBytes(TokenEscapedDoubleQuotedString)

	} else if l.MatchMinMaxBytes(bytesEscapedSingleQuotedString, 2, 2) {
		l.EmitTokenWithBytes(TokenEscapedSingleQuotedString)

	} else if l.MatchOneOrMoreBytes(bytesDoubleQuotedStrings) {
		l.EmitTokenWithBytes(TokenDoubleQuotedString)

	} else if l.MatchOneOrMoreBytes(bytesSingleQuotedStrings) {
		l.EmitTokenWithBytes(TokenSingleQuotedString)

	} else if l.NonMatchOneOrMoreBytes(bytesNonWord) {
		l.EmitTokenWithBytes(TokenWord)

	} else if l.MatchOneOrMoreBytes(bytesIntegers) {
		l.EmitTokenWithBytes(TokenNumbers)

		// Space run
	} else if l.MatchOneOrMoreBytes(bytesSpace) {
		l.EmitTokenWithBytes(TokenSpace)

		// Line Feed

	} else if charNewLine == l.PeekRune(0) {
		l.NextRune()
		l.EmitTokenWithBytes(TokenNewline)
		l.NewLine()

		// Carriage-Return with optional line-feed immediately following
	} else if charReturn == l.PeekRune(0) {
		l.NextRune()
		if charNewLine == l.PeekRune(0) {
			l.NextRune()
		}
		l.EmitTokenWithBytes(TokenNewline)
		l.NewLine()
	} else if charSemicolon == l.PeekRune(0) {
		l.NextRune()
		l.EmitTokenWithBytes(TokenSemiColon)
	} else {
		panic("Failed to parse line at '" + string(l.PeekRune(0)))
	}

	return lexFunc
}
