<template>
  <div class="hero-install-section">
    <div class="install-box">
      <div class="install-code" @click="copyToClipboard">
        <span class="install-prompt">$</span>
        <span class="install-command">curl -fsSL https://trymist.cloud/install.sh | bash</span>
        <button class="copy-button" :class="{ copied }" @click.stop="copyToClipboard">
          <span v-if="!copied">📋</span>
          <span v-else>✓</span>
        </button>
      </div>
    </div>
    <div class="hero-actions">
      <a href="/guide/getting-started" class="action-button primary">Get Started</a>
      <a href="https://github.com/corecollectives/mist" class="action-button secondary" target="_blank">View on GitHub</a>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const copied = ref(false)

const copyToClipboard = () => {
  const command = 'curl -sSL https://trymist.cloud/install.sh | bash'
  navigator.clipboard.writeText(command).then(() => {
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  })
}
</script>

<style scoped>
.hero-install-section {
  max-width: 800px;
  margin: 2rem auto 0;
  padding: 0 1.5rem;
}

.install-box {
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-border);
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.install-box:hover {
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.install-code {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 1rem;
  cursor: pointer;
  user-select: none;
  position: relative;
}

.install-prompt {
  color: var(--vp-c-brand-1);
  font-weight: bold;
  opacity: 0.8;
}

.install-command {
  flex: 1;
  color: var(--vp-c-text-1);
  font-weight: 500;
}

.copy-button {
  background: var(--vp-c-brand-1);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 48px;
}

.copy-button:hover {
  background: var(--vp-c-brand-2);
  transform: scale(1.05);
}

.copy-button.copied {
  background: #10b981;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.action-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 1rem;
  text-decoration: none;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.action-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.action-button.primary {
  background: var(--vp-c-brand-1);
  color: white;
}

.action-button.primary:hover {
  background: var(--vp-c-brand-2);
}

.action-button.secondary {
  background: var(--vp-c-bg-soft);
  color: var(--vp-c-text-1);
  border: 1px solid var(--vp-c-border);
}

.action-button.secondary:hover {
  border-color: var(--vp-c-brand-1);
  color: var(--vp-c-brand-1);
}

@media (max-width: 768px) {
  .install-code {
    font-size: 0.875rem;
    flex-wrap: wrap;
  }
  
  .copy-button {
    position: absolute;
    right: 0;
    top: 0;
  }
  
  .install-command {
    padding-right: 60px;
  }
}
</style>
