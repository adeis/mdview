<script>
  import { marked } from 'marked';
  
  export let markdown = "";
  export let fontClass = "font-serif"; // font-serif or font-sans
  export let pdfMode = true;

  // Configure marked options
  marked.setOptions({
    gfm: true,
    breaks: true,
    headerIds: true,
    mangle: false
  });

  $: html = marked.parse(markdown || "*No content*");
</script>

<div class="viewer-container" class:pdf-mode={pdfMode}>
  <div class="print-page" class:pdf-sheet={pdfMode}>
    <article class="markdown-body {fontClass}">
      {@html html}
    </article>
  </div>
</div>

<style>
  .viewer-container {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    padding: 2rem;
    display: flex;
    justify-content: center;
    background-color: var(--pdf-bg);
    transition: background-color 0.3s ease;
  }

  .viewer-container.pdf-mode {
    background-color: var(--pdf-bg);
  }

  /* When NOT in PDF mode, it takes full width without margins */
  .print-page {
    width: 100%;
    max-width: 900px;
    transition: all 0.3s ease;
  }

  /* PDF-like sheet styling */
  .print-page.pdf-sheet {
    background-color: var(--pdf-paper-bg);
    box-shadow: var(--pdf-paper-shadow);
    border-radius: 8px;
    padding: 3rem 4rem;
    min-height: 297mm; /* Standard A4 proportion */
    border: 1px solid var(--border-color);
  }

  /* Make sure printing looks like clean pages */
  @media print {
    .viewer-container {
      padding: 0 !important;
      background: white !important;
      overflow: visible !important;
    }
    .print-page.pdf-sheet {
      box-shadow: none !important;
      border: none !important;
      padding: 0 !important;
      background: transparent !important;
    }
  }

  /* Custom responsiveness */
  @media (max-width: 768px) {
    .viewer-container {
      padding: 1rem;
    }
    .print-page.pdf-sheet {
      padding: 2rem;
    }
  }
</style>
