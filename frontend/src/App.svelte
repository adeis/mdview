<script>
  import { onMount } from 'svelte';
  import { 
    GetInitialFile, 
    SaveFile, 
    OpenFileDialog, 
    SaveFileDialog,
    ExportPdfDialog
  } from '../wailsjs/go/main/App.js';
  import Editor from './components/Editor.svelte';
  import Viewer from './components/Viewer.svelte';
  import { marked } from 'marked';
  
  import { jsPDF } from 'jspdf';
  import html2canvas from 'html2canvas';
  window.html2canvas = html2canvas;

  // Configure marked options
  marked.setOptions({
    gfm: true,
    breaks: true,
    headerIds: true,
    mangle: false
  });
  
  // Icon Imports
  import { 
    FolderOpen, 
    Save, 
    FilePlus, 
    Eye, 
    Edit3, 
    Columns, 
    Sun, 
    Moon, 
    FileText, 
    Download 
  } from 'lucide-svelte';

  // App States
  let markdown = "";
  let filePath = "";
  let fileName = "Untitled.md";
  let lastSavedMarkdown = "";
  
  let viewMode = "split"; // editor, viewer, split
  let fontClass = "font-serif"; // font-serif, font-sans
  let pdfMode = true;
  let theme = "dark";

  // Statistics
  $: wordCount = markdown ? markdown.trim().split(/\s+/).filter(w => w.length > 0).length : 0;
  $: charCount = markdown ? markdown.length : 0;
  $: isDirty = markdown !== lastSavedMarkdown;

  onMount(async () => {
    // Determine OS theme preference
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    theme = prefersDark ? 'dark' : 'light';
    applyTheme();

    // Check if opened with a file from CLI
    try {
      const fileInfo = await GetInitialFile();
      if (fileInfo && fileInfo.path) {
        filePath = fileInfo.path;
        markdown = fileInfo.content || "";
        fileName = fileInfo.name || "Untitled.md";
        lastSavedMarkdown = markdown;
        
        // Open in Split mode (Editor + Viewer) on file load so editor is visible
        viewMode = "split";
      }
    } catch (e) {
      console.error("Failed to fetch initial file:", e);
    }
  });

  function applyTheme() {
    if (theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }

  function toggleTheme() {
    theme = theme === 'light' ? 'dark' : 'light';
    applyTheme();
  }

  function confirmDiscard() {
    if (isDirty) {
      return confirm("You have unsaved changes. Do you want to discard them?");
    }
    return true;
  }

  function newFile() {
    if (!confirmDiscard()) return;
    markdown = "";
    filePath = "";
    fileName = "Untitled.md";
    lastSavedMarkdown = "";
  }

  async function openFile() {
    if (!confirmDiscard()) return;
    try {
      const fileInfo = await OpenFileDialog();
      if (fileInfo && fileInfo.path) {
        filePath = fileInfo.path;
        markdown = fileInfo.content;
        fileName = fileInfo.name;
        lastSavedMarkdown = fileInfo.content;
      }
    } catch (err) {
      alert("Error opening file: " + err);
    }
  }

  async function saveFile() {
    if (!filePath) {
      await saveFileAs();
      return;
    }
    try {
      await SaveFile(filePath, markdown);
      lastSavedMarkdown = markdown;
    } catch (err) {
      alert("Error saving file: " + err);
    }
  }

  async function saveFileAs() {
    try {
      const fileInfo = await SaveFileDialog(markdown);
      if (fileInfo && fileInfo.path) {
        filePath = fileInfo.path;
        fileName = fileInfo.name;
        lastSavedMarkdown = markdown;
      }
    } catch (err) {
      alert("Error saving file: " + err);
    }
  }

  // Reactive parsed HTML for the hidden PDF export container
  $: parsedHtml = marked.parse(markdown || "*No content*");

  async function exportPdf() {
    const element = document.getElementById('pdf-export-container');
    if (!element) {
      alert("Export container not found!");
      return;
    }

    try {
      // Ensure all web fonts are fully loaded before capturing
      await document.fonts.ready;

      // 1. Detect elements crossing page boundaries and insert HTML/table row spacers
      const pageHeightPx = 1123; // A4 height at 96 DPI
      const pageMarginPx = 54; // standard top/bottom page margin
      const maxContentHeight = pageHeightPx - (pageMarginPx * 2); // 1015px printable area

      // Find critical block elements to avoid splitting: table rows, code pre blocks, quotes, and headings
      const elementsToAvoid = Array.from(element.querySelectorAll('tr, pre, blockquote, h1, h2, h3'));
      const spacers = [];

      for (let i = 0; i < elementsToAvoid.length; i++) {
        const el = elementsToAvoid[i];
        const rect = el.getBoundingClientRect();
        const containerRect = element.getBoundingClientRect();
        const elTop = rect.top - containerRect.top;
        const elBottom = elTop + rect.height;

        // Skip elements that are too tall to fit on a single page anyway
        if (rect.height >= maxContentHeight) continue;

        const currentPage = Math.floor(elTop / pageHeightPx);
        const pageBoundary = (currentPage + 1) * pageHeightPx;

        // If the element overlaps with the A4 page boundary, push it to the next page
        if (elTop < pageBoundary && elBottom > pageBoundary) {
          const spacerHeight = (pageBoundary + pageMarginPx) - elTop;
          let spacer;

          if (el.tagName === 'TR') {
            // Valid HTML spacer for table rows: insert a new row with an empty tall cell
            spacer = document.createElement('tr');
            spacer.className = 'pdf-page-spacer';
            const td = document.createElement('td');
            td.colSpan = el.cells.length; // span all columns
            td.style.height = `${spacerHeight}px`;
            td.style.border = 'none';
            td.style.padding = '0';
            spacer.appendChild(td);
          } else {
            // Standard block spacer for divs/paragraphs/headings
            spacer = document.createElement('div');
            spacer.style.height = `${spacerHeight}px`;
            spacer.style.clear = 'both';
            spacer.className = 'pdf-page-spacer';
          }

          // Insert spacer in front of the element to push it down
          el.parentNode.insertBefore(spacer, el);
          spacers.push(spacer);
        }
      }

      // 2. Render the off-screen A4 container (now fully paginated) to canvas
      const canvas = await html2canvas(element, {
        scale: 2, // Crisp high-resolution quality
        useCORS: true,
        logging: false,
        backgroundColor: '#ffffff'
      });

      // 3. Remove all spacers immediately after capture to restore original DOM state
      spacers.forEach(spacer => {
        if (spacer.parentNode) {
          spacer.parentNode.removeChild(spacer);
        }
      });

      // 4. Slice the high-res canvas image into A4 PDF pages
      const imgData = canvas.toDataURL('image/jpeg', 0.95);
      
      const doc = new jsPDF({
        orientation: 'portrait',
        unit: 'pt',
        format: 'a4'
      });

      const imgWidth = 595.28;
      const pageHeight = 841.89;
      const imgHeight = (canvas.height * imgWidth) / canvas.width;
      
      let heightLeft = imgHeight;
      let position = 0;

      // Render page 1
      doc.addImage(imgData, 'JPEG', 0, position, imgWidth, imgHeight);
      heightLeft -= pageHeight;

      // Slice remaining height into subsequent pages
      while (heightLeft > 0) {
        position = heightLeft - imgHeight;
        doc.addPage();
        doc.addImage(imgData, 'JPEG', 0, position, imgWidth, imgHeight);
        heightLeft -= pageHeight;
      }

      // 5. Save generated PDF
      const pdfBase64 = doc.output('datauristring').split(',')[1];
      const defaultPdfName = fileName.replace(/\.md$/, '') + '.pdf';
      const fileInfo = await ExportPdfDialog(pdfBase64, defaultPdfName);
      if (fileInfo && fileInfo.path) {
        console.log("PDF exported successfully to: " + fileInfo.path);
      }
    } catch (err) {
      alert("Failed to export PDF: " + err);
    }
  }
</script>

<div class="app-layout">
  <!-- Top Navigation Toolbar -->
  <header class="glass no-print">
    <div class="brand">
      <FileText size={20} class="brand-icon" />
      <span class="brand-name">mdviewer</span>
      <span class="file-badge" class:dirty={isDirty}>
        {fileName}{isDirty ? ' • modified' : ''}
      </span>
    </div>
    
    <div class="actions">
      <!-- File controls -->
      <button class="btn" on:click={newFile} title="New Document">
        <FilePlus size={15} />
        <span>New</span>
      </button>
      
      <button class="btn" on:click={openFile} title="Open Document">
        <FolderOpen size={15} />
        <span>Open</span>
      </button>
      
      <button class="btn btn-primary" on:click={saveFile} title="Save Document">
        <Save size={15} />
        <span>Save</span>
      </button>
      
      <div class="divider" />
      
      <!-- View controls -->
      <div class="segmented-control">
        <button 
          class="control-btn" 
          class:active={viewMode === 'viewer'} 
          on:click={() => viewMode = 'viewer'} 
          title="Reader Mode"
        >
          <Eye size={15} />
          <span>Read</span>
        </button>
        <button 
          class="control-btn" 
          class:active={viewMode === 'split'} 
          on:click={() => viewMode = 'split'} 
          title="Split View"
        >
          <Columns size={15} />
          <span>Split</span>
        </button>
        <button 
          class="control-btn" 
          class:active={viewMode === 'editor'} 
          on:click={() => viewMode = 'editor'} 
          title="Editor Mode"
        >
          <Edit3 size={15} />
          <span>Write</span>
        </button>
      </div>
      
      <div class="divider" />
      
      <!-- Display preferences -->
      <button 
        class="btn font-toggle" 
        on:click={() => fontClass = fontClass === 'font-serif' ? 'font-sans' : 'font-serif'} 
        title="Toggle Reader Font (Serif/Sans-serif)"
      >
        <span>{fontClass === 'font-serif' ? 'Serif' : 'Sans'}</span>
      </button>
      
      <button 
        class="btn" 
        class:active={pdfMode} 
        on:click={() => pdfMode = !pdfMode} 
        title="Toggle PDF Paper Layout"
      >
        <span>Paper</span>
      </button>
      
      <button class="btn" on:click={exportPdf} title="Export to PDF">
        <Download size={15} />
        <span>Export</span>
      </button>
      
      <div class="divider" />
      
      <!-- Theme toggle -->
      <button class="btn theme-toggle" on:click={toggleTheme} title="Toggle Light/Dark Theme">
        {#if theme === 'dark'}
          <Sun size={15} />
        {:else}
          <Moon size={15} />
        {/if}
      </button>
    </div>
  </header>

  <!-- Workspace split panel area -->
  <main class="workspace">
    {#if viewMode === 'editor' || viewMode === 'split'}
      <div class="pane editor-pane" class:full-width={viewMode === 'editor'}>
        <Editor bind:markdown />
      </div>
    {/if}
    
    {#if viewMode === 'viewer' || viewMode === 'split'}
      <div class="pane viewer-pane" class:full-width={viewMode === 'viewer'}>
        <Viewer {markdown} {fontClass} {pdfMode} />
      </div>
    {/if}
  </main>

  <!-- Bottom Status Bar -->
  <footer class="status-bar glass no-print">
    <div class="path-info">
      <span>{filePath || 'Unsaved Document'}</span>
    </div>
    <div class="doc-stats">
      <span class="stat-item">Words: <strong>{wordCount}</strong></span>
      <span class="divider-v" />
      <span class="stat-item">Characters: <strong>{charCount}</strong></span>
    </div>
  </footer>

  <!-- Hidden off-screen container for PDF rendering -->
  <div id="pdf-export-container" class="markdown-body font-serif no-print">
    {@html parsedHtml}
  </div>
</div>

<style>
  .app-layout {
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
    overflow: hidden;
  }

  header {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 1.5rem;
    border-bottom: 1px solid var(--border-color);
    z-index: 100;
  }

  .brand {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  :global(.brand-icon) {
    color: var(--accent);
  }

  .brand-name {
    font-weight: 700;
    font-size: 1.1rem;
    letter-spacing: -0.025em;
  }

  .file-badge {
    font-size: 0.75rem;
    font-weight: 500;
    background-color: var(--border-color);
    color: var(--text-secondary);
    padding: 0.2rem 0.6rem;
    border-radius: 9999px;
    margin-left: 0.5rem;
    transition: all 0.2s ease;
  }

  .file-badge.dirty {
    background-color: var(--accent-light);
    color: var(--accent);
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .btn {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    background-color: transparent;
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    padding: 0.45rem 0.8rem;
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn:hover {
    background-color: var(--border-color);
  }

  .btn.active {
    background-color: var(--accent);
    color: white;
    border-color: var(--accent);
  }

  .btn-primary {
    background-color: var(--accent);
    color: white;
    border-color: var(--accent);
  }

  .btn-primary:hover {
    background-color: var(--accent-hover);
    border-color: var(--accent-hover);
  }

  .font-toggle {
    font-family: inherit;
    font-weight: 600;
    min-width: 55px;
    justify-content: center;
  }

  .divider {
    width: 1px;
    height: 20px;
    background-color: var(--border-color);
    margin: 0 0.5rem;
  }

  /* Segmented view switcher control */
  .segmented-control {
    display: flex;
    background-color: var(--bg-app);
    border: 1px solid var(--border-color);
    padding: 2px;
    border-radius: 8px;
  }

  .control-btn {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    border: none;
    background: transparent;
    color: var(--text-secondary);
    padding: 0.4rem 0.8rem;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .control-btn.active {
    background-color: var(--bg-panel);
    color: var(--text-primary);
    box-shadow: var(--shadow-sm);
  }

  /* Main editor and viewer space */
  .workspace {
    flex: 1;
    display: flex;
    width: 100%;
    overflow: hidden;
  }

  .pane {
    height: 100%;
    overflow: hidden;
    transition: width 0.3s ease;
  }

  .editor-pane {
    width: 50%;
  }

  .viewer-pane {
    width: 50%;
  }

  .pane.full-width {
    width: 100% !important;
  }

  /* Footer Status Bar */
  .status-bar {
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 1.5rem;
    border-top: 1px solid var(--border-color);
    font-size: 0.75rem;
    color: var(--text-secondary);
  }

  .path-info {
    font-family: 'Fira Code', monospace;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 60%;
  }

  .doc-stats {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .divider-v {
    width: 1px;
    height: 12px;
    background-color: var(--border-color);
  }

  /* Print optimizations */
  @media print {
    .app-layout {
      height: auto !important;
      overflow: visible !important;
    }
    
    .workspace {
      display: block !important;
      height: auto !important;
      overflow: visible !important;
    }

    .viewer-pane {
      width: 100% !important;
      height: auto !important;
      overflow: visible !important;
    }
  }
</style>
