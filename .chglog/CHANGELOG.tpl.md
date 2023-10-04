{{ range .Versions }}
<a name="{{ .Tag.Name }}"></a>
## {{ .Tag.Name }}

> {{ datetime "2006-01-02" .Tag.Date }}

- Full diff - {{ if .Tag.Previous }}**[{{ .Tag.Previous.Name }}...{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }})**{{ else }}{{ .Tag.Name }}{{ end }}  

{{ range .CommitGroups -}}
### {{ .Title }}

{{ range .Commits -}}
* {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts

{{ range .RevertCommits -}}
* {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}