# ⚡ ZipUp-ReadyToGo

**ZipUp-ReadyToGo** is a lightweight command-line tool built in Go that automates repetitive file preparation tasks — including renaming, organizing, and packaging data — so it's ready for deployment in seconds.

Originally created to eliminate a tedious, manual workflow in my office (I've made this during my intern time) (often involving 50+ repetitive steps), this tool was built during personal time and now serves as a practical portfolio piece showcasing automation and Go CLI development.

---

## ✨ Features

- 📁 Batch renames and organizes files
- 📦 Prepares folders for server compatibility
- ⚙️ Fully automated – no more manual setup
- 🧰 CLI-first design, minimal setup

---

## 🚀 Quick Start

### ✅ Build the executable:

```bash
go build main.go
```
<br>

## 🛠️ Usage

Once compiled, run the tool via CLI:


./ZipUp-ReadyToGo /path/to/your/data

ℹ️ You can pass in a directory path and optionally use a renameList.json to fine-tune behavior.


## 🧩 Why I Built This

This project was born out of necessity: we had to manually prepare and rename batches of files and folders for server processing — sometimes doing this 50+ times in a row.

So, I built ZipUp-ReadyToGo to do all of it automatically — saving time, reducing errors, and turning a painful task into a one-command operation.

## 💡 Suggestions Welcome!

Although development on this application has finished, I’m still interested in hearing your feedback and suggestions. If you have any ideas for how to improve the tool, please feel free to open an issue.

I appreciate your feedback! 🥳

## 👨‍💻 Author

Made with care during my free time :)

## 📄 License

This project is licensed under the GNU General Public License v3.0 (GPLv3).
You are free to use, modify, and distribute the code, but any modified versions must also be open-sourced under the same license.

For full details, see the LICENSE file or visit gnu.org/licenses/gpl-3.0.