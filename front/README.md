# DODle - Front-end

A Vue.js application for the DODle guessing game, similar to Wordle but for guessing people by their characteristics.

## Features

- **Person Selection**: Choose from a list of available persons
- **Guess Submission**: Submit your guess and get feedback
- **Visual Feedback**: Green for correct characteristics, Red for incorrect
- **Guess History**: View all previous guesses with color-coded results
- **Hint System**: Get hints when you're stuck
- **Responsive Design**: Works on desktop and mobile devices

## Game Rules

1. The goal is to guess the "Person of the Day"
2. Select a person from the dropdown and submit your guess
3. You'll get feedback on each characteristic:
   - **Green**: Correct characteristic
   - **Red**: Incorrect characteristic
4. Use the feedback to make better guesses
5. Get hints if you need help
6. Keep guessing until you find the correct person!

## Characteristics

Each person has the following characteristics:
- **First Name**
- **Last Name**
- **Gender**
- **Type** (role/profession)
- **Workplace**

## Development

### Prerequisites

- Node.js (v18 or higher)
- npm

### Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Build for production:
   ```bash
   npm run build
   ```

## Docker Deployment

You can build and run the DODle frontend in a Docker container.

### Building the Docker image

```bash
# Navigate to the front directory
cd front

# Build the Docker image
docker build -t dodle-frontend .
```

### Running the container

```bash
# Run the frontend container
docker run -p 80:80 dodle-frontend
```

### Using Docker Compose

The project includes a docker-compose.yml file in the root directory that sets up the complete application (frontend, API, and MongoDB):

```bash
# Navigate to the project root
cd ..

# Start all services
docker-compose up

# Or build and start in detached mode
docker-compose up -d --build
```

### Accessing the application

Once the container is running, access the application at [http://localhost](http://localhost)

### API Integration

The front-end communicates with the Go API backend on port 8080. The Vite development server automatically proxies API requests to avoid CORS issues.

API endpoints used:
- `GET /public/v1/persons` - Get all available persons
- `POST /public/v1/guess/person/submit` - Submit a guess
- `GET /public/v1/guess/person/hint` - Get a hint

## Technologies Used

- **Vue.js 3** - Progressive JavaScript framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Fast build tool
- **Vue Router** - Client-side routing
- **CSS3** - Styling with modern CSS features

## Project Structure

```
src/
├── components/
│   └── GuessGame.vue     # Main game component
├── services/
│   └── api.ts           # API service functions
├── types/
│   └── person.ts        # TypeScript interfaces
├── views/
│   └── GameView.vue     # Main game view
├── App.vue              # Root component
└── main.ts              # Application entry point
```

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
npm run test:unit
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```
