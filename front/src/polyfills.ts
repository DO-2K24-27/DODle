// src/polyfills.ts
// This file contains polyfills for browser environments
import CryptoJS from 'crypto-js';

// Declare types for global extensions
declare global {
  interface Window {
    global: typeof globalThis;
  }
  
  interface Crypto {
    hash?: (algorithm: string, data: ArrayBuffer) => Uint8Array;
  }
}

// Create window.global for use in the app
window.global = window as typeof globalThis;

// Handle missing crypto functionality
if (window.crypto) {
  // Add hash function if it doesn't exist in the browser's crypto implementation
  if (!('hash' in window.crypto)) {
    Object.defineProperty(window.crypto, 'hash', {
      value: function(algorithm: string, data: ArrayBuffer): Uint8Array {
        // Use crypto-js to perform actual hashing
        let wordArray;
        const dataView = new Uint8Array(data);
        const dataHex = Array.from(dataView)
          .map(b => b.toString(16).padStart(2, '0'))
          .join('');
        
        // Choose hashing algorithm based on the algorithm parameter
        switch (algorithm.toLowerCase()) {
          case 'sha-256':
            wordArray = CryptoJS.SHA256(CryptoJS.enc.Hex.parse(dataHex));
            break;
          case 'sha-1':
            wordArray = CryptoJS.SHA1(CryptoJS.enc.Hex.parse(dataHex));
            break;
          case 'md5':
            wordArray = CryptoJS.MD5(CryptoJS.enc.Hex.parse(dataHex));
            break;
          default:
            console.warn(`Hash algorithm ${algorithm} not supported, using SHA-256`);
            wordArray = CryptoJS.SHA256(CryptoJS.enc.Hex.parse(dataHex));
        }
        
        // Convert WordArray to Uint8Array
        const hexStr = wordArray.toString(CryptoJS.enc.Hex);
        const result = new Uint8Array(hexStr.length / 2);
        for (let i = 0; i < hexStr.length; i += 2) {
          result[i / 2] = parseInt(hexStr.substr(i, 2), 16);
        }
        
        return result;
      },
      writable: false,
      configurable: true
    });
  }
}

export default {};
