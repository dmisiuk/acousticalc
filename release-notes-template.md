## Debug Template Variables

**Tag**: {{ .Tag }}
**TagSubject**: {{ .TagSubject }}
**TagContents**: {{ .TagContents }}
**TagBody**: {{ .TagBody }}
**Version**: {{ .Version }}
**Branch**: {{ .Branch }}
**Commit**: {{ .ShortCommit }}

---

## Actual Release Notes

{{ .TagContents }}