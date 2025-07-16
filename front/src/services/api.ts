import type { Person, GuessResult } from '@/types/person';

// Use relative path for API calls when using nginx proxy, or full URL for direct calls
const API_BASE_URL = 'PLACEHOLDER_API_URL'; // This will be replaced with actual API URL

export const apiService = {
  async getPersons(): Promise<{ persons: Person[] }> {
    const response = await fetch(`${API_BASE_URL}/public/v1/persons`);
    if (!response.ok) {
      throw new Error('Failed to fetch persons');
    }
    return response.json();
  },

  async submitGuess(person: Person): Promise<GuessResult> {
    const response = await fetch(`${API_BASE_URL}/public/v1/guess/person/submit`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(person),
    });
    
    if (!response.ok) {
      throw new Error('Failed to submit guess');
    }
    return response.json();
  },

  async getHint(): Promise<string> {
    const response = await fetch(`${API_BASE_URL}/public/v1/guess/person/hint`);
    if (!response.ok) {
      throw new Error('Failed to get hint');
    }
    return response.json();
  },

  async getYesterdaysPerson(): Promise<Person> {
    const response = await fetch(`${API_BASE_URL}/public/v1/guess/person/yesterday`);
    if (!response.ok) {
      throw new Error('Failed to fetch yesterday\'s person');
    }
    return response.json();
  },
  
  async getGuessID(): Promise<string> {
    const response = await fetch(`${API_BASE_URL}/public/v1/guess/id`);
    if (!response.ok) {
      throw new Error('Failed to fetch guess ID');
    }
    const data = await response.json();
    return data.id; // Extract the ID from the wrapper object
  },
};
