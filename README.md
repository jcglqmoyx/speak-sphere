[简体中文](README_zh.md) | **English**

# Speak Sphere 🌍 Language Learning Platform

**Speak Sphere** is a modern language learning application that integrates "BuBei Words" and "QingTing English" functionalities, built with Rust and the egui framework to provide a smooth desktop learning experience.

## 🌟 Features

### 📚 Vocabulary Module (Integrated BuBei Words)
- **Smart Memory System**: Based on Ebbinghaus forgetting curve
- **Multi-mode Learning**: Flash cards, spelling practice, dictation tests
- **Rich Vocabulary**: Covers English words for all learning stages
- **Pronunciation Support**: Standard English pronunciation practice
- **Progress Tracking**: Visualized learning progress

### 🎧 Audiobook Reading (Integrated QingTing English)
- **Immersive Reading**: Read while listening to improve language sense
- **Sync Subtitles**: Real-time text display
- **Voice Control**: Play/pause/speed adjustment support
- **Word Marking**: One-click word marking and study plan integration
- **Multi-language Support**: English, Japanese, French, and more

## 🚀 Quick Start

### System Requirements
- **OS**: Windows 10/11, macOS 10.15+, or Linux
- **Memory**: At least 4GB RAM
- **Storage**: At least 100MB free space

### Installation

#### Build from Source
```bash
# Clone the repository
git clone https://github.com/jcglqmoyx/speak-sphere.git
cd speak-sphere

# Build the project
cargo build --release

# Run the application
cargo run
```

#### Download Pre-built Binaries
Visit [Releases page](https://github.com/jcglqmoyx/speak-sphere/releases) to download pre-built binaries for your platform.

### User Guide

1. **Launch Application**: Double-click the app icon or run `target/release/speak-sphere`
2. **Select Feature**: Choose "Vocabulary" or "Audiobook" from the main interface
3. **Start Learning**: Follow instructions to enter corresponding learning modules
4. **Customize Settings**: Adjust learning parameters as needed

## 🛠️ Development Guide

### Tech Stack
- **Backend**: Rust
- **Frontend**: egui (immediate mode GUI)
- **Audio**: Audio playback library (to be integrated)
- **Data Storage**: SQLite (planned)

### Development Environment Setup
```bash
# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Clone the project
git clone https://github.com/jcglqmoyx/speak-sphere.git
cd speak-sphere

# Run development version
cargo run

# Run tests
cargo test
```

### Project Structure
```
speak-sphere/
├── src/
│   ├── main.rs          # Main application entry
│   ├── config.rs        # Configuration management
│   └── utils/           # Utility modules
│       ├── mod.rs
│       ├── font.rs      # Font management
│       └── display_size.rs
├── Cargo.toml          # Dependency configuration
└── README.md          # Project documentation
```

## 📦 Dependencies

Main dependencies:
- `egui` - Immediate mode GUI framework
- `eframe` - egui application framework
- `rdev` - System interaction functionality
- `once_cell` - Single initialization container

## 🤝 Contributing

We welcome all forms of contributions!

### How to Contribute
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Standards
- Follow Rust official coding standards
- Use English for commit messages
- Include test cases for new features
- Ensure all tests pass

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Contact

- **Project Homepage**: [GitHub Repository](https://github.com/jcglqmoyx/speak-sphere)
- **Issue Tracking**: [Issues](https://github.com/jcglqmoyx/speak-sphere/issues)
- **Email**: jcglqmoyx@gmail.com

## 🙏 Acknowledgments

Thanks to the following open source projects:
- [egui](https://github.com/emilk/egui) - Excellent immediate mode GUI framework
- [BuBei Words](https://www.bbdc.net/) - Vocabulary learning concept reference
- [QingTing English](https://www.qt-ing.com/) - Audiobook learning concept reference

---

**Speak Sphere** - Making language learning more efficient and fun! 🚀
