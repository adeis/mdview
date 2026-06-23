<script>
  import { onMount } from 'svelte';

  export let markdown = "";
  export let placeholder = "Start typing your markdown here...";

  let textarea;
  let gutter;

  $: linesCount = markdown.split('\n').length;
  $: lines = Array.from({ length: Math.max(linesCount, 1) }, (_, i) => i + 1);

  function handleScroll() {
    if (textarea && gutter) {
      gutter.scrollTop = textarea.scrollTop;
    }
  }

  function handleKeyDown(e) {
    // Enable Tab indentation inside textarea
    if (e.key === 'Tab') {
      e.preventDefault();
      const start = textarea.selectionStart;
      const end = textarea.selectionEnd;
      
      // Insert 4 spaces
      markdown = markdown.substring(0, start) + "    " + markdown.substring(end);
      
      // Update cursor position after Svelte updates DOM
      setTimeout(() => {
        textarea.selectionStart = textarea.selectionEnd = start + 4;
      }, 0);
    }
  }

  // Ensure scroll is synced on mount
  onMount(() => {
    handleScroll();
  });
</script>

<div class="editor-container">
  <div class="line-gutter" bind:this={gutter}>
    {#each lines as line}
      <div class="line-number">{line}</div>
    {/each}
  </div>
  
  <textarea
    bind:this={textarea}
    bind:value={markdown}
    {placeholder}
    on:scroll={handleScroll}
    on:keydown={handleKeyDown}
    spellcheck="false"
  />
</div>

<style>
  .editor-container {
    display: flex;
    width: 100%;
    height: 100%;
    background-color: var(--bg-panel);
    border-right: 1px solid var(--border-color);
    font-family: 'Fira Code', 'Courier New', Courier, monospace;
    font-size: 14px;
    line-height: 1.6;
    overflow: hidden;
    position: relative;
  }

  .line-gutter {
    width: 48px;
    height: 100%;
    padding: 1.5rem 0;
    overflow: hidden;
    text-align: right;
    background-color: rgba(0, 0, 0, 0.02);
    border-right: 1px solid var(--border-color);
    user-select: none;
    color: var(--text-secondary);
    opacity: 0.6;
  }

  .dark .line-gutter {
    background-color: rgba(255, 255, 255, 0.02);
  }

  .line-number {
    padding-right: 12px;
    height: 22.4px; /* match line-height * font-size */
  }

  textarea {
    flex: 1;
    height: 100%;
    border: none;
    outline: none;
    resize: none;
    background: transparent;
    color: var(--text-primary);
    padding: 1.5rem 1rem;
    font-family: inherit;
    font-size: inherit;
    line-height: inherit;
    overflow-y: auto;
    white-space: pre;
  }
</style>
