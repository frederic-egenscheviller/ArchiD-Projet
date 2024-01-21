import {Injectable, Renderer2, RendererFactory2} from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class ThemeService {
  darkModeMediaQuery: MediaQueryList;
  renderer: Renderer2;

  constructor(rendererFactory: RendererFactory2) {
    this.renderer = rendererFactory.createRenderer(null, null);
    this.darkModeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
  }

  isDarkMode(): boolean {
    return this.darkModeMediaQuery.matches;
  }

  watchDarkMode(callback: (darkMode: boolean) => void) {
    this.darkModeMediaQuery.addEventListener('change', (e) => {
      callback(e.matches);
    });
  }

  toggleBodyClass(darkMode: boolean) {
    const body = document.body;

    if (darkMode) {
      this.renderer.addClass(body, 'dark-theme');
    } else {
      this.renderer.removeClass(body, 'dark-theme');
    }
  }
}
