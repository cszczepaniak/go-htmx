build: css templ
	@go build .

css:
	@echo "Rebuilding tailwind styles"
	@tailwindcss -i web/app.css -o web/dist/styles.css

templ:
	@echo "Regenerating templates"
	@templ generate
