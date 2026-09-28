<script lang="ts">
	import QuestionList from './QuestionList.svelte';
	import type { Question } from '$lib/types/question';
	export let questions: Question[] = [];
</script>

<QuestionList {questions} emptyMessage="No active questions yet." />
