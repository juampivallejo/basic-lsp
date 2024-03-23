package analysis

import (
	"basiclsp/lsp"
	"fmt"
	"strings"
)

type State struct {
	// Map file names to contents
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnostics(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnostics(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	// In real life this would look up the type analysis code

	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File %s, Characters: %d", uri, len(document)),
		},
	}

}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// In real life this would look up the type definition

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      1,
					Character: 2,
				},
			},
		},
	}

}
func (s *State) TextDocumentCodeAction(id int, uri string) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri]
	actions := []lsp.CodeAction{}

	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "NeoVim",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS **** with a superior editor",
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})
			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS ****",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Censor Shitty Editors",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})

		}

	}
	response := lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
	return response
}

func LineRange(row, idx, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      row,
			Character: idx,
		},
		End: lsp.Position{
			Line:      row,
			Character: end,
		},
	}
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	// Here should ask for good completions based on what the user is trying to type
	completions := []lsp.CompletionItem{
		{
			Label:         "NeoVim (BTW)",
			Detail:        "Best Editor",
			Documentation: "Not for noobs?",
		},
	}

	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: completions,
	}
	return response

}

func getDiagnostics(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			d := lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Severity: 1,
				Source:   "Common Sense",
				Message:  "Replace for superior editor",
			}
			diagnostics = append(diagnostics, d)
		}
	}

	return diagnostics
}
