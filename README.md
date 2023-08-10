# Crypto News Summarizer 💼 📈

Leveraging the power of OpenAI GPT to process cryptocurrency news and delivering concise summaries directly to your Telegram and Twitter!

_Note: This project is built using Go Workspace, so dependencies are defined for each service independently_

**Language:** GoLang  
**Architecture:** Microservice

## Features 🌟

- **News Summarization**: Condense lengthy crypto news into digestible summaries.
- **Multi-Platform Notification**: Directly posts summaries to Telegram and Twitter.
- **Customizable Configurations**: Control bot behavior and settings directly via Telegram.
- **Localization**: Supports multiple languages, ensuring everyone gets the news in the language they prefer.

## Microservices & Packages 📦

### Services

1. **Configuration-bot**: Control system settings & manage user access.
2. **GPT Service**: Defines the GPT bot wrapper.
```go
type Bot interface {
    Ask(ctx context.Context, prompt, context string, language string) (*Message, error)
}
```
3. **Migrator**: Manage database migrations.
4. **Parser**: Query news sources and store news snippets in the database.
5. **Telegram-bot**: Fetch summarized news from the database and dispatch to Telegram channels.
6. **Twitter-bot**: Fetch summarized news and tweet them out.

### Packages

- **Common**: Logic that's universal across services (DB access, error handling, utility functions).
- **Docker**: Compose files tailored for various environments (local, dev, prod).
- **Localization**: Language mapping files (e.g., `pl.locale.yaml` for Polish).
- **Scripts**: Shell scripts supporting AWS CI/CD with GitHub Actions.
- **Templates**: Define the structure of posts and commands (e.g., `news.post.tmpl`).

### Configuration
- **config.yaml**: is a sample configuration file. Fill the `...` with your keys values.
- **docker/config.docker.yaml**: is a sample docker configuration . It follows the same format that is in the config.yaml. But unlike the first one, this file is used to launch in docker environment.

## Setup 🚀

*Details about how to set up and run the project (make sure your go version is 1.20+)*

1. `git clone https://github.com/StepanTita/crypto-news.git`
2. `go work sync`
3. `docker-compose -f ./docker/docker-compose.local.yaml up -d`

*To stop the project:*
`docker-compose -f ./docker/docker-compose.local.yaml down -v`

## License 📄

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.
