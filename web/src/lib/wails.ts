/**
 * Wails Runtime Integration
 *
 * This module provides a unified interface for interacting with the Wails desktop runtime.
 * When running as a web app, it provides fallback implementations or throws appropriate errors.
 */

// Check if we're running in Wails desktop environment
export const isDesktopApp = (): boolean => {
  if (typeof window === 'undefined') return false;
  const win = window as unknown as WailsWindow;
  return 'go' in window && win.go && 'main' in win.go;
};

// Type definitions for Wails runtime
interface WailsWindow extends Window {
  go: {
    main: {
      App: AppBindings;
    };
  };
  runtime: WailsRuntime;
}

interface WailsRuntime {
  EventsOn(eventName: string, callback: (...args: unknown[]) => void): () => void;
  EventsOff(eventName: string): void;
  EventsEmit(eventName: string, ...args: unknown[]): void;
  WindowSetTitle(title: string): void;
  WindowMinimize(): void;
  WindowMaximize(): void;
  WindowUnmaximize(): void;
  WindowToggleMaximise(): void;
  WindowSetMinSize(width: number, height: number): void;
  WindowSetMaxSize(width: number, height: number): void;
  WindowSetSize(width: number, height: number): void;
  WindowGetSize(): Promise<{ w: number; h: number }>;
  WindowCenter(): void;
  WindowShow(): void;
  WindowHide(): void;
  WindowIsFullscreen(): Promise<boolean>;
  WindowIsMaximised(): Promise<boolean>;
  WindowIsMinimised(): Promise<boolean>;
  WindowIsNormal(): Promise<boolean>;
  Quit(): void;
  BrowserOpenURL(url: string): void;
}

interface AppBindings {
  GetAPIBaseURL(): Promise<string>;
  SetAPIBaseURL(url: string): Promise<void>;
  GetAppVersion(): Promise<string>;
  GetPrinters(): Promise<PrinterInfo[]>;
  GetDefaultPrinter(): Promise<string>;
  GetThermalPrinter(): Promise<string>;
  SetThermalPrinter(printerName: string): Promise<void>;
  PrintDocument(pdfBase64: string): Promise<void>;
  PrintThermalTicket(pdfBase64: string): Promise<void>;
  PrintToSpecificPrinter(pdfBase64: string, printerName: string): Promise<void>;
  ShowSaveDialog(defaultFilename: string, filterDescription: string, filterPattern: string): Promise<string>;
  SaveFile(contentBase64: string, defaultFilename: string): Promise<void>;
  GetConfigPath(): Promise<string>;
  GetPlatform(): Promise<string>;
  OpenExternalURL(url: string): Promise<void>;
  GetDataDirectory(): Promise<string>;
  IsDesktopApp(): Promise<boolean>;
}

export interface PrinterInfo {
  name: string;
  isDefault: boolean;
  status: string;
}

// Get the Wails window object
const getWailsWindow = (): WailsWindow | null => {
  if (isDesktopApp()) {
    return window as unknown as WailsWindow;
  }
  return null;
};

// Get the Wails app bindings
const getAppBindings = (): AppBindings | null => {
  const wailsWindow = getWailsWindow();
  return wailsWindow?.go?.main?.App ?? null;
};

// Get the Wails runtime
const getRuntime = (): WailsRuntime | null => {
  const wailsWindow = getWailsWindow();
  return wailsWindow?.runtime ?? null;
};

/**
 * Desktop App API - provides methods bound from the Go backend
 */
