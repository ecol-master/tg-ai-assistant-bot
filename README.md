# Telegram Personal Assistant Bot

![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)

This is a Telegram bot written in Golang that serves as a personal assistant. It leverages the [Telegram Bot API](https://core.telegram.org/bots/api) for communication with users and the [OpenAI GPT-4 Turbo Preview](https://beta.openai.com/docs/introduction/gpt-4) for generating responses.

## Features

- Responds to user messages and commands.
- Implements basic command: `/help`
- Can generate intelligent responses using GPT-4 Turbo Preview for any user input.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/telegram-personal-assistant-bot.git
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Update the bot token and OpenAI API key in the `main.go` file.

4. Build and run the bot:

   ```bash
   go build
   ./your-bot-executable
   ```

## Getting Started


### Prerequisites

- Go installed on your system.
- Docker installed on your system.
- You have a token from your personal **Telegram-Bot** and personal key to **OpenAI API**.

### How to run Project

**1. Set up `.env` file**

Paste your secret keys into such variables in config file:
   ```
   TELEGRAM_TOKEN=<your key>
   OPENAI_TOKEN=<your key>
   ...
   EMAIL_USERNAME="your email" # your email address
   EMAIL_PASSWORD="password for email app"
   ```
   **Note**: `EMAIL_PASSWORD` variable is not a regular password from email. These password i need to create in your email settings: **Security** -> **App passwords**.  Next, create an application and copy the resulting code to a variable `EMAIL_PASSWORD`

**2. Set up Google OAuth to work with Calendar API**

The guide about creating `token.json` and `credentials.json` files i write in this [file (How to start with Google OAuth)](./docs/google_oauth.md).

**3. Starting postgresql database**:
   ``` bash
   docker run --name some-postgres -e POSTGRES_DB=planpilot -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:latest
   ```

**4. Starting Golang Application:**
```bash
cd cmd
go run main.go
```

### Usage

1. Start a chat with your bot on Telegram.

2. Send messages to the bot, and it will respond to your messages using OpenAI's GPT-4 for intelligent responses.

3. You can also use the following commands:
   - `/start`: Starting the Bot
   - `/help`: Get information about available commands.

### Customization

You can customize the bot's responses, add more commands, and adapt its behavior as needed by modifying the `main.go` file. For more advanced interactions with GPT-3.5 Turbo, refer to the [OpenAI API documentation](https://beta.openai.com/docs/introduction/gpt-3).

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Acknowledgments

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [OpenAI GPT-4 Turbo](https://beta.openai.com/docs/introduction/gpt-3)
