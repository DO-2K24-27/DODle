<template>
    <div class="game-container">
        <header class="game-header">
            <h1>DODle</h1>
            <p class="game-subtitle">
                Find the person by guessing their characteristics. Green = correct, Red = incorrect.
            </p>
        </header>

        <!-- Yesterday's Person Section -->
        <div v-if="yesterdaysPerson && !gameWon" class="yesterdays-person">
            <div class="yesterday-content">
                <p>Yesterday's person was: {{ yesterdaysPerson.firstname }} {{ yesterdaysPerson.lastname }}</p>
            </div>
        </div>

        <div class="game-content">
            <!-- Person Selection -->
            <div class="person-selector" v-if="!gameWon">
                <h2>Make your guess</h2>
                <div class="form-group">
                    <label for="person-select">Select a person:</label>
                    <select id="person-select" v-model="selectedPersonIndex" class="person-select" :disabled="loading">
                        <option value="">Choose a person...</option>
                        <option v-for="(person, index) in availablePersons" :key="index" :value="getOriginalIndex(person)">
                            {{ person.firstname }} {{ person.lastname }}
                        </option>
                    </select>
                </div>

                <button @click="submitGuess" :disabled="selectedPersonIndex === '' || loading" class="submit-button">
                    {{ loading ? 'Submitting...' : 'Submit Guess' }}
                </button>

                <button @click="getHint" :disabled="loading" class="hint-button">
                    {{ loading ? 'Getting Hint...' : 'Get Hint' }}
                </button>
            </div>

            <!-- Hint Display -->
            <div v-if="hint" class="hint-display">
                <h3>💡 Hint</h3>
                <p>{{ hint }}</p>
            </div>

            <!-- Game Status -->
            <div v-if="gameWon" class="game-status won">
                <h2>🎉 Congratulations! You found the person!</h2>
                <p>It was {{ correctPerson?.firstname }} {{ correctPerson?.lastname }}!</p>
            </div>

            <!-- Error Display -->
            <div v-if="error" class="error-message">
                <p>{{ error }}</p>
            </div>

            <!-- Guess History -->
            <div v-if="guessHistory.length > 0" class="guess-history">
                <h2>Your Guesses</h2>
                <div class="history-table">
                    <div class="table-header">
                        <div class="header-cell">Guess</div>
                        <div class="header-cell">First Name</div>
                        <div class="header-cell">Last Name</div>
                        <div class="header-cell">Gender</div>
                        <div class="header-cell">Type</div>
                        <div class="header-cell">Workplace</div>
                    </div>

                    <div v-for="(historyItem, index) in guessHistory" :key="index" class="table-row">
                        <div class="cell guess-number" :class="(getFieldClass(historyItem.result.person.firstname, historyItem.guessedPerson.firstname) == 'correct' && getFieldClass(historyItem.result.person.lastname, historyItem.guessedPerson.lastname) == 'correct' ? 'correct' : 'incorrect')">{{ index + 1 }}</div>
                        <div class="cell"
                            :class="getFieldClass(historyItem.result.person.firstname, historyItem.guessedPerson.firstname)">
                            {{ historyItem.guessedPerson.firstname }}
                        </div>
                        <div class="cell"
                            :class="getFieldClass(historyItem.result.person.lastname, historyItem.guessedPerson.lastname)">
                            {{ historyItem.guessedPerson.lastname }}
                        </div>
                        <div class="cell"
                            :class="getFieldClass(historyItem.result.person.gender, historyItem.guessedPerson.gender)">
                            {{ historyItem.guessedPerson.gender }}
                        </div>
                        <div class="cell"
                            :class="getFieldClass(historyItem.result.person.type, historyItem.guessedPerson.type)">
                            {{ historyItem.guessedPerson.type }}
                        </div>
                        <div class="cell"
                            :class="getFieldClass(historyItem.result.person.workplace, historyItem.guessedPerson.workplace)">
                            {{ historyItem.guessedPerson.workplace }}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { apiService } from '@/services/api';
import type { Person, GuessHistory } from '@/types/person';

const persons = ref<Person[]>([]);
const selectedPersonIndex = ref<string>('');
const guessHistory = ref<GuessHistory[]>([]);
const loading = ref(false);
const error = ref<string>('');
const hint = ref<string>('');
const gameWon = ref(false);
const correctPerson = ref<Person | null>(null);
const yesterdaysPerson = ref<Person | null>(null);
const notification = ref<string>('');
const showNotification = ref<boolean>(false);

// Storage keys
const STORAGE_KEY_HISTORY = 'dodle_game_history';
const STORAGE_KEY_GUESS_ID = 'dodle_guess_id';
const STORAGE_KEY_HINT = 'dodle_hint';

// Computed property to get persons that haven't been guessed yet
const availablePersons = computed(() => {
    const guessedPersons = guessHistory.value.map(history => 
        `${history.guessedPerson.firstname} ${history.guessedPerson.lastname}`
    );
    return persons.value.filter(person => 
        !guessedPersons.includes(`${person.firstname} ${person.lastname}`)
    );
});

