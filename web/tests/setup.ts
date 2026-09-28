import '@testing-library/jest-dom/vitest';
import { cleanup, fireEvent } from '@testing-library/svelte';
import { afterEach } from 'vitest';

// Keep the existing component tests compatible with both event-name spellings.
const fireEventWithAliases = fireEvent as typeof fireEvent & { mouseup: typeof fireEvent.mouseUp };
fireEventWithAliases.mouseup = fireEvent.mouseUp;

afterEach(() => cleanup());
