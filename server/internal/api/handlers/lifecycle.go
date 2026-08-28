package handlers

// Lifecycle mutations intentionally use PUT /api/v1/questions/{questionId}.
// Keeping one canonical question endpoint prevents divergent answer copies;
// the capture handler delegates validation to the authoritative store.
