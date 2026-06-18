# Pure Reader 🎯

Pure Reader is a minimalist, open-source web utility built in Go designed to eliminate digital noise, ads, cookie banners, and clutter from web pages, delivering a clean and distraction-free reading experience.

Built with performance and privacy at its core, it acts as an ephemeral proxy that processes web texts instantly without storing any user data.

## ✨ Features

- **Distraction-Free Reading:** Automatically strips away ads, pop-ups, newsletter forms, and chaotic layouts, leaving only the main text and relevant images.
- **Privacy-First (Zero Logs):** No database attached. The application processes URLs on the fly in memory and immediately purges the data. Your reading history remains entirely yours.
- **System-Aware Theme Toggle:** Supports Dark, Light, and System default themes with smooth transitions.
- **Bilingual Interface:** Toggle seamlessly between English and Turkish with a single click.
- **Lightweight & Blazing Fast:** Powered by Go and Fiber, compiled into a single binary with minimal resource consumption.
- **Built-in Rate Limiting:** Protected against automated abuse and DDoS attempts.

## 🛠️ Built With

- [Go (Golang)](https://go.dev/) - Core language
- [Fiber](https://gofiber.io/) - High-performance web framework
- [go-readability](https://github.com/go-shiori/go-readability) - Metrics-based text extraction algorithm

## 🚀 Getting Started

### Prerequisites

Make sure you have Go installed on your system (version 1.16 or higher).

### Installation & Running Locally

1. Clone the repository:
   ```bash
   git clone [https://github.com/kdmrbss/pure-reader.git](https://github.com/kdmrbss/pure-reader.git)
   cd pure-reader