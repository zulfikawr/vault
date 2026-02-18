import { defineStore } from 'pinia';

export type Theme = 'gruvbox' | 'monokai' | 'nord' | 'dracula' | 'solarized';

interface UIState {
  theme: Theme;
  sidebarCollapsed: boolean;
  isMobileMenuOpen: boolean;
  compactMode: boolean;
}

export const useUIStore = defineStore('ui', {
  state: (): UIState => ({
    theme: (localStorage.getItem('theme') as Theme) || 'gruvbox',
    sidebarCollapsed: localStorage.getItem('sidebarCollapsed') === 'true',
    isMobileMenuOpen: false,
    compactMode: localStorage.getItem('compactMode') === 'true',
  }),
  actions: {
    setTheme(theme: Theme) {
      this.theme = theme;
      localStorage.setItem('theme', theme);
      document.documentElement.className = `theme-${theme}`;
    },
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed;
      localStorage.setItem('sidebarCollapsed', String(this.sidebarCollapsed));
    },
    toggleCompactMode() {
      this.compactMode = !this.compactMode;
      localStorage.setItem('compactMode', String(this.compactMode));
    },
    toggleMobileMenu() {
      this.isMobileMenuOpen = !this.isMobileMenuOpen;
    },
    closeMobileMenu() {
      this.isMobileMenuOpen = false;
    },
    initTheme() {
      document.documentElement.className = `theme-${this.theme}`;
    }
  },
});
