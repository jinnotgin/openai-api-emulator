# OpenAI API Emulator

This repository contains a simple server written in Go that emulates the behaviour of the OpenAI API. It's a great tool for testing your applications that interact with the OpenAI API without making actual API calls.

## Description

The server is designed to mimic the OpenAI's chat API for large language models . It generates fake responses to API calls, allowing you to test your application's handling of API responses without incurring any cost. The server supports both regular and streaming responses.

The server is not meant to replicate the AI capabilities of the OpenAI API. It generates static responses and is intended purely for testing purposes.

## Features

- Generates unique IDs for each response similar to the OpenAI API.
- Supports both regular and streaming responses.
- Estimates the number of tokens used in the prompt. 
- Verbose logging for debugging purposes.

## Usage

To start the server, simply run `go run server.go`. The server will start listening on port 8383.

## API Endpoints

The server supports the following API endpoint:

- `/v1/chat/completions`: This endpoint mimics the chat models of the OpenAI API (https://platform.openai.com/docs/api-reference/chat/create). It accepts POST requests with a JSON body containing the chat messages and optional parameters like model and temperature. The server responds with a JSON object that mimics the structure of the OpenAI API responses.

## Response Structure

The response from the server has the following structure (assuming a regular non-streaming response):

```json
{
  "id": <random_id>,
  "choices": [
    {
      "finish_reason": "STOP",
      "index": 0,
      "message": {
        "content": "Blank response from OpenAI API emulator.",
        "role": "assistant"
      }
    }
  ],
  "created": <unix_timestamp>,
  "model": "gpt-3.5-turbo",
  "object": "chat.completion",
  "usage": {
    "prompt_tokens": <estimated_token_count>,
    "completion_tokens": 8,
    "total_tokens": <total_token_count>
  }
}
```

Please note that the server always returns a predefined message in the response and the message `content` field does not depend on the input message.

## License

This project is licensed under the MIT License.

Please note: This server does not use the OpenAI API and is not affiliated with OpenAI.
