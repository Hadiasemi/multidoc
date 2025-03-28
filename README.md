# haddixscript

A Go-based utility that processes input through multiple AI models concurrently (OpenAI, Claude, and Gemini) and provides a summarized comparison of their responses.

## Description

HaddixScript allows you to send the same prompt to different AI models simultaneously and compare their responses. The tool:

1. Takes text input from stdin
2. Processes the input through multiple AI models in parallel:
   - GPT-4o (OpenAI)
   - o3-mini (OpenAI)
   - o1-mini (OpenAI)
   - Gemini 2.0 Flash (Google)
   - Claude 3.7 Sonnet (Anthropic)
3. Collects responses from all models with timing information
4. Generates a summary of all model outputs using o1-mini
5. Displays the summary along with timing details

This tool is useful for comparing how different AI models respond to the same input, benchmarking response times, and getting a consolidated view of different AI capabilities.

## Requirements

- Go 1.18 or later
- API keys for:
  - OpenAI API
  - Google Gemini API
  - Anthropic Claude API

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/xssdoctor/haddixscript.git
   cd haddixscript
   ```

2. Install the required Go dependencies:

   ```
   go mod download
   ```

3. API Key Configuration:

   The application will automatically create a config directory at `~/.config/multidoc/multidoc/` with a default `.env` file on first run.

   Edit this file to add your API keys:

   ```
   ~/.config/multidoc/multidoc/.env
   ```

   Add your API keys:

   ```
   OPENAI_API_KEY=your_openai_api_key
   GEMINI_API_KEY=your_gemini_api_key
   CLAUDE_API_KEY=your_claude_api_key
   ```

   Alternatively, you can create a `.env` file in the project directory with the same format.

## Usage

You can pipe text to the program or use interactive input:

```bash
# Using a pipe
echo "Compare the benefits of REST vs GraphQL" | go run app.go

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
- A summary comparing all model responses
- Timing information for each model and the total execution

## Build

To build an executable:

```bash
go build -o haddixscript
```

Then you can run it directly:

```bash
./haddixscript
```

When you run the application for the first time, it will check if the configuration directory exists at `~/.config/multidoc/multidoc/` and create it if necessary, along with a template `.env` file. The application will then exit, allowing you to add your API keys before running it again.

## License

MIT

## Author

xssdoctor
