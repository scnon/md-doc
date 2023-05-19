package model

const DocTmpl = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>{{.Title}}</title>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.2.0/github-markdown.min.css" integrity="sha512-Ya9H+OPj8NgcQk34nCrbehaA0atbzGdZCI2uCbqVRELgnlrh8vQ2INMnkadVMSniC54HChLIh5htabVuKJww8g==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<style>
img,
svg {
	max-width: 100%;
	display: block;
}

html {
	color-scheme: light dark;
}

body {
	font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
	font-size: 1.125rem;
	line-height: 1.5;
}

.content {
	width: min(128ch, 100% - 4rem);
	margin-inline: auto;
	padding-left: 2rem;
	padding-right: 2rem;
}

.markdown-body {
	padding: 12px;
}
</style>
</head>
<body>
<div class="content">
<div class="markdown-body">
{{.Content}}
</div>
</div>
</body>
</html>
`