// Helper function to get original index from persons array
const getOriginalIndex = (person: Person) => {
    return persons.value.findIndex(p => 
        p.firstname === person.firstname && p.lastname === person.lastname
    );
};

// Save game state to local storage
const saveGameState = (guessId: string) => {
    // Save guess ID
    localStorage.setItem(STORAGE_KEY_GUESS_ID, guessId);
    
    // Save guess history
    localStorage.setItem(STORAGE_KEY_HISTORY, JSON.stringify(guessHistory.value));
    
    // Save hint
    localStorage.setItem(STORAGE_KEY_HINT, hint.value);
};

// Load game state from local storage
const loadGameState = async () => {
    try {
        // Get current guess ID from API
        const currentGuessId = await apiService.getGuessID();
        
        // Get saved guess ID from local storage
        const savedGuessId = localStorage.getItem(STORAGE_KEY_GUESS_ID);
        
        // If IDs don't match or no saved ID, clear storage and start new game
        if (!savedGuessId || savedGuessId !== currentGuessId) {
            clearGameState();
            return false;
        }
        
        // IDs match, load the saved game state
        const savedHistoryJson = localStorage.getItem(STORAGE_KEY_HISTORY);
        const savedHint = localStorage.getItem(STORAGE_KEY_HINT);
        
        if (savedHistoryJson) {
            guessHistory.value = JSON.parse(savedHistoryJson);
            
            // Check if the game was already won
            const lastGuess = guessHistory.value[guessHistory.value.length - 1];
            if (lastGuess && lastGuess.result.correct) {
                gameWon.value = true;
                correctPerson.value = lastGuess.guessedPerson;
            }
        }
        
        if (savedHint) {
            hint.value = savedHint;
        }
        
        return true;
    } catch (err) {
        console.error('Error loading game state:', err);
        clearGameState();
        return false;
    }
};

// Clear game state from local storage
const clearGameState = () => {
    localStorage.removeItem(STORAGE_KEY_GUESS_ID);
    localStorage.removeItem(STORAGE_KEY_HISTORY);
    localStorage.removeItem(STORAGE_KEY_HINT);
    guessHistory.value = [];
    hint.value = '';
    gameWon.value = false;
    correctPerson.value = null;
};

// Watch for changes to guess history and save it
watch([guessHistory, hint], async () => {
    if (guessHistory.value.length > 0) {
        try {
            const currentGuessId = await apiService.getGuessID();
            saveGameState(currentGuessId);
        } catch (err) {
            console.error('Error saving game state:', err);
        }
    }
}, { deep: true });

onMounted(async () => {
    await Promise.all([
        loadPersons(),
        loadYesterdaysPerson()
    ]);
    
    // Load saved game state after persons are loaded
    const gameStateLoaded = await loadGameState();
    
    if (gameStateLoaded && guessHistory.value.length > 0) {
        notification.value = "Game state loaded from your previous session";
        showNotification.value = true;
        
        // Hide notification after 5 seconds
        setTimeout(() => {
            showNotification.value = false;
        }, 5000);
    }
});


const loadPersons = async () => {
    try {
        loading.value = true;
        error.value = '';
        const response = await apiService.getPersons();
        persons.value = response.persons;
    } catch (err) {
        error.value = 'Failed to load persons. Please try again.';
        console.error('Error loading persons:', err);
    } finally {
        loading.value = false;
    }
};

const loadYesterdaysPerson = async () => {
    try {
        const response = await apiService.getYesterdaysPerson();
        yesterdaysPerson.value = response;
    } catch (err) {
        console.error('Error loading yesterday\'s person:', err);
        // We don't set error state here as it's not critical to game functionality
    }
};

const submitGuess = async () => {
    if (selectedPersonIndex.value === '') return;

    const selectedPerson = persons.value[parseInt(selectedPersonIndex.value)];

    try {
        loading.value = true;
        error.value = '';

        const result = await apiService.submitGuess(selectedPerson);

        guessHistory.value.push({
            guess: selectedPerson,
            result: result,
            guessedPerson: selectedPerson
        });

        if (result.correct) {
            gameWon.value = true;
            correctPerson.value = selectedPerson;
        }

        // No need to explicitly call saveGameState here as it's handled by the watcher

        selectedPersonIndex.value = '';

    } catch (err) {
        error.value = 'Failed to submit guess. Please try again.';
        console.error('Error submitting guess:', err);
    } finally {
        loading.value = false;
    }
};

const getHint = async () => {
    try {
        loading.value = true;
        error.value = '';

        const response = await apiService.getHint();
        hint.value = response;
        
        // No need to explicitly call saveGameState here as it's handled by the watcher

    } catch (err) {
        error.value = 'Failed to get hint. Please try again.';
        console.error('Error getting hint:', err);
    } finally {
        loading.value = false;
    }
};

