# Welcome to mdviewer! 🚀

This is a beautiful, premium, PDF-like Markdown viewer and editor built with **Golang** and **Svelte** (via **Wails**).

> **Aesthetic Design:** With modern typography, glassmorphism UI, a dual-panel layout, and dark/light modes, this app is built to help you write and read books or blog posts with maximum comfort.

---

## Key Features

1. **Dual Panel Layout**
    - **Read Mode:** Fullscreen reading view with PDF-like sheet padding.
    - **Split View:** Side-by-side editing and real-time rendering.
    - **Write Mode:** Focus-driven editor with line numbers and custom tab indentation.
2. **Typography Options**
    - Toggle between *Serif* (`Lora` for classic book reading) and *Sans* (`Inter` for technical blogs).
    - Custom line spacing and spacing parameters.
3. **Save and Export**
    - Save files directly to disk.
    - **Export to PDF:** Powered by native OS print styles (hides UI toolbars automatically).

---

## Sample Markdown Elements

### Code Block with Highlights

```javascript
// Syncing gutter scrolling with editor textarea
function handleScroll() {
  if (textarea && gutter) {
    gutter.scrollTop = textarea.scrollTop;
  }
}
```

### Alternating Data Table

| Mode | Purpose | Layout |
| :--- | :--- | :--- |
| **Reader** | Relaxed reading (book layout) | Single panel |
| **Split** | Writing & Previewing | 50/50 split panel |
| **Writer** | Pure content creation | Fullscreen text editor |

### Quote Box

> "Writing is the painting of the voice."
> — *Voltaire*

### Checklist
- [x] Write Go backend logic
- [x] Style Svelte frontend theme
- [x] Package into single executable app
- [ ] Write my first book! 📚

Feel free to click **Open** to load other documents or **Export** to save this as a beautiful PDF!
