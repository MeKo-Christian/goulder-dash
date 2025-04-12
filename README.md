# Goulder Dash

**A nutz Boulder Dash clone written in Go.**  
Powered by [gonutz/prototype](https://github.com/gonutz/prototype) for pixel-perfect nostalgia!

---

## 🎮 About

**Goulder Dash** is a classic-style Boulder Dash clone written in pure Go using the `prototype` game framework.  
Guide your explorer through dirt-filled caves, collect sparkling gems, avoid falling rocks, and reach the exit!

---

## 🚀 Features

- 🪨 Falling rock physics with timed updates
- 💎 Gem collection
- 🎮 Arrow key movement and rock pushing
- 🧱 Procedural and handcrafted levels
- 🕹 Minimal dependencies – pure Go fun

---

## 🧰 Requirements

- Go 1.21+
- A C compiler (required by `prototype`'s GLFW backend)
- `libglfw3` and `libglfw3-dev` (on Linux)

```bash
# Ubuntu/Debian
sudo apt install libglfw3 libglfw3-dev
```

---

## 🛠 How to Run

```bash
git clone https://github.com/yourusername/goulder-dash
cd goulder-dash
go run .
```

---

## 🎨 Assets

The game uses a custom sprite sheet located in assets/sprites.png.
Sprites are 64×64 pixels, scaled to 32×32 during gameplay.

---

## 📦 Future Plans

- 💥 Add explosions when crushed
- 🤖 Add simple enemy AI
- ⛏ Improve digging logic and particle effects
- 🎶 Add sound effects and retro music
- 🗺 Multiple levels and level transitions

---

## ✨ Credits

- Inspired by the original Boulder Dash (1984)
- Built using [gonutz/prototype](https://github.com/gonutz/prototype)
- Code by [Christian-W. Budde](https://pcjv.de)
- Sprites created by ChatGPT and adapted.

---

## 🧪 License

MIT License.
Use the code, expand it, or make your own nutz clone.
