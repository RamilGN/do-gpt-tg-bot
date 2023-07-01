.PHONY: deploy
deploy:
	@doctl serverless deploy .

.PHONY: lint
lint:
	@golangci-lint run ./packages/core/bot --out-format tab

.PHONY: prepare-to-deploy
prepare-to-deploy: tg-set-commands tg-set-webhook

.PHONY: url
url:
	@doctl serverless function get --url core/bot

.PHONY: invoke
invoke:
	@doctl serverless functions invoke core/bot

.PHONY: watch
watch:
	@doctl serverless watch .

.PHONY: logs
logs:
	@doctl serverless activations logs --follow --function core/bot

.PHONY: tg-set-webhook
tg-set-webhook:
	@echo "Setting webhook..."
	@curl -sS \
		-X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setWebhook" \
		-H "Content-Type: application/json" \
		-d '{"url":"'$$(doctl serverless function get --url core/bot)'"}' | jq

.PHONY: tg-get-webhook
tg-get-webhook:
	@curl -sS "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getWebhookInfo" | jq

.PHONY: tg-set-commands
tg-set-commands:
	@echo "Setting commands..."
	@curl -sS \
		-X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setMyCommands" \
		-H "Content-Type: application/json" \
		-d '{"commands":[{"command":"ask", "description":"ask bot"}]}' | jq

.PHONY: tg-get-commands
tg-get-commands:
	@curl -sS "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getMyCommands" | jq
