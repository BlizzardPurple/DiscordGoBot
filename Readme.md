# DiscordBot Code Documentation

## About the Developer
Developer Name: Garvit Jain   
Email: garvitrita2002@gmail.com

## Project Title
DiscordBot with Bad Language and Misgender Detection

## Project Description
This project provides a Discord bot designed to monitor and moderate chat environments within a Discord server. The bot has the capability to detect bad language usage and misgendering, which can help enforce a positive and inclusive community culture.

## Features
- **Bad Language Detection**: The bot scans incoming messages for a predefined list of inappropriate phrases and responds accordingly.
- **Misgender Detection**: Detects the usage of gender-specific terms and encourages the use of gender-neutral alternatives.
- **Role Management**: Manages roles based on user reactions, ensuring that a user can only have one role at a time.
- **Database Integration**: Stores user responses and records them in a PostgreSQL database.
- **User Interaction**: Provides direct messaging capabilities to collect user preferences and store them in the database.
- **Random Country Generation**: Responds with a random country name when prompted.
- **Registration Prompt**: Sends a registration prompt to users, allowing them to input personal details.
- **Welcome Greetings**: Greets new members to the server.

## Installation and Setup
To run this Discord bot, you will need:
- Go programming language installed on your system.
- Access to a Discord server to invite the bot.
- A PostgreSQL database for storing user data.

### Environment Variables
Before running the bot, ensure that you have the following environment variables set:
- `BOT_TOKEN`: Your Discord bot token.
- `DSN`: Connection string to your PostgreSQL database.
- `ADMIN_ID`: The Discord user ID of the server administrator.
- `FIRE_NATION`: The role ID for the Fire Nation role in your Discord server.
- `WATER_NATION`: The role ID for the Water Nation role in your Discord server.

### Running the Bot
1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Run the command `go run .` to start the bot.

## Usage
Once the bot is running, it will listen for commands and reactions in the configured Discord server. Here are some commands and reactions you can use:
- **!gobot hello**: The bot will respond with a greeting directed at the user.
- **!gobot country**: The bot will reply with a random country name.
- **!gobot register**: The bot will prompt the user to register and input their preferences.
- **!gobot answers [ID]**: Retrieve and display user preferences from the database using the provided ID.
- Emoji reactions: Users can react with specific emojis to gain roles.

## Contribution
Contributions to this project are welcome. Please fork the repository, make your changes, and submit a pull request.

---

Please note that the actual implementation of the bot may vary based on the specific requirements of the server and the developer's preferences. Always test your bot thoroughly to ensure it functions as expected and adheres to the community guidelines.