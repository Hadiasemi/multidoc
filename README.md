# multidoc

A Go-based utility that processes input through multiple AI models concurrently (OpenAI, Claude, and Gemini) and provides a combination of their responses. Keeping each one's individual best parts. This is useful for research and is sometimes refferred to as multipromting or metaprompting. 

## Description

Multidoc allows you to send the same prompt to different AI models simultaneously and combine their responses. The tool:

1. Takes text input from stdin
2. Processes the input through multiple AI models in parallel:
   - GPT-4o (OpenAI)
   - o3-mini (OpenAI)
   - o1-mini (OpenAI)
   - Gemini 2.0 Flash (Google)
   - Claude 3.7 Sonnet (Anthropic)
3. Collects responses from all models with timing information
4. Generates a combination of all model outputs using o1-mini
5. Displays the output along with timing details

## Requirements

- Go 1.18 or later
- API keys for:
  - OpenAI API
  - Google Gemini API
  - Anthropic Claude API

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
   CLAUDE_API_KEY=your_claude_api_key
   ```

   Alternatively, you can create a `.env` file in the project directory with the same format.

## Usage

You can pipe text to the program or use interactive input:

```bash
# Using a pipe
echo "Compare the benefits of REST vs GraphQL" | go run multidoc.go

# Using interactive input
go run multidoc.go
# Then type your prompt and press Ctrl+D when finished
```

For larger prompts, you can use a text file:

```bash
cat prompt.txt | go run multidoc.go
```

The output will display:

- Progress updates as each model responds
- A combination of the models best outputs
- Timing information for each model and the total execution

## Build

To build an executable:

```bash
go build -o multidoc
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