const getFieldClass = (resultField: string, guessedField: string) => {
    if (resultField && resultField === guessedField) {
        return 'correct';
    } else if (resultField === '') {
        return 'incorrect';
    }
    return '';
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&family=Outfit:wght@100..900&display=swap');

.game-container {
    width: 100%;
    margin: 0;
    font-family: Outfit;
    color: #ffffff;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding-bottom: 100px;
}

.game-header {
    text-align: center;
    margin-bottom: 30px;
    padding: 20px;
    background: #1d1d20;
    width: 100%;
}

.game-header h1 {
    font-size: 2.5em;
    margin-bottom: 10px;
    text-transform: uppercase;
}

.game-subtitle {
    color: #7f8c8d;
    font-size: 1.1em;
    margin-bottom: 20px;
}

.yesterdays-person {
    width: 500px;
    background: #1d1d20;
    padding: 12px 20px;
    margin-bottom: 30px;
    display: flex;
    justify-content: center;
    text-align: center;
    border-radius: 8px;
}

.yesterday-content {
    max-width: 800px;
}

.yesterdays-person p {
    color: #7f8c8d;
    font-size: 1em;
    font-weight: 500;
}

.yesterday-name {
    font-size: 1.2em;
    font-weight: 700;
    margin-bottom: 5px;
    color: #ffffff;
}

.yesterday-details {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    gap: 15px;
    color: #b8b8b8;
    font-size: 0.9em;
}

.detail strong {
    color: #ffffff;
}

.game-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 30px;
    width: fit-content;
}

.person-selector {
    background: #1d1d20;
    padding: 20px;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 500px;
}

.person-selector h2 {
    margin-bottom: 15px;
}

.form-group {
    margin-bottom: 15px;
}

.form-group label {
    display: block;
    margin-bottom: 5px;
    font-weight: 500;
}

.person-select {
    width: 300px;
    padding: 10px;
    border: 1px solid #ced4da;
    border-radius: 4px;
    font-size: 16px;
    background-color: white;
}

.person-select:disabled {
    background-color: #f8f9fa;
    cursor: not-allowed;
}

.submit-button,
.hint-button {
    padding: 12px 24px;
    margin-right: 10px;
    margin-bottom: 10px;
    border: none;
    border-bottom: 2px solid #2b2b2b;
    border-radius: 4px;
    font-size: 16px;
    cursor: pointer;
    transition: background-color 0.2s;
    width: 200px;
    font-weight: 500;
}

.submit-button {
    background-color: #28a745;
    color: white;
}

.submit-button:hover:not(:disabled) {
    background-color: #1e7e34;
}

.hint-button {
    background-color: #3a3a3c;
    border: none;
    color: white;
}

.hint-button:hover:not(:disabled) {
    background-color: #2c2c2e;
}

.submit-button:disabled,
.hint-button:disabled {
    background-color: #6c757d;
    cursor: not-allowed;
}

.hint-display {
    background: #1d1d20;
    padding: 15px;
    border-radius: 8px;
    border-left: 4px solid #007bff;
}

.hint-display h3 {
    margin-bottom: 10px;
}

.game-status.won {
    background: #d4edda;
    color: #155724;
    padding: 20px;
    border-radius: 8px;
    text-align: center;
    border: 1px solid #c3e6cb;
}

.notification {
    background: #cce5ff;
    color: #004085;
    padding: 15px;
    border-radius: 8px;
    border: 1px solid #b8daff;
    margin-bottom: 15px;
    animation: fadeInOut 5s ease-in-out;
}

@keyframes fadeInOut {
    0% { opacity: 0; }
    10% { opacity: 1; }
    90% { opacity: 1; }
    100% { opacity: 0; }
}

.error-message {
    background: #f8d7da;
    color: #721c24;
    padding: 15px;
    border-radius: 8px;
    border: 1px solid #f5c6cb;
}

.guess-history {
    background: #1d1d20;
    padding: 20px;
    border-radius: 8px;
}

.guess-history h2 {
    margin-bottom: 20px;
}

.history-table {
    display: flex;
    flex-direction: column;
    gap: 1px;
    background: #dee2e6;
    border-radius: 4px;
    overflow: hidden;
}

.table-header {
    display: grid;
    grid-template-columns: 80px 1fr 1fr 1fr 1fr 1fr;
    gap: 1px;
    background: #262627;
}

.header-cell {
    padding: 12px;
    font-weight: 600;
    text-align: center;
    background: #3a3a3c;
}

.table-row {
    display: grid;
    grid-template-columns: 80px 1fr 1fr 1fr 1fr 1fr;
    gap: 1px;
}

.cell {
    padding: 10px;
    background: #3a3a3c;
    text-align: center;
    transition: background-color 0.2s;
    text-transform: uppercase;
    color: white;
    font-weight: 500;
}

.guess-number {
    font-weight: 600;
}

.cell.correct {
    background-color: #3eaa42;
    color: white;
    font-weight: 600;
}

.cell.incorrect {
    background-color: #3a3a3c;
}

@media (max-width: 768px) {

    .table-header,
    .table-row {
        grid-template-columns: 60px 1fr 1fr;
    }

    .header-cell:nth-child(n+4),
    .cell:nth-child(n+4) {
        display: none;
    }
}
</style>
