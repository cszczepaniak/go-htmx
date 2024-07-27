build:
	@templ generate
	@go build .

css:
	@tailwindcss -i web/app.css -o web/dist/styles.css
