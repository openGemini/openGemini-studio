# openGemini Studio

<p align="center">
  <a href="README.md">English</a> | <a href="README_zh-CN.md">简体中文</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/Vue.js-3.x-4FC08D?logo=vue.js" alt="Vue.js">
  <img src="https://img.shields.io/badge/TypeScript-5.x-3178C6?logo=typescript" alt="TypeScript">
</p>

A modern, cross-platform desktop client for [openGemini](https://github.com/openGemini/openGemini) time-series database. Built with Wails, Vue 3, and TypeScript to provide a native desktop experience for database management, query execution, and data visualization.

## 📋 Table of Contents

- [Features](#-features)
- [Screenshots](#-screenshots)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#-usage)
- [Building from Source](#-building-from-source)
- [Development](#-development)
- [Architecture](#-architecture)
- [Contributing](#-contributing)
- [License](#-license)
- [Community](#-community)
- [Acknowledgments](#-acknowledgments)

## ✨ Features

### Core Functionality
- **Multi-Connection Management**: Connect to multiple openGemini instances simultaneously with support for HTTP/HTTPS protocols
- **SQL Query Editor**: Advanced InfluxQL query editor with syntax highlighting and auto-completion
- **Query Execution**: Execute queries with real-time results display and performance metrics
- **Data Visualization**: View query results in a customizable table format with resizable columns
- **Query History**: Persistent query history with success/failure tracking and execution time

### Database Operations
- **Database Explorer**: Browse databases, measurements, retention policies, and fields in a tree view
- **Metadata Management**: View and manage database metadata including retention policies
- **Data Writing**: Support for INSERT operations with line protocol
- **Batch Operations**: Execute multiple queries in sequence

### User Experience
- **Cross-Platform**: Native desktop application for Windows, macOS, and Linux
- **Internationalization**: Built-in support for English and Chinese (zh-CN)
- **Theme Support**: Light, dark, and system-following theme modes
- **Customizable**: Configurable fonts, history limits, and data directories
- **Debug Mode**: Optional debug mode for detailed logging and diagnostics

### Security & Connectivity
- **TLS/SSL Support**: Secure connections with custom CA certificates
- **Certificate Management**: Client certificate and key support
- **Authentication**: Username/password authentication
- **Flexible Configuration**: Per-connection security settings

## 📸 Screenshots

![image](images/en.png)

## 🚀 Getting Started

### Prerequisites

**For End Users:**
- No prerequisites needed - download and run the pre-built binary for your platform

**For Developers:**
- Go 1.19 or higher
- Node.js 16 or higher
- npm or yarn
- Wails CLI v2

### Installation

#### Download Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/openGemini/openGemini-studio/releases) page:

- **Windows**: `openGemini-studio-windows-amd64.exe`
- **macOS**: `openGemini-studio-darwin-universal`
- **Linux**: `openGemini-studio-linux-amd64`

#### Install from Source

```bash
# Clone the repository
git clone https://github.com/openGemini/openGemini-studio.git
cd openGemini-studio

# Install Wails CLI (if not already installed)
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Install dependencies and build
wails build
```

## 💡 Usage

### Creating a Connection

1. Launch openGemini Studio
2. Click the "Add Connection" button in the sidebar
3. Enter connection details:
   - **Connection Name**: Friendly name for the connection
   - **Address**: openGemini server address (e.g., `localhost:8086`)
   - **Protocol**: HTTP or HTTPS
   - **Authentication**: Enable and provide credentials if required
4. Click "Test Connection" to verify
5. Click "Save Connection" to persist

### Executing Queries

1. Select a database from the connection tree
2. Enter your InfluxQL query in the editor
3. Click "Execute Query" or press `Ctrl+Enter` (Windows/Linux) or `Cmd+Enter` (macOS)
4. View results in the table below

### Managing Query History

- Click the "Query History" button to view past queries
- Click on a history item to append it to the current query editor
- Clear history using the "Clear History" button

### Configuring Settings

Access settings via the gear icon:
- **Language**: Switch between English and Chinese
- **Theme**: Choose light, dark, or system theme
- **Custom Font**: Set a custom font family
- **Max History Count**: Configure the number of queries to retain (10-500)
- **Debug Mode**: Enable detailed logging for troubleshooting

## 🔨 Building from Source

### Prerequisites

Install the following dependencies:

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verify installation
wails doctor
```

### Build Commands

```bash
# Development build
wails build

# Production build with optimizations
wails build -clean

# Build for specific platform
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64
```

### Output

Built binaries are placed in the `build/bin/` directory.

## 🛠 Development

### Running in Development Mode

```bash
# Start development server with hot reload
wails dev
```

This runs a Vite development server with hot module replacement. The application will automatically reload when you modify frontend code. Backend changes require a restart.

### Project Structure

```
openGemini-studio/
├── app.go                  # Main application logic
├── data_struct.go          # Data models and structures
├── http_client.go          # OpenGemini HTTP client
├── database.go             # BoltDB persistence layer
├── logger.go               # Logging infrastructure
├── frontend/               # Vue.js frontend
│   ├── src/
│   │   ├── components/     # Vue components
│   │   ├── composables/    # Vue composables
│   │   ├── locales/        # i18n translations
│   │   └── types/          # TypeScript types
│   └── wailsjs/            # Generated Wails bindings
└── build/                  # Build output directory
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run Vite dev server
npm run dev

# Build frontend
npm run build

# Type checking
npm run type-check
```

### Backend Development

The backend is written in Go and uses:
- **Wails v2**: Desktop application framework
- **BoltDB**: Embedded database for configuration storage
- **openGemini Go Client**: Official openGemini client library

## 🏗 Architecture

### Technology Stack

**Frontend:**
- Vue 3 (Composition API)
- TypeScript
- Vue I18n (internationalization)
- Tailwind CSS

**Backend:**
- Go 1.24+
- Wails v2
- BoltDB (embedded database)
- openGemini Go client

### Data Flow

1. User interacts with Vue.js frontend
2. Frontend calls Go methods via Wails bindings
3. Go backend communicates with openGemini server
4. Results are returned to frontend for display
5. Configuration and history are persisted to BoltDB

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

### Ways to Contribute

- 🐛 **Report Bugs**: Open an issue with details about the bug
- 💡 **Suggest Features**: Share your ideas in the discussions
- 📝 **Improve Documentation**: Help us improve README, guides, and code comments
- 🔧 **Submit Pull Requests**: Fix bugs or implement new features

### Contribution Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Commit your changes (`git commit -s -S -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Style

- **Go**: Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- **TypeScript/Vue**: Follow [Vue.js Style Guide](https://vuejs.org/style-guide/)
- **Commits**: Use conventional commit messages

## 📄 License

Copyright 2025 openGemini Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## 🌐 Community

### Get Help

- **GitHub Issues**: [Report bugs or request features](https://github.com/openGemini/openGemini-studio/issues)
- **Discussions**: [Ask questions and share ideas](https://github.com/openGemini/openGemini-studio/discussions)
- **openGemini Community**: [Join the main openGemini community](https://github.com/openGemini/openGemini)

### Stay Updated

- Watch this repository for updates
- Star the project if you find it useful
- Follow [@openGemini](https://github.com/openGemini) on GitHub

## 🙏 Acknowledgments

This project is built with amazing open source technologies:

- [openGemini](https://github.com/openGemini/openGemini) - High-performance time-series database
- [Wails](https://wails.io/) - Build desktop apps using Go & Web Technologies
- [Vue.js](https://vuejs.org/) - The Progressive JavaScript Framework
- [TypeScript](https://www.typescriptlang.org/) - JavaScript with syntax for types
- [BoltDB](https://github.com/etcd-io/bbolt) - Embedded key/value database

---

<p align="center">
  Made with ❤️ by the openGemini community
</p>
