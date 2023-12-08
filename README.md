# Party Image Bot

## Overview

This application provides a complete server and client slide implementation of a bot that will accept an SMS prompt for an image, render with DALL-E, then ask for an additional emblisshed prompt to be created by ChatGPT4, and then rendered again with DALL-E.

## Example

![Cats](/examples/example-output.png)

## Architecture

The approach is to provide a self contained stack that handle the client and server. To this end I used [go-app.dev](https://go-app.dev) to build a Progressive Web App (PWA) that builds a Web Assembly (WASM) component for the browser. It may receive an SMS via a POST request on the API end point api/v1/sms and then uses websockets on api/v1/ws to update the generated content to the browser session.

## Required environment variable for OpenAI Developer key

Obtain an developer API key from OpenAI and set an enviornment variable 'OPENAI_API_KEY'.