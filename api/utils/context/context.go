package ctxutil

// Define a custom type for context keys to avoid collisions
type ContextKey string

const MongoClientKey ContextKey = "mongoClient"
