Source: [{{ .Title }}]({{ .URL }})

Files changes: {{ .Files.TotalCount }}

{{ range $idx, $file := .Files.Nodes -}}
* {{ .Path }} (+{{ .Additions }}, -{{ .Deletions }})
{{ end }}
{{- if .Files.PageInfo.HasNextPage -}}TODO! files has next!{{- end }}

{{ if not ( isSamePrev "user" .Author.Login ) -}}
![{{ .Author.Login }}]({{ avatarURL .Author.Login }}) {{ .Author.Login }} - 
{{- end -}}
{{ .CreatedAt | date }}

{{ .Body }}

{{- range $idx, $comment := .Comments.Nodes }}

---

{{ if not ( isSamePrev "user" .Author.Login ) -}}
![{{ .Author.Login }}]({{ avatarURL .Author.Login }}) {{ .Author.Login }} - 
{{- end -}}
{{ .CreatedAt | date }}

{{ .Body }}
{{- end -}}
{{- if .Comments.PageInfo.HasNextPage -}}TODO! comments has next!{{- end }}
