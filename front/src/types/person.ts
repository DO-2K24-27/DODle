export interface Person {
  firstname: string;
  lastname: string;
  gender: string;
  type: string;
  workplace: string;
  image: string;
  hint: string;
}

export interface GuessResult {
  correct: boolean;
  person: Person;
}

export interface GuessHistory {
  guess: Person;
  result: GuessResult;
  guessedPerson: Person;
}
