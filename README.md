# multidoc

A Go-based utility that processes input through multiple AI models concurrently (OpenAI, Claude, Gemini, and optionally Lambda.Chat) and provides a combination of their responses. Keeping each one's individual best parts. This is useful for research and is sometimes refferred to as multipromting or metaprompting.

## Description

Multidoc allows you to send the same prompt to different AI models simultaneously and combine their responses. The tool:

1. Takes text input from stdin
2. Processes the input through multiple AI models in parallel:
   - GPT-4o (OpenAI)
   - GPT-4.1 (OpenAI)
   - GPT-o3 (OpenAI)
   - GPT-o4-mini (OpenAI)
   - Claude 3.7 Sonnet (Anthropic)
   - Gemini 2.5 Pro Experimental (Google)
   - Gemini 2.5 Flash Preview (Google)
   - Optionally, scrapes Lambda.Chat using Playwright (requires Node.js and dependencies, enable with `-lc` flag)
3. Collects responses from all sources with timing information
4. Generates a combination of all model outputs using GPT-o4-mini
5. Displays the output along with timing details

## Requirements

- Go 1.18 or later
- API keys for:
  - OpenAI API
  - Google Gemini API
  - Anthropic Claude API
- (Optional) Node.js and `npm install` run in the `lambda_scraper` directory if using the `-lc` flag.

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/xssdoctor/multidoc.git
   cd multidoc
   ```

2. Install the required Go dependencies:

   ```
   go mod download
   ```

3. API Key Configuration:

   The application will automatically create a config directory at `~/.config/multidoc` with a default `.env` file on first run.

   Edit this file to add your API keys:

   ```
   ~/.config/multidoc/.env
   ```

   Add your API keys:

   ```
   OPENAI_API_KEY=your_openai_api_key
   GEMINI_API_KEY=your_gemini_api_key
   ANTHROPIC_API_KEY=your_anthropic_api_key
   ```

   Alternatively, you can create a `.env` file in the project directory with the same format.

## Usage

You can pipe text to the program or use interactive input:

```bash
# Using a pipe
echo "Compare the benefits of REST vs GraphQL" | go run app.go

# Using the Lambda.Chat scraper
echo "Explain the concept of metaprompting" | go run app.go -lc

# Using interactive input
go run app.go
# Then type your prompt and press Ctrl+D when finished
```

For larger prompts, you can use a text file:

```bash
cat prompt.txt | go run app.go
```

The output will display:

- Progress updates as each model responds
- A combination of the models best outputs
- Timing information for each model and the total execution

## Build

To build an executable:

```bash
go build -o multidoc app.go
```

Then you can run it directly:

```bash
./multidoc
```

When you run the application for the first time, it will check if the configuration directory exists at `~/.config/multidoc` and create it if necessary, along with a template `.env` file. The application will then exit, allowing you to add your API keys before running it again.

## License

MIT

## Author

xssdoctor & jhaddix