export const desktopApp = {
  /**
   * Get the configured API base URL
   */
  async getAPIBaseURL(): Promise<string> {
    const app = getAppBindings();
    if (!app) {
      // Fallback for web: use environment variable or default
      return import.meta.env.VITE_API_BASE_URL || '/api/v1';
    }
    return app.GetAPIBaseURL();
  },

  /**
   * Set the API base URL
   */
  async setAPIBaseURL(url: string): Promise<void> {
    const app = getAppBindings();
    if (!app) {
      console.warn('setAPIBaseURL is only available in desktop mode');
      return;
    }
    return app.SetAPIBaseURL(url);
  },

  /**
   * Get the application version
   */
  async getAppVersion(): Promise<string> {
    const app = getAppBindings();
    if (!app) {
      return '0.0.1-web';
    }
    return app.GetAppVersion();
  },

  /**
   * Get list of available printers
   */
  async getPrinters(): Promise<PrinterInfo[]> {
    const app = getAppBindings();
    if (!app) {
      throw new Error('Printer functions are only available in desktop mode');
    }
    return app.GetPrinters();
  },

  /**
   * Get the default system printer
   */
  async getDefaultPrinter(): Promise<string> {
    const app = getAppBindings();
    if (!app) {
      throw new Error('Printer functions are only available in desktop mode');
    }
    return app.GetDefaultPrinter();
  },

  /**
   * Get the configured thermal printer
   */
  async getThermalPrinter(): Promise<string> {
    const app = getAppBindings();
    if (!app) {
      throw new Error('Printer functions are only available in desktop mode');
    }
    return app.GetThermalPrinter();
  },

  /**
   * Set the thermal printer for receipts
   */
  async setThermalPrinter(printerName: string): Promise<void> {
    const app = getAppBindings();
    if (!app) {
      throw new Error('Printer functions are only available in desktop mode');
    }
    return app.SetThermalPrinter(printerName);
  },

  /**
   * Print a PDF document to the default printer
   * @param pdfBytes - The PDF content as Uint8Array or ArrayBuffer
   */
  async printDocument(pdfBytes: Uint8Array | ArrayBuffer): Promise<void> {
    const app = getAppBindings();
    if (!app) {
      throw new Error('Print functions are only available in desktop mode');
    }
    const base64 = arrayBufferToBase64(pdfBytes);
    return app.PrintDocument(base64);
  },

  /**
   * Print a PDF to the thermal printer
   * @param pdfBytes - The PDF content as Uint8Array or ArrayBuffer
   */
  async printThermalTicket(pdfBytes: Uint8Array | ArrayBuffer): Promise<void> {
    const app = getAppBindings();
    if (!app) {
      throw new Error('Print functions are only available in desktop mode');
    }
    const base64 = arrayBufferToBase64(pdfBytes);
    return app.PrintThermalTicket(base64);
  },

  /**
   * Print a PDF to a specific printer
   * @param pdfBytes - The PDF content as Uint8Array or ArrayBuffer
   * @param printerName - The name of the printer to use
   */
  async printToSpecificPrinter(pdfBytes: Uint8Array | ArrayBuffer, printerName: string): Promise<void> {
    const app = getAppBindings();
    if (!app) {
      throw new Error('Print functions are only available in desktop mode');
    }
    const base64 = arrayBufferToBase64(pdfBytes);
    return app.PrintToSpecificPrinter(base64, printerName);
  },

  /**
   * Show a native save file dialog
   */
  async showSaveDialog(defaultFilename: string, filterDescription: string, filterPattern: string): Promise<string> {
    const app = getAppBindings();
    if (!app) {
      throw new Error('File dialogs are only available in desktop mode');
    }
    return app.ShowSaveDialog(defaultFilename, filterDescription, filterPattern);
  },

  /**
   * Save file content using a native dialog
   * @param content - The file content as Uint8Array or ArrayBuffer
   * @param defaultFilename - Default filename for the save dialog
   */
  async saveFile(content: Uint8Array | ArrayBuffer, defaultFilename: string): Promise<void> {
    const app = getAppBindings();
    if (!app) {
      // Fallback for web: use browser download
      const blob = new Blob([content]);
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = defaultFilename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
      return;
    }
    const base64 = arrayBufferToBase64(content);
    return app.SaveFile(base64, defaultFilename);
  },

  /**
   * Open a URL in the default browser
   */
  async openExternalURL(url: string): Promise<void> {
    const app = getAppBindings();
    if (!app) {
      window.open(url, '_blank');
      return;
    }
    return app.OpenExternalURL(url);
  },

  /**
   * Get the current platform
   */
  async getPlatform(): Promise<string> {
    const app = getAppBindings();
    if (!app) {
      return 'web';
    }
    return app.GetPlatform();
  },
};

/**
 * Wails Window API - provides window control methods
 */
export const windowAPI = {
  /**
   * Set the window title
   */
  setTitle(title: string): void {
    const runtime = getRuntime();
    if (runtime) {
      runtime.WindowSetTitle(title);
    } else {
      document.title = title;
    }
  },

  /**
   * Minimize the window
   */
  minimize(): void {
    const runtime = getRuntime();
    runtime?.WindowMinimize();
  },

  /**
   * Maximize the window
   */
  maximize(): void {
    const runtime = getRuntime();
    runtime?.WindowMaximize();
  },

  /**
   * Toggle maximize state
   */
  toggleMaximize(): void {
    const runtime = getRuntime();
    runtime?.WindowToggleMaximise();
  },

  /**
   * Center the window
   */
  center(): void {
    const runtime = getRuntime();
    runtime?.WindowCenter();
  },

  /**
   * Quit the application
   */
  quit(): void {
    const runtime = getRuntime();
    runtime?.Quit();
  },

  /**
   * Open URL in browser
   */
  openURL(url: string): void {
    const runtime = getRuntime();
    if (runtime) {
      runtime.BrowserOpenURL(url);
    } else {
      window.open(url, '_blank');
    }
  },
};

/**
 * Event system for communication between Go backend and React frontend
 */
export const events = {
  /**
   * Subscribe to an event
   */
  on(eventName: string, callback: (...args: unknown[]) => void): () => void {
    const runtime = getRuntime();
    if (runtime) {
      return runtime.EventsOn(eventName, callback);
    }
    // No-op for web
    return () => {};
  },

  /**
   * Unsubscribe from an event
   */
  off(eventName: string): void {
    const runtime = getRuntime();
    runtime?.EventsOff(eventName);
  },

  /**
   * Emit an event
   */
  emit(eventName: string, ...args: unknown[]): void {
    const runtime = getRuntime();
    runtime?.EventsEmit(eventName, ...args);
  },
};

// Helper function to convert ArrayBuffer to base64
function arrayBufferToBase64(buffer: Uint8Array | ArrayBuffer): string {
  const bytes = buffer instanceof Uint8Array ? buffer : new Uint8Array(buffer);
  let binary = '';
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}

// Export type for use in components
export type { AppBindings, WailsRuntime };
