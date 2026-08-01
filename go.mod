module github.com/septagon-oss/pk-ui

go 1.26

require (
	github.com/septagon-oss/pk-design v0.4.0
	maragu.dev/gomponents v1.3.0
)

require (
	github.com/septagon-oss/styleengine v0.1.1
	github.com/septagon-oss/tw v0.2.2
)

require (
	github.com/tdewolff/minify/v2 v2.24.14 // indirect
	github.com/tdewolff/parse/v2 v2.8.14 // indirect
)

retract (
	v0.2.3
	v0.0.0 // broken: contained local replace directives
)
